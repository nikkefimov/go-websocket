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

func main() {
	// Initialize database.
	initRedis()

	// Serve static pages.
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// Websocket route handler.
	http.HandleFunc("/ws", HandleWebSocket)

	// Launch the HTTP server.
	log.Println("Starting WebSocket server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
