package handlers

import (
	"net/http"
	"messenger/pkg/models"
	"messenger/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContactHandler struct {
	db *db.Database
}

func NewContactHandler(database *db.Database) *ContactHandler {
	return &ContactHandler{db: database}
}

// Получить список контактов текущего пользователя
func (h *ContactHandler) GetContacts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	var contacts []models.Contact
	err := h.db.DB.Preload("Contact").Where("user_id = ?", userUUID).Find(&contacts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contacts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"contacts": contacts})
}

// Добавить контакт по username
func (h *ContactHandler) AddContact(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	var req struct {
		Username string `json:"username" binding:"required"`
		Nickname string `json:"nickname"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Найти пользователя по username
	var contactUser models.User
	if err := h.db.DB.Where("username = ? AND is_active = ?", req.Username, true).First(&contactUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		}
		return
	}

	// Нельзя добавить себя
	if contactUser.ID == userUUID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot add yourself as a contact"})
		return
	}

	// Проверить, не добавлен ли уже контакт
	var existing models.Contact
	if err := h.db.DB.Where("user_id = ? AND contact_id = ?", userUUID, contactUser.ID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Contact already exists"})
		return
	}

	contact := models.Contact{
		UserID:    userUUID,
		ContactID: contactUser.ID,
		Nickname:  req.Nickname,
		IsBlocked: false,
	}
	if err := h.db.DB.Create(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add contact"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Contact added"})
}

// Удалить контакт
func (h *ContactHandler) RemoveContact(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	contactID := c.Param("id")
	contactUUID, err := uuid.Parse(contactID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	if err := h.db.DB.Where("user_id = ? AND contact_id = ?", userUUID, contactUUID).Delete(&models.Contact{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove contact"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact removed"})
}

// Заблокировать контакт
func (h *ContactHandler) BlockContact(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	contactID := c.Param("id")
	contactUUID, err := uuid.Parse(contactID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	if err := h.db.DB.Model(&models.Contact{}).
		Where("user_id = ? AND contact_id = ?", userUUID, contactUUID).
		Update("is_blocked", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to block contact"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact blocked"})
}

// Разблокировать контакт
func (h *ContactHandler) UnblockContact(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
		return
	}

	contactID := c.Param("id")
	contactUUID, err := uuid.Parse(contactID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	if err := h.db.DB.Model(&models.Contact{}).
		Where("user_id = ? AND contact_id = ?", userUUID, contactUUID).
		Update("is_blocked", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unblock contact"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact unblocked"})
}