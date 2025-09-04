package hub

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	LinkID     string `json:"link_id"`
	ClickCount int    `json:"click_count"`
}

type Hub struct {
	clients    map[int]*websocket.Conn
	mu         sync.Mutex
	register   chan *Client
	unregister chan *Client
}

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

func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

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
		go h.UnregisterClient(&Client{UserID: userID, Conn: conn})
	}
}
