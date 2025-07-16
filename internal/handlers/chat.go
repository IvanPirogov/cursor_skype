package handlers

import (
	"net/http"
	"messenger/pkg/models"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	// TODO: Add chat service dependency
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

func (h *ChatHandler) GetChats(c *gin.Context) {
	// TODO: Implement get user chats
	c.JSON(http.StatusOK, gin.H{"chats": []models.Chat{}})
}

func (h *ChatHandler) CreateChat(c *gin.Context) {
	// TODO: Implement create chat
	c.JSON(http.StatusCreated, gin.H{"message": "Chat created"})
}

func (h *ChatHandler) GetChat(c *gin.Context) {
	// TODO: Implement get chat by ID
	c.JSON(http.StatusOK, gin.H{"chat": models.Chat{}})
}

func (h *ChatHandler) UpdateChat(c *gin.Context) {
	// TODO: Implement update chat
	c.JSON(http.StatusOK, gin.H{"message": "Chat updated"})
}

func (h *ChatHandler) DeleteChat(c *gin.Context) {
	// TODO: Implement delete chat
	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted"})
}

func (h *ChatHandler) AddChatMember(c *gin.Context) {
	// TODO: Implement add chat member
	c.JSON(http.StatusOK, gin.H{"message": "Member added"})
}

func (h *ChatHandler) RemoveChatMember(c *gin.Context) {
	// TODO: Implement remove chat member
	c.JSON(http.StatusOK, gin.H{"message": "Member removed"})
}