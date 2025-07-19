package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"messenger/internal/websocket"
)

func main() {
	// Создаем простой hub без базы данных
	hub := websocket.NewHub(nil)
	go hub.Run()

	// Настраиваем Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS middleware
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

	// WebSocket endpoint (без аутентификации для тестирования)
	router.GET("/ws", func(c *gin.Context) {
		// Создаем простой токен для тестирования
		testUserID := uuid.New()
		
		conn, err := websocket.Upgrade(c.Writer, c.Request)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}

		client := &websocket.Client{
			ID:     uuid.New(),
			Hub:    hub,
			Conn:   conn,
			Send:   make(chan []byte, 256),
			UserID: testUserID,
		}

		hub.Register <- client

		go client.WritePump()
		go client.ReadPump()
	})

	// Статические файлы
	router.StaticFile("/", "./web/index.html")
	router.StaticFile("/index.html", "./web/index.html")
	router.StaticFile("/chat.html", "./web/chat.html")
	router.StaticFile("/websocket-test.html", "./web/websocket-test.html")
	router.Static("/static", "./web/static")

	// Запускаем сервер
	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Test server starting on port 8080")
	log.Printf("WebSocket endpoint: ws://localhost:8080/ws")
	log.Printf("Test page: http://localhost:8080/websocket-test.html")
	
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed to start:", err)
	}
}