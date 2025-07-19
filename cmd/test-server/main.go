package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"messenger/internal/websocket"
)

func main() {
	hub := websocket.NewHub(nil)
	go hub.Run()

	router := gin.Default()
	
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.GET("/ws", func(c *gin.Context) {
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

	router.StaticFile("/websocket-test.html", "./web/websocket-test.html")
	router.Static("/static", "./web/static")

	log.Printf("Test server starting on port 8080")
	log.Printf("Test page: http://localhost:8080/websocket-test.html")
	
	router.Run(":8080")
}