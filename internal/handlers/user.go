package handlers

import (
	"net/http"
	"messenger/pkg/models"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	// TODO: Add user service dependency
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user_id":  c.GetString("user_id"),
		"username": c.GetString("username"),
		"email":    c.GetString("email"),
	})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	// TODO: Implement get users
	c.JSON(http.StatusOK, gin.H{"users": []models.User{}})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	// TODO: Implement get user by ID
	c.JSON(http.StatusOK, gin.H{"user": models.User{}})
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	// TODO: Implement update user profile
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}

func (h *UserHandler) UpdateStatus(c *gin.Context) {
	// TODO: Implement update user status
	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}