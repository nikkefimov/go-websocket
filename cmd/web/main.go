package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		//	ReadBufferSize:  4096,
		//	WriteBufferSize: 4096,
		// CheckOrigin is used to allow connections from any origin, its importatnt for testing,
		// in production should validate the origin to prevent cross-origin attacks.
		CheckOrigin: func(r *http.ReadRequest) bool {
			return true
		},
	}

	// Global Redis client variable to connect the Redis server.
	redisClient *redis.Client
)

func main() {
	// Intialize redis client to connect to the redis server at localhost:6379
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // redis server.
	})
	// Redis check. Ping the redis server to check if it is available
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		// If redis is not available, log the error and stop the app.
		log.Fatalf("Database connect issue: %v", err)
	}

	// Create a chat server.
	chatServer := NewChatServer()

	// Start a goroutine to handle incoming messages and broadcast messages.
	go chatServer.handleMessages()

	// Set up websocket route to handle connections on /ws path.
	http.HandleFunc("/ws", chatServer.handleWebSocketConnection)

	// Serve static HTML page for the frontend.
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// Start the HTTP server to listen for incoming requests on the port 8081.
	log.Println("Server started on port :8081")
	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		// If server error, log the error and stop the app.
		log.Fatal(err)
	}
}
