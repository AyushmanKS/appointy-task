// backend/internal/link/handlers.go
package link

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/AyushmanKS/appointy-task/internal/auth"
	"github.com/AyushmanKS/appointy-task/internal/database"
	"github.com/AyushmanKS/appointy-task/internal/hub" // Import the hub
)

// CreateLinkHandler handles the creation of a new short URL for an authenticated user.
func CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Could not retrieve user ID from token", http.StatusInternalServerError)
		return
	}

	var data struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL(data.URL)
	query := "INSERT INTO urls (id, original_url, user_id) VALUES ($1, $2, $3)"
	_, err := database.DB.ExecContext(r.Context(), query, shortURL, data.URL, userID)
	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	fullShortURL := fmt.Sprintf("https://%s/r/%s", r.Host, shortURL)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"short_url": fullShortURL})
}

// RedirectHandler finds the original URL and redirects. It also records the click asynchronously.
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/r/"):]
	var originalURL string
	query := "SELECT original_url FROM urls WHERE id = $1"

	err := database.DB.QueryRowContext(r.Context(), query, id).Scan(&originalURL)
	if err != nil {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	go recordClick(id, r)
	http.Redirect(w, r, originalURL, http.StatusFound)
}

// recordClick saves the click and broadcasts the new total count via the WebSocket Hub.
func recordClick(linkID string, r *http.Request) {
	// First, insert the click record.
	queryInsert := "INSERT INTO clicks (url_id, ip_address, user_agent) VALUES ($1, $2, $3)"
	_, err := database.DB.ExecContext(context.Background(), queryInsert, linkID, r.RemoteAddr, r.UserAgent())
	if err != nil {
		log.Printf("Failed to record click for link %s: %v", linkID, err)
		return
	}

	// Now, query the new total count and find out who owns the link.
	var totalClicks int
	var userID int
	queryCount := `
		SELECT count(c.id), u.user_id
		FROM clicks c
		JOIN urls u ON c.url_id = u.id
		WHERE c.url_id = $1
		GROUP BY u.user_id`

	err = database.DB.QueryRowContext(context.Background(), queryCount, linkID).Scan(&totalClicks, &userID)
	if err != nil {
		log.Printf("Failed to query new click count for link %s: %v", linkID, err)
		return
	}

	// Tell the global hub to send an update to this specific user.
	hub.GlobalHub.BroadcastUpdate(userID, linkID, totalClicks)
}

// GetAnalyticsHandler retrieves click data for a link.
func GetAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Could not retrieve user ID from token", http.StatusInternalServerError)
		return
	}
	linkID := r.URL.Path[len("/analytics/"):]

	var totalClicks int
	query := `SELECT count(c.id) FROM clicks c JOIN urls u ON c.url_id = u.id WHERE c.url_id = $1 AND u.user_id = $2`
	err := database.DB.QueryRowContext(r.Context(), query, linkID, userID).Scan(&totalClicks)
	if err != nil {
		http.Error(w, "Could not retrieve analytics", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{"total_clicks": totalClicks})
}

// GetLinksHandler retrieves all links for the authenticated user.
func GetLinksHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Could not retrieve user ID from token", http.StatusInternalServerError)
		return
	}

	type Link struct {
		ShortID     string `json:"short_id"`
		OriginalURL string `json:"original_url"`
	}
	var links []Link

	query := "SELECT id, original_url FROM urls WHERE user_id = $1 ORDER BY creation_date DESC"
	rows, err := database.DB.QueryContext(r.Context(), query, userID)
	if err != nil {
		http.Error(w, "Could not retrieve links", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var link Link
		if err := rows.Scan(&link.ShortID, &link.OriginalURL); err != nil {
			http.Error(w, "Error scanning links", http.StatusInternalServerError)
			return
		}
		links = append(links, link)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}

func generateShortURL(originalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(originalURL))
	return hex.EncodeToString(hasher.Sum(nil))[:8]
}
