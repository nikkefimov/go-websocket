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
}
