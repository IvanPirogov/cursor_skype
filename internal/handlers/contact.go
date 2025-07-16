package handlers

import (
	"net/http"
	"messenger/pkg/models"
	"github.com/gin-gonic/gin"
)

type ContactHandler struct {
	// TODO: Add contact service dependency
}

func NewContactHandler() *ContactHandler {
	return &ContactHandler{}
}

func (h *ContactHandler) GetContacts(c *gin.Context) {
	// TODO: Implement get contacts
	c.JSON(http.StatusOK, gin.H{"contacts": []models.Contact{}})
}

func (h *ContactHandler) AddContact(c *gin.Context) {
	// TODO: Implement add contact
	c.JSON(http.StatusCreated, gin.H{"message": "Contact added"})
}

func (h *ContactHandler) RemoveContact(c *gin.Context) {
	// TODO: Implement remove contact
	c.JSON(http.StatusOK, gin.H{"message": "Contact removed"})
}

func (h *ContactHandler) BlockContact(c *gin.Context) {
	// TODO: Implement block contact
	c.JSON(http.StatusOK, gin.H{"message": "Contact blocked"})
}

func (h *ContactHandler) UnblockContact(c *gin.Context) {
	// TODO: Implement unblock contact
	c.JSON(http.StatusOK, gin.H{"message": "Contact unblocked"})
}