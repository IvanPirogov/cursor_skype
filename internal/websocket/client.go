package websocket

import (
	"bytes"
	"encoding/json"
	"log"
	"time"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		
		// Handle incoming message
		c.handleIncomingMessage(message)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleIncomingMessage(message []byte) {
	log.Printf("Received message from client %s: %s", c.UserID, string(message))
	
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	// Set the user ID from the client
	msg.UserID = c.UserID
	msg.Timestamp = time.Now().Unix()

	log.Printf("Processing message type: %s from user: %s", msg.Type, c.UserID)

	switch msg.Type {
	case MessageTypeChat:
		c.handleChatMessage(&msg)
	case MessageTypeTyping:
		c.handleTypingMessage(&msg)
	case MessageTypeCallOffer:
		c.handleCallOffer(&msg)
	case MessageTypeCallAnswer:
		c.handleCallAnswer(&msg)
	case MessageTypeCallReject:
		c.handleCallReject(&msg)
	case MessageTypeCallEnd:
		c.handleCallEnd(&msg)
	case MessageTypeMessageRead:
		c.handleMessageRead(&msg)
	case MessageTypeNewContact:
		c.handleNewContact(&msg)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

func (c *Client) handleChatMessage(msg *Message) {
	// Send to chat members only
	if chatData, ok := msg.Data.(map[string]interface{}); ok {
		if chatID, ok := chatData["chat_id"].(string); ok {
			data, _ := json.Marshal(msg)
			c.Hub.SendToChatMembers(chatID, data, c.UserID)
		}
	}
}

func (c *Client) handleTypingMessage(msg *Message) {
	// Send typing indicator to chat members only
	if chatData, ok := msg.Data.(map[string]interface{}); ok {
		if chatID, ok := chatData["chat_id"].(string); ok {
			data, _ := json.Marshal(msg)
			c.Hub.SendToChatMembers(chatID, data, c.UserID)
		}
	}
}

func (c *Client) handleCallOffer(msg *Message) {
	// Handle call offer - send to specific user
	if callData, ok := msg.Data.(map[string]interface{}); ok {
		if _, ok := callData["target_user_id"].(string); ok {
			// Send to specific user
			data, _ := json.Marshal(msg)
			// Implementation depends on how you want to handle direct messages
			c.Hub.broadcast <- data
		}
	}
}

func (c *Client) handleCallAnswer(msg *Message) {
	// Handle call answer
	data, _ := json.Marshal(msg)
	c.Hub.broadcast <- data
}

func (c *Client) handleCallReject(msg *Message) {
	// Handle call rejection
	data, _ := json.Marshal(msg)
	c.Hub.broadcast <- data
}

func (c *Client) handleCallEnd(msg *Message) {
	// Handle call end
	data, _ := json.Marshal(msg)
	c.Hub.broadcast <- data
}

func (c *Client) handleMessageRead(msg *Message) {
	// Handle message read receipt
	data, _ := json.Marshal(msg)
	c.Hub.broadcast <- data
}

func (c *Client) handleNewContact(msg *Message) {
	// Handle new contact notification
	// This message is already sent to the specific user, so we just log it
	// log.Printf("New contact notification sent to user %s", c.UserID)
}