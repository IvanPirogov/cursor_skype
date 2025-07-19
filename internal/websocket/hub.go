package websocket

import (
	"encoding/json"
	"log"
	"messenger/internal/auth"
	"messenger/pkg/models"
	"net/http"
	"sync"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients    map[uuid.UUID]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mutex      sync.RWMutex
	db         *gorm.DB
}

type Client struct {
	ID     uuid.UUID
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	UserID uuid.UUID
}

type Message struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	UserID    uuid.UUID   `json:"user_id"`
	Timestamp int64       `json:"timestamp"`
}

const (
	MessageTypeChat         = "chat"
	MessageTypeUserStatus   = "user_status"
	MessageTypeTyping       = "typing"
	MessageTypeCallOffer    = "call_offer"
	MessageTypeCallAnswer   = "call_answer"
	MessageTypeCallReject   = "call_reject"
	MessageTypeCallEnd      = "call_end"
	MessageTypeNewMessage   = "new_message"
	MessageTypeMessageRead  = "message_read"
	MessageTypeUserJoined   = "user_joined"
	MessageTypeUserLeft     = "user_left"
	MessageTypeNewContact   = "new_contact"
)

func NewHub(db *gorm.DB) *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		db:         db,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client.UserID] = client
			h.mutex.Unlock()
			
			// Update user status to online
			h.db.Model(&models.User{}).Where("id = ?", client.UserID).Update("status", models.StatusOnline)
			
			// Notify others about user joining
			h.broadcastUserStatus(client.UserID, models.StatusOnline)
			
			log.Printf("Client %s connected", client.UserID)

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
				
				// Update user status to offline
				h.db.Model(&models.User{}).Where("id = ?", client.UserID).Update("status", models.StatusOffline)
				
				// Notify others about user leaving
				h.broadcastUserStatus(client.UserID, models.StatusOffline)
				
				log.Printf("Client %s disconnected", client.UserID)
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			h.mutex.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client.UserID)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

func (h *Hub) SendToUser(userID uuid.UUID, message []byte) {
	log.Printf("Attempting to send message to user %s", userID)
	
	h.mutex.RLock()
	client, exists := h.clients[userID]
	h.mutex.RUnlock()
	
	if exists {
		log.Printf("User %s found, sending message", userID)
		select {
		case client.Send <- message:
			log.Printf("Message sent successfully to user %s", userID)
		default:
			log.Printf("Failed to send message to user %s, closing connection", userID)
			close(client.Send)
			h.mutex.Lock()
			delete(h.clients, userID)
			h.mutex.Unlock()
		}
	} else {
		log.Printf("User %s not found in connected clients", userID)
	}
}

func (h *Hub) broadcastUserStatus(userID uuid.UUID, status models.UserStatus) {
	message := Message{
		Type:      MessageTypeUserStatus,
		UserID:    userID,
		Timestamp: getCurrentTimestamp(),
		Data: map[string]interface{}{
			"user_id": userID,
			"status":  status,
		},
	}
	
	data, _ := json.Marshal(message)
	h.broadcast <- data
}

func (h *Hub) IsUserOnline(userID uuid.UUID) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	_, exists := h.clients[userID]
	return exists
}

func (h *Hub) GetOnlineUsers() []uuid.UUID {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	users := make([]uuid.UUID, 0, len(h.clients))
	for userID := range h.clients {
		users = append(users, userID)
	}
	return users
}

func (h *Hub) HandleWebSocket(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			return
		}

		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}

		client := &Client{
			ID:     uuid.New(),
			Hub:    h,
			Conn:   conn,
			Send:   make(chan []byte, 256),
			UserID: claims.UserID,
		}

		client.Hub.register <- client

		go client.writePump()
		go client.readPump()
	}
}

func getCurrentTimestamp() int64 {
	return 0 // Implement proper timestamp
}

