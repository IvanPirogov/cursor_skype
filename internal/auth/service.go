package auth

import (
	"errors"
	"fmt"
	"time"
	"messenger/pkg/models"
	"messenger/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	db        *gorm.DB
	jwtSecret string
	jwtExpiry int
}

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type AuthResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}

func NewService(db *gorm.DB, jwtConfig *config.JWTConfig) *Service {
	return &Service{
		db:        db,
		jwtSecret: jwtConfig.Secret,
		jwtExpiry: jwtConfig.ExpiresIn,
	}
}

func (s *Service) Register(req RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Status:    models.StatusOffline,
		IsActive:  true,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create user session
	session := models.UserSession{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(s.jwtExpiry) * time.Hour),
	}

	if err := s.db.Create(&session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Clear password from response
	user.Password = ""

	return &AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *Service) Login(req LoginRequest) (*AuthResponse, error) {
	var user models.User
	if err := s.db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Update user status and last seen
	now := time.Now()
	user.Status = models.StatusOnline
	user.LastSeen = &now
	if err := s.db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user status: %w", err)
	}

	// Create user session
	session := models.UserSession{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(s.jwtExpiry) * time.Hour),
	}

	if err := s.db.Create(&session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Clear password from response
	user.Password = ""

	return &AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *Service) Logout(userID uuid.UUID, token string) error {
	// Update user status to offline
	now := time.Now()
	if err := s.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"status":    models.StatusOffline,
		"last_seen": &now,
	}).Error; err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	// Delete session
	if err := s.db.Where("user_id = ? AND token = ?", userID, token).Delete(&models.UserSession{}).Error; err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Check if session exists
	var session models.UserSession
	if err := s.db.Where("user_id = ? AND token = ? AND expires_at > ?", claims.UserID, tokenString, time.Now()).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found or expired")
		}
		return nil, fmt.Errorf("failed to find session: %w", err)
	}

	return claims, nil
}

func (s *Service) generateToken(user models.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.jwtExpiry) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}