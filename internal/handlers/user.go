package handlers

import (
	"net/http"
	"strconv"
	"time"
	"messenger/pkg/models"
	"messenger/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *db.Database
}

func NewUserHandler(database *db.Database) *UserHandler {
	return &UserHandler{
		db: database,
	}
}

// GetMe возвращает информацию о текущем пользователе
func (h *UserHandler) GetMe(c *gin.Context) {
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

	var user models.User
	err := h.db.DB.Where("id = ? AND is_active = ?", userUUID, true).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		}
		return
	}

	// Скрываем пароль
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetUsers возвращает список пользователей с пагинацией и поиском
func (h *UserHandler) GetUsers(c *gin.Context) {
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
	search := c.Query("search")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	status := c.Query("status")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	var users []models.User
	query := h.db.DB.Where("is_active = ?", true)

	// Исключаем текущего пользователя из результатов
	query = query.Where("id != ?", userUUID)

	// Поиск по имени пользователя, email или полному имени
	if search != "" {
		query = query.Where("username ILIKE ? OR email ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Фильтр по статусу
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err = query.Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Скрываем пароли
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// GetUser возвращает информацию о конкретном пользователе
func (h *UserHandler) GetUser(c *gin.Context) {
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

	targetUserID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	err = h.db.DB.Where("id = ? AND is_active = ?", targetUserID, true).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		}
		return
	}

	// Скрываем пароль
	user.Password = ""

	// Если пользователь запрашивает свой профиль, возвращаем полную информацию
	if userUUID == targetUserID {
		c.JSON(http.StatusOK, gin.H{"user": user})
		return
	}

	// Для других пользователей возвращаем публичную информацию
	publicUser := struct {
		ID        uuid.UUID           `json:"id"`
		Username  string              `json:"username"`
		FirstName string              `json:"first_name"`
		LastName  string              `json:"last_name"`
		Avatar    string              `json:"avatar"`
		Status    models.UserStatus   `json:"status"`
		LastSeen  *time.Time          `json:"last_seen"`
		CreatedAt time.Time           `json:"created_at"`
	}{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		Status:    user.Status,
		LastSeen:  user.LastSeen,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"user": publicUser})
}

// UpdateMe обновляет профиль текущего пользователя
func (h *UserHandler) UpdateMe(c *gin.Context) {
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
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Avatar    string `json:"avatar"`
		Email     string `json:"email" binding:"omitempty,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Получаем текущего пользователя
	var user models.User
	err := h.db.DB.Where("id = ? AND is_active = ?", userUUID, true).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		}
		return
	}

	// Проверяем уникальность email, если он изменяется
	if request.Email != "" && request.Email != user.Email {
		var existingUser models.User
		err = h.db.DB.Where("email = ? AND id != ?", request.Email, userUUID).First(&existingUser).Error
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email uniqueness"})
			return
		}
	}

	// Обновляем поля
	updates := make(map[string]interface{})
	if request.FirstName != "" {
		updates["first_name"] = request.FirstName
	}
	if request.LastName != "" {
		updates["last_name"] = request.LastName
	}
	if request.Avatar != "" {
		updates["avatar"] = request.Avatar
	}
	if request.Email != "" {
		updates["email"] = request.Email
	}

	err = h.db.DB.Model(&user).Updates(updates).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Загружаем обновленного пользователя
	err = h.db.DB.First(&user, userUUID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user"})
		return
	}

	// Скрываем пароль
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateStatus обновляет статус текущего пользователя
func (h *UserHandler) UpdateStatus(c *gin.Context) {
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
		Status models.UserStatus `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Валидация статуса
	validStatuses := []models.UserStatus{
		models.StatusOnline,
		models.StatusOffline,
		models.StatusAway,
		models.StatusBusy,
		models.StatusInvisible,
	}

	isValidStatus := false
	for _, status := range validStatuses {
		if request.Status == status {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	// Обновляем статус и время последнего действия
	now := time.Now()
	updates := map[string]interface{}{
		"status":    request.Status,
		"last_seen": &now,
	}

	err := h.db.DB.Model(&models.User{}).Where("id = ?", userUUID).Updates(updates).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	// Загружаем обновленного пользователя
	var user models.User
	err = h.db.DB.Where("id = ?", userUUID).First(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user"})
		return
	}

	// Скрываем пароль
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{"user": user})
}