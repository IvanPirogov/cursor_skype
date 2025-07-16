package handlers

import (
	"net/http"
	"messenger/pkg/models"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	// TODO: Add message service dependency
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	// TODO: Implement get messages
	c.JSON(http.StatusOK, gin.H{"messages": []models.Message{}})
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	// TODO: Implement send message
	c.JSON(http.StatusCreated, gin.H{"message": "Message sent"})
}

func (h *MessageHandler) GetMessage(c *gin.Context) {
	// TODO: Implement get message by ID
	c.JSON(http.StatusOK, gin.H{"message": models.Message{}})
}

func (h *MessageHandler) UpdateMessage(c *gin.Context) {
	// TODO: Implement update message
	c.JSON(http.StatusOK, gin.H{"message": "Message updated"})
}

func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	// TODO: Implement delete message
	c.JSON(http.StatusOK, gin.H{"message": "Message deleted"})
}

func (h *MessageHandler) MarkMessageAsRead(c *gin.Context) {
	// TODO: Implement mark message as read
	c.JSON(http.StatusOK, gin.H{"message": "Message marked as read"})
}