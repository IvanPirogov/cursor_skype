package router

import (
	"messenger/internal/auth"
	"messenger/internal/config"
	"messenger/internal/db"
	"messenger/internal/handlers"
	"messenger/internal/middleware"
	"messenger/internal/websocket"
	"github.com/gin-gonic/gin"
)

func Setup(authService *auth.Service, hub *websocket.Hub, cfg *config.Config, database *db.Database) *gin.Engine {
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

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(database)
	chatHandler := handlers.NewChatHandler(database, hub)
	messageHandler := handlers.NewMessageHandler(database)
	contactHandler := handlers.NewContactHandler()
	callHandler := handlers.NewCallHandler()
	uploadHandler := handlers.NewUploadHandler()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Static files
	router.Static("/static", "./web/static")

	// Web pages
	router.StaticFile("/", "./web/index.html")
	router.StaticFile("/index.html", "./web/index.html")
	router.StaticFile("/chat.html", "./web/chat.html")

	// WebSocket endpoint
	router.GET("/ws", hub.HandleWebSocket(authService))

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", middleware.AuthMiddleware(authService), authHandler.Logout)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.GetMe)
				users.GET("/", userHandler.GetUsers)
				users.GET("/:id", userHandler.GetUser)
				users.PUT("/me", userHandler.UpdateMe)
				users.PUT("/status", userHandler.UpdateStatus)
			}

			// Chat routes
			chats := protected.Group("/chats")
			{
				chats.GET("/", chatHandler.GetChats)
				chats.POST("/", chatHandler.CreateChat)
				chats.GET("/:id", chatHandler.GetChat)
				chats.PUT("/:id", chatHandler.UpdateChat)
				chats.DELETE("/:id", chatHandler.DeleteChat)
				chats.POST("/:id/members", chatHandler.AddChatMember)
				chats.DELETE("/:id/members/:user_id", chatHandler.RemoveChatMember)
			}

			// Message routes
			messages := protected.Group("/messages")
			{
				messages.GET("/", messageHandler.GetMessages)
				messages.POST("/", messageHandler.SendMessage)
				messages.GET("/:id", messageHandler.GetMessage)
				messages.PUT("/:id", messageHandler.UpdateMessage)
				messages.DELETE("/:id", messageHandler.DeleteMessage)
				messages.POST("/:id/read", messageHandler.MarkMessageAsRead)
			}

			// Contact routes
			contacts := protected.Group("/contacts")
			{
				contacts.GET("/", contactHandler.GetContacts)
				contacts.POST("/", contactHandler.AddContact)
				contacts.DELETE("/:id", contactHandler.RemoveContact)
				contacts.PUT("/:id/block", contactHandler.BlockContact)
				contacts.PUT("/:id/unblock", contactHandler.UnblockContact)
			}

			// Call routes
			calls := protected.Group("/calls")
			{
				calls.GET("/", callHandler.GetCalls)
				calls.POST("/", callHandler.InitiateCall)
				calls.PUT("/:id/answer", callHandler.AnswerCall)
				calls.PUT("/:id/reject", callHandler.RejectCall)
				calls.PUT("/:id/end", callHandler.EndCall)
			}

			// File upload
			protected.POST("/upload", uploadHandler.UploadFile)
		}
	}

	return router
}