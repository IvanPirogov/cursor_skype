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
	h.mutex.RLock()
	client, exists := h.clients[userID]
	h.mutex.RUnlock()
	
	if exists {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			h.mutex.Lock()
			delete(h.clients, userID)
			h.mutex.Unlock()
		}
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

// SendNewContactNotification отправляет уведомление пользователю о добавлении нового контакта
func (h *Hub) SendNewContactNotification(userID uuid.UUID, contact models.Contact) {
	message := Message{
		Type:      MessageTypeNewContact,
		UserID:    userID,
		Timestamp: getCurrentTimestamp(),
		Data: map[string]interface{}{
			"contact_id": contact.ContactID,
			"contact": map[string]interface{}{
				"id":         contact.Contact.ID,
				"username":   contact.Contact.Username,
				"first_name": contact.Contact.FirstName,
				"last_name":  contact.Contact.LastName,
				"avatar":     contact.Contact.Avatar,
				"status":     contact.Contact.Status,
			},
			"nickname": contact.Nickname,
		},
	}
	
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling new contact notification: %v", err)
		return
	}
	
	h.SendToUser(userID, data)
}

// SendNewContactNotificationSimple отправляет уведомление о новом контакте с минимальными данными
func (h *Hub) SendNewContactNotificationSimple(userID uuid.UUID, contactID uuid.UUID, contactUsername string, contactFirstName string, contactLastName string, nickname string) {
	message := Message{
		Type:      MessageTypeNewContact,
		UserID:    userID,
		Timestamp: getCurrentTimestamp(),
		Data: map[string]interface{}{
			"contact_id": contactID,
			"contact": map[string]interface{}{
				"id":         contactID,
				"username":   contactUsername,
				"first_name": contactFirstName,
				"last_name":  contactLastName,
			},
			"nickname": nickname,
		},
	}
	
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling new contact notification: %v", err)
		return
	}
	
	h.SendToUser(userID, data)
}