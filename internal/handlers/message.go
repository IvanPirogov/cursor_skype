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

type MessageHandler struct {
	db *db.Database
}

func NewMessageHandler(database *db.Database) *MessageHandler {
	return &MessageHandler{
		db: database,
	}
}

// GetMessages возвращает сообщения для чата или пользователя
func (h *MessageHandler) GetMessages(c *gin.Context) {
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

	// Получаем параметры запроса
	chatIDStr := c.Query("chat_id")
	receiverIDStr := c.Query("receiver_id")
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	var messages []models.Message
	query := h.db.DB.Preload("Sender").
		Preload("Receiver").
		Preload("Chat").
		Preload("ReplyTo").
		Preload("Files").
		Preload("Reactions.User").
		Order("created_at ASC")

	// Фильтруем по чату или получателю
	if chatIDStr != "" {
		chatID, err := uuid.Parse(chatIDStr)
		if err != nil {
			// Если chat_id невалидный (например, временный), возвращаем пустой список сообщений
			c.JSON(http.StatusOK, gin.H{"messages": []models.Message{}})
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
			// Продолжаем загрузку сообщений
		} else {
			// Для приватных и групповых чатов проверяем членство
			var member models.ChatMember
			err = h.db.DB.Where("chat_id = ? AND user_id = ? AND is_active = ?", chatID, userUUID, true).
				First(&member).Error
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this chat"})
				return
			}
		}

		query = query.Where("chat_id = ?", chatID)
	} else if receiverIDStr != "" {
		receiverID, err := uuid.Parse(receiverIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receiver ID"})
			return
		}

		// Получаем сообщения между двумя пользователями
		query = query.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userUUID, receiverID, receiverID, userUUID)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either chat_id or receiver_id must be provided"})
		return
	}

	err = query.Limit(limit).Offset(offset).Find(&messages).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// SendMessage отправляет новое сообщение
func (h *MessageHandler) SendMessage(c *gin.Context) {
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
		Content    string       `json:"content" binding:"required"`
		Type       models.MessageType `json:"type"`
		ChatID     *uuid.UUID   `json:"chat_id"`
		ReceiverID *uuid.UUID   `json:"receiver_id"`
		ReplyToID  *uuid.UUID   `json:"reply_to_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Проверяем, что указан либо chat_id, либо receiver_id
	if request.ChatID == nil && request.ReceiverID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either chat_id or receiver_id must be provided"})
		return
	}

	// Если указан chat_id, проверяем права доступа
	if request.ChatID != nil {
		// Сначала получаем информацию о чате
		var chat models.Chat
		err := h.db.DB.Where("id = ? AND is_active = ?", *request.ChatID, true).
			First(&chat).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat"})
			}
			return
		}

		// Для публичных чатов отправка сообщений разрешена всем
		if chat.Type == models.ChatTypePublic {
			// Продолжаем отправку сообщения
		} else {
			// Для приватных и групповых чатов проверяем членство
			var member models.ChatMember
			err = h.db.DB.Where("chat_id = ? AND user_id = ? AND is_active = ?", *request.ChatID, userUUID, true).
				First(&member).Error
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this chat"})
				return
			}
		}
	}

	// Если указан reply_to_id, проверяем существование сообщения
	if request.ReplyToID != nil {
		var replyMessage models.Message
		err := h.db.DB.Where("id = ?", *request.ReplyToID).First(&replyMessage).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Reply message not found"})
			return
		}
	}

	// Создаем сообщение
	message := models.Message{
		SenderID:   userUUID,
		ReceiverID: request.ReceiverID,
		ChatID:     request.ChatID,
		Content:    request.Content,
		Type:       request.Type,
		Status:     models.MessageStatusSent,
		ReplyToID:  request.ReplyToID,
	}

	if message.Type == "" {
		message.Type = models.MessageTypeText
	}

	err := h.db.DB.Create(&message).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	// Загружаем созданное сообщение с полной информацией
	err = h.db.DB.Preload("Sender").
		Preload("Receiver").
		Preload("Chat").
		Preload("ReplyTo").
		Preload("Files").
		Preload("Reactions.User").
		First(&message, message.ID).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created message"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": message})
}

// GetMessage возвращает конкретное сообщение по ID
func (h *MessageHandler) GetMessage(c *gin.Context) {
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

	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var message models.Message
	err = h.db.DB.Preload("Sender").
		Preload("Receiver").
		Preload("Chat").
		Preload("ReplyTo").
		Preload("Files").
		Preload("Reactions.User").
		Where("id = ?", messageID).
		First(&message).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch message"})
		}
		return
	}

	// Проверяем права доступа к сообщению
	if message.ChatID != nil {
		// Сообщение в чате - проверяем членство
		var member models.ChatMember
		err = h.db.DB.Where("chat_id = ? AND user_id = ? AND is_active = ?", *message.ChatID, userUUID, true).
			First(&member).Error
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this message"})
			return
		}
	} else if message.ReceiverID != nil {
		// Приватное сообщение - проверяем, что пользователь отправитель или получатель
		if message.SenderID != userUUID && *message.ReceiverID != userUUID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this message"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

// UpdateMessage обновляет сообщение
func (h *MessageHandler) UpdateMessage(c *gin.Context) {
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

	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var request struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Получаем сообщение
	var message models.Message
	err = h.db.DB.Where("id = ?", messageID).First(&message).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch message"})
		}
		return
	}

	// Проверяем, что пользователь является отправителем
	if message.SenderID != userUUID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only message sender can edit the message"})
		return
	}

	// Обновляем сообщение
	err = h.db.DB.Model(&message).Updates(map[string]interface{}{
		"content":   request.Content,
		"is_edited": true,
	}).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message"})
		return
	}

	// Загружаем обновленное сообщение
	err = h.db.DB.Preload("Sender").
		Preload("Receiver").
		Preload("Chat").
		Preload("ReplyTo").
		Preload("Files").
		Preload("Reactions.User").
		First(&message, messageID).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

// DeleteMessage удаляет сообщение
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
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

	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	// Получаем сообщение
	var message models.Message
	err = h.db.DB.Where("id = ?", messageID).First(&message).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch message"})
		}
		return
	}

	// Проверяем права на удаление
	if message.ChatID != nil {
		// Сообщение в чате - проверяем роль пользователя
		var member models.ChatMember
		err = h.db.DB.Where("chat_id = ? AND user_id = ? AND is_active = ?", *message.ChatID, userUUID, true).
			First(&member).Error
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this message"})
			return
		}

		// Только отправитель или админ может удалить сообщение
		if message.SenderID != userUUID && member.Role != models.ChatMemberRoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only message sender or admin can delete the message"})
			return
		}
	} else {
		// Приватное сообщение - только отправитель может удалить
		if message.SenderID != userUUID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only message sender can delete the message"})
			return
		}
	}

	// Удаляем сообщение (мягкое удаление)
	err = h.db.DB.Delete(&message).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}

// MarkMessageAsRead отмечает сообщение как прочитанное
func (h *MessageHandler) MarkMessageAsRead(c *gin.Context) {
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

	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	// Получаем сообщение
	var message models.Message
	err = h.db.DB.Where("id = ?", messageID).First(&message).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch message"})
		}
		return
	}

	// Проверяем права доступа к сообщению
	if message.ChatID != nil {
		// Сообщение в чате - проверяем членство
		var member models.ChatMember
		err = h.db.DB.Where("chat_id = ? AND user_id = ? AND is_active = ?", *message.ChatID, userUUID, true).
			First(&member).Error
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this message"})
			return
		}
	} else if message.ReceiverID != nil {
		// Приватное сообщение - проверяем, что пользователь получатель
		if *message.ReceiverID != userUUID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this message"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message type"})
		return
	}

	// Проверяем, не отмечено ли уже сообщение как прочитанное
	var existingRead models.MessageRead
	err = h.db.DB.Where("message_id = ? AND user_id = ?", messageID, userUUID).
		First(&existingRead).Error

	if err == nil {
		// Уже отмечено как прочитанное
		c.JSON(http.StatusOK, gin.H{"message": "Message already marked as read"})
		return
	}

	// Создаем запись о прочтении
	messageRead := models.MessageRead{
		MessageID: messageID,
		UserID:    userUUID,
	}

	err = h.db.DB.Create(&messageRead).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark message as read"})
		return
	}

	// Обновляем статус сообщения
	err = h.db.DB.Model(&message).Update("status", models.MessageStatusRead).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message marked as read successfully"})
}

// GetPrivateHistory возвращает историю приватного чата по имени
func (h *MessageHandler) GetPrivateHistory(c *gin.Context) {
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
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chat name required"})
		return
	}
	// Ищем приватный чат с таким именем
	var chat models.Chat
	err := h.db.DB.Where("name = ? AND type = ? AND is_active = ?", name, models.ChatTypePrivate, true).First(&chat).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"chat": nil, "messages": []models.Message{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat"})
		return
	}
	// Проверяем, что пользователь — участник чата
	var member models.ChatMember
	err = h.db.DB.Where("chat_id = ? AND user_id = ? AND is_active = ?", chat.ID, userUUID, true).First(&member).Error
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this chat"})
		return
	}
	// Загружаем сообщения
	var messages []models.Message
	h.db.DB.Preload("Sender").Preload("Files").Where("chat_id = ?", chat.ID).Order("created_at ASC").Find(&messages)
	c.JSON(http.StatusOK, gin.H{"chat": chat, "messages": messages})
}