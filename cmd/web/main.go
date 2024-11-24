package main

import (
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

// Global variables.
var redisClient *redis.Client
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var clients = make(map[*websocket.Conn]bool) // To keep track of connected clients.
var broadcast = make(chan Message)           // Channel for broadcasting messages.

// Message structure for broadcast
type Message struct {
	UserID    string `json:"user_id"`
	MessageID string `json:"message_id"`
	Message   string `json:"message"`
}

func main() {
	// Initialize database.
	initRedis()

	// Start broadcsting goroutine.
	go handleBroadcast()

	// Serve static pages.
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// Websocket route handler.
	http.HandleFunc("/ws", HandleWebSocket)

	// Launch the HTTP server.
	log.Println("Starting WebSocket server on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
