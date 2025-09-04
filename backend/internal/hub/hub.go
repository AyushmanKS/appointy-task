// backend/internal/hub/hub.go
package hub

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Message defines the structure of a real-time update.
type Message struct {
	LinkID     string `json:"link_id"`
	ClickCount int    `json:"click_count"`
}

// Hub manages WebSocket clients. The fields are lowercase (private).
type Hub struct {
	clients    map[int]*websocket.Conn
	mu         sync.Mutex
	register   chan *Client
	unregister chan *Client
}

// Client represents a connected user.
type Client struct {
	UserID int
	Conn   *websocket.Conn
}

var GlobalHub *Hub

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int]*websocket.Conn),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func init() {
	GlobalHub = NewHub()
}

// Run starts the hub's internal event loop.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client.Conn
			h.mu.Unlock()
			log.Printf("Client registered: UserID %d", client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				client.Conn.Close()
				log.Printf("Client unregistered: UserID %d", client.UserID)
			}
			h.mu.Unlock()
		}
	}
}

// --- NEW PUBLIC METHODS ---

// RegisterClient is an exported (public) method to register a client.
// Other packages will call this method.
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient is an exported (public) method to unregister a client.
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// BroadcastUpdate sends a message to a specific connected user.
func (h *Hub) BroadcastUpdate(userID int, linkID string, count int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	conn, ok := h.clients[userID]
	if !ok {
		return
	}

	message := Message{LinkID: linkID, ClickCount: count}
	err := conn.WriteJSON(message)
	if err != nil {
		log.Printf("Error sending update to UserID %d: %v", userID, err)
		// Unregister the client if there's an error writing to the connection.
		go h.UnregisterClient(&Client{UserID: userID, Conn: conn})
	}
}
