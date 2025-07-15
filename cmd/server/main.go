package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"messenger/internal/auth"
	"messenger/internal/config"
	"messenger/internal/db"
	"messenger/internal/middleware"
	"messenger/internal/websocket"
	"messenger/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database
	database, err := db.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Run migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize services
	authService := auth.NewService(database.DB, &cfg.JWT)
	
	// Initialize WebSocket hub
	hub := websocket.NewHub(database.DB)
	go hub.Run()

	// Setup router
	router := setupRouter(authService, hub, cfg)

	// Start server
	server := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        router,
		ReadTimeout:    time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func setupRouter(authService *auth.Service, hub *websocket.Hub, cfg *config.Config) *gin.Engine {
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// WebSocket endpoint
	router.GET("/ws", hub.HandleWebSocket(authService))

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", handleRegister(authService))
			auth.POST("/login", handleLogin(authService))
			auth.POST("/logout", middleware.AuthMiddleware(authService), handleLogout(authService))
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", handleGetMe())
				users.GET("/", handleGetUsers())
				users.GET("/:id", handleGetUser())
				users.PUT("/me", handleUpdateMe())
				users.PUT("/status", handleUpdateStatus())
			}

			// Chat routes
			chats := protected.Group("/chats")
			{
				chats.GET("/", handleGetChats())
				chats.POST("/", handleCreateChat())
				chats.GET("/:id", handleGetChat())
				chats.PUT("/:id", handleUpdateChat())
				chats.DELETE("/:id", handleDeleteChat())
				chats.POST("/:id/members", handleAddChatMember())
				chats.DELETE("/:id/members/:user_id", handleRemoveChatMember())
			}

			// Message routes
			messages := protected.Group("/messages")
			{
				messages.GET("/", handleGetMessages())
				messages.POST("/", handleSendMessage())
				messages.GET("/:id", handleGetMessage())
				messages.PUT("/:id", handleUpdateMessage())
				messages.DELETE("/:id", handleDeleteMessage())
				messages.POST("/:id/read", handleMarkMessageAsRead())
			}

			// Contact routes
			contacts := protected.Group("/contacts")
			{
				contacts.GET("/", handleGetContacts())
				contacts.POST("/", handleAddContact())
				contacts.DELETE("/:id", handleRemoveContact())
				contacts.PUT("/:id/block", handleBlockContact())
				contacts.PUT("/:id/unblock", handleUnblockContact())
			}

			// Call routes
			calls := protected.Group("/calls")
			{
				calls.GET("/", handleGetCalls())
				calls.POST("/", handleInitiateCall())
				calls.PUT("/:id/answer", handleAnswerCall())
				calls.PUT("/:id/reject", handleRejectCall())
				calls.PUT("/:id/end", handleEndCall())
			}

			// File upload
			protected.POST("/upload", handleFileUpload())
		}
	}

	return router
}

// Auth handlers
func handleRegister(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req auth.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := authService.Register(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, response)
	}
}

func handleLogin(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req auth.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := authService.Login(req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func handleLogout(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		token := c.GetString("token")

		parsedUserID, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		if err := authService.Logout(parsedUserID, token); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}

// User handlers
func handleGetMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user_id":  c.GetString("user_id"),
			"username": c.GetString("username"),
			"email":    c.GetString("email"),
		})
	}
}

func handleGetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement get users
		c.JSON(http.StatusOK, gin.H{"users": []models.User{}})
	}
}

func handleGetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement get user by ID
		c.JSON(http.StatusOK, gin.H{"user": models.User{}})
	}
}

func handleUpdateMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement update user profile
		c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
	}
}

func handleUpdateStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement update user status
		c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
	}
}

// Chat handlers
func handleGetChats() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement get user chats
		c.JSON(http.StatusOK, gin.H{"chats": []models.Chat{}})
	}
}

func handleCreateChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement create chat
		c.JSON(http.StatusCreated, gin.H{"message": "Chat created"})
	}
}

func handleGetChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement get chat by ID
		c.JSON(http.StatusOK, gin.H{"chat": models.Chat{}})
	}
}

func handleUpdateChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement update chat
		c.JSON(http.StatusOK, gin.H{"message": "Chat updated"})
	}
}

func handleDeleteChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement delete chat
		c.JSON(http.StatusOK, gin.H{"message": "Chat deleted"})
	}
}

func handleAddChatMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement add chat member
		c.JSON(http.StatusOK, gin.H{"message": "Member added"})
	}
}

func handleRemoveChatMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement remove chat member
		c.JSON(http.StatusOK, gin.H{"message": "Member removed"})
	}
}

// Message handlers
func handleGetMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement get messages
		c.JSON(http.StatusOK, gin.H{"messages": []models.Message{}})
	}
}

func handleSendMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement send message
		c.JSON(http.StatusCreated, gin.H{"message": "Message sent"})
	}
}

func handleGetMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement get message by ID
		c.JSON(http.StatusOK, gin.H{"message": models.Message{}})
	}
}

func handleUpdateMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement update message
		c.JSON(http.StatusOK, gin.H{"message": "Message updated"})
	}
}

func handleDeleteMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement delete message
		c.JSON(http.StatusOK, gin.H{"message": "Message deleted"})
	}
}

func handleMarkMessageAsRead() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement mark message as read
		c.JSON(http.StatusOK, gin.H{"message": "Message marked as read"})
	}
}

// Contact handlers
func handleGetContacts() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement get contacts
		c.JSON(http.StatusOK, gin.H{"contacts": []models.Contact{}})
	}
}

func handleAddContact() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement add contact
		c.JSON(http.StatusCreated, gin.H{"message": "Contact added"})
	}
}

func handleRemoveContact() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement remove contact
		c.JSON(http.StatusOK, gin.H{"message": "Contact removed"})
	}
}

func handleBlockContact() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement block contact
		c.JSON(http.StatusOK, gin.H{"message": "Contact blocked"})
	}
}

func handleUnblockContact() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement unblock contact
		c.JSON(http.StatusOK, gin.H{"message": "Contact unblocked"})
	}
}

// Call handlers
func handleGetCalls() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement get calls
		c.JSON(http.StatusOK, gin.H{"calls": []models.Call{}})
	}
}

func handleInitiateCall() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement initiate call
		c.JSON(http.StatusCreated, gin.H{"message": "Call initiated"})
	}
}

func handleAnswerCall() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement answer call
		c.JSON(http.StatusOK, gin.H{"message": "Call answered"})
	}
}

func handleRejectCall() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement reject call
		c.JSON(http.StatusOK, gin.H{"message": "Call rejected"})
	}
}

func handleEndCall() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement end call
		c.JSON(http.StatusOK, gin.H{"message": "Call ended"})
	}
}

// File upload handler
func handleFileUpload() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement file upload
		c.JSON(http.StatusOK, gin.H{"message": "File uploaded"})
	}
}