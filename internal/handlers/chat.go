package handlers

import (
	"net/http"
	"strconv"
	"messenger/pkg/models"
	"messenger/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatHandler struct {
	db *db.Database
}

func NewChatHandler(database *db.Database) *ChatHandler {
	return &ChatHandler{
		db: database,
	}
}

// GetChats возвращает список чатов для текущего пользователя
func (h *ChatHandler) GetChats(c *gin.Context) {
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

	// Получаем чаты, где пользователь является участником
	var chatMembers []models.ChatMember
	err := h.db.DB.Preload("Chat.Creator").
		Preload("Chat.Members.User").
		Preload("Chat.Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Where("user_id = ? AND is_active = ?", userUUID, true).
		Find(&chatMembers).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user chats"})
		return
	}

	// Создаем map для отслеживания уже добавленных чатов
	chatMap := make(map[uuid.UUID]models.Chat)
	
	// Добавляем чаты, где пользователь является участником
	for _, member := range chatMembers {
		chatMap[member.Chat.ID] = member.Chat
	}

	// Получаем все публичные каналы
	var publicChannels []models.Chat
	err = h.db.DB.Preload("Creator").
		Preload("Members.User").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Where("type = ? AND is_active = ?", models.ChatTypePublic, true).
		Find(&publicChannels).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch public channels"})
		return
	}

	// Добавляем публичные каналы (если они еще не добавлены)
	for _, channel := range publicChannels {
		if _, exists := chatMap[channel.ID]; !exists {
			chatMap[channel.ID] = channel
		}
	}

	// Преобразуем map в slice
	chats := make([]models.Chat, 0, len(chatMap))
	for _, chat := range chatMap {
		chats = append(chats, chat)
	}

	c.JSON(http.StatusOK, gin.H{"chats": chats})
}

// CreateChat создает новый чат
func (h *ChatHandler) CreateChat(c *gin.Context) {
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

	var request struct {
		Name        string      `json:"name" binding:"required"`
		Description string      `json:"description"`
		Type        models.ChatType `json:"type"`
		MemberIDs   []uuid.UUID `json:"member_ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Устанавливаем тип по умолчанию как public, если не указан
	chatType := request.Type
	if chatType == "" {
		chatType = models.ChatTypePublic
	}

	// Создаем чат
	chat := models.Chat{
		Name:        request.Name,
		Description: request.Description,
		Type:        chatType,
		CreatedBy:   userUUID,
		IsActive:    true,
	}

	err := h.db.DB.Create(&chat).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	// Добавляем создателя как участника с ролью админа
	creatorMember := models.ChatMember{
		ChatID:   chat.ID,
		UserID:   userUUID,
		Role:     models.ChatMemberRoleAdmin,
		IsActive: true,
	}

	err = h.db.DB.Create(&creatorMember).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add creator to chat"})
		return
	}

	// Добавляем других участников только для приватных и групповых чатов
	if chatType != models.ChatTypePublic {
		for _, memberID := range request.MemberIDs {
			if memberID != userUUID {
				member := models.ChatMember{
					ChatID:   chat.ID,
					UserID:   memberID,
					Role:     models.ChatMemberRoleMember,
					IsActive: true,
				}
				h.db.DB.Create(&member)
			}
		}
	}

	// Загружаем созданный чат с участниками
	err = h.db.DB.Preload("Creator").
		Preload("Members.User").
		First(&chat, chat.ID).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created chat"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"chat": chat})
}

// GetChat возвращает информацию о конкретном чате
func (h *ChatHandler) GetChat(c *gin.Context) {
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

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	// Сначала получаем информацию о чате
	var chat models.Chat
	err = h.db.DB.Where("id = ? AND is_active = ?", chatID, true).
		First(&chat).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat"})
		}
		return
	}

	// Для публичных чатов доступ разрешен всем
	if chat.Type == models.ChatTypePublic {
		// Продолжаем загрузку чата
	} else {
		// Для приватных и групповых чатов проверяем членство
		var member models.ChatMember
		err = h.db.DB.Where("chat_id = ? AND user_id = ? AND is_active = ?", chatID, userUUID, true).
			First(&member).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this chat"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check chat access"})
			}
			return
		}
	}

	// Загружаем чат с полной информацией
	err = h.db.DB.Preload("Creator").
		Preload("Members.User").
		Preload("Messages.Sender").
		Preload("Messages.Reactions").
		Where("id = ? AND is_active = ?", chatID, true).
		First(&chat).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"chat": chat})
}

// UpdateChat обновляет информацию о чате
func (h *ChatHandler) UpdateChat(c *gin.Context) {
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

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	var request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Avatar      string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Сначала получаем информацию о чате
	var chat models.Chat
	err = h.db.DB.Where("id = ? AND is_active = ?", chatID, true).
		First(&chat).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat"})
		}
		return
	}

	// Для публичных чатов только создатель может редактировать
	if chat.Type == models.ChatTypePublic {
		if chat.CreatedBy != userUUID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only chat creator can update public chat"})
			return
		}
	} else {
		// Для приватных и групповых чатов проверяем права админа
		var member models.ChatMember
		err = h.db.DB.Where("chat_id = ? AND user_id = ? AND is_active = ?", chatID, userUUID, true).
			First(&member).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this chat"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check chat access"})
			}
			return
		}

		// Проверяем роль пользователя
		if member.Role != models.ChatMemberRoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can update chat"})
			return
		}
	}

	// Обновляем чат
	updates := make(map[string]interface{})
	if request.Name != "" {
		updates["name"] = request.Name
	}
	if request.Description != "" {
		updates["description"] = request.Description
	}
	if request.Avatar != "" {
		updates["avatar"] = request.Avatar
	}

	err = h.db.DB.Model(&models.Chat{}).Where("id = ?", chatID).Updates(updates).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update chat"})
		return
	}

	// Загружаем обновленный чат
	err = h.db.DB.Preload("Creator").
		Preload("Members.User").
		First(&chat, chatID).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chat": chat})
}

// DeleteChat удаляет чат (деактивирует)
func (h *ChatHandler) DeleteChat(c *gin.Context) {
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

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	// Проверяем права доступа (только создатель может удалить чат)
	var chat models.Chat
	err = h.db.DB.Where("id = ? AND created_by = ? AND is_active = ?", chatID, userUUID, true).
		First(&chat).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only chat creator can delete the chat"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check chat access"})
		}
		return
	}

	// Деактивируем чат
	err = h.db.DB.Model(&models.Chat{}).Where("id = ?", chatID).Update("is_active", false).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chat"})
		return
	}

	// Деактивируем всех участников
	err = h.db.DB.Model(&models.ChatMember{}).Where("chat_id = ?", chatID).Update("is_active", false).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate chat members"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted successfully"})
}

// AddChatMember добавляет нового участника в чат
func (h *ChatHandler) AddChatMember(c *gin.Context) {
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

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	var request struct {
		UserID uuid.UUID `json:"user_id" binding:"required"`
		Role   models.ChatMemberRole `json:"role"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Проверяем права доступа (только админ может добавлять участников)
	var member models.ChatMember
	err = h.db.DB.Where("chat_id = ? AND user_id = ? AND is_active = ?", chatID, userUUID, true).
		First(&member).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this chat"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check chat access"})
		}
		return
	}

	if member.Role != models.ChatMemberRoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can add members"})
		return
	}

	// Проверяем, не является ли пользователь уже участником
	var existingMember models.ChatMember
	err = h.db.DB.Where("chat_id = ? AND user_id = ?", chatID, request.UserID).
		First(&existingMember).Error

	if err == nil {
		// Пользователь уже участник, активируем его
		err = h.db.DB.Model(&existingMember).Updates(map[string]interface{}{
			"is_active": true,
			"role":      request.Role,
			"left_at":   nil,
		}).Error
	} else if err == gorm.ErrRecordNotFound {
		// Создаем нового участника
		newMember := models.ChatMember{
			ChatID:   chatID,
			UserID:   request.UserID,
			Role:     request.Role,
			IsActive: true,
		}
		err = h.db.DB.Create(&newMember).Error
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing member"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member to chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added successfully"})
}

// RemoveChatMember удаляет участника из чата
func (h *ChatHandler) RemoveChatMember(c *gin.Context) {
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

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	memberID, err := uuid.Parse(c.Param("member_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
		return
	}

	// Проверяем права доступа (только админ может удалять участников)
	var adminMember models.ChatMember
	err = h.db.DB.Where("chat_id = ? AND user_id = ? AND is_active = ?", chatID, userUUID, true).
		First(&adminMember).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this chat"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check chat access"})
		}
		return
	}

	if adminMember.Role != models.ChatMemberRoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can remove members"})
		return
	}

	// Проверяем, что удаляемый участник не является создателем чата
	var chat models.Chat
	err = h.db.DB.Where("id = ?", chatID).First(&chat).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat"})
		return
	}

	if chat.CreatedBy == memberID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot remove chat creator"})
		return
	}

	// Удаляем участника (деактивируем)
	err = h.db.DB.Model(&models.ChatMember{}).
		Where("chat_id = ? AND user_id = ?", chatID, memberID).
		Updates(map[string]interface{}{
			"is_active": false,
			"left_at":   gorm.Expr("NOW()"),
		}).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove member from chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed successfully"})
}