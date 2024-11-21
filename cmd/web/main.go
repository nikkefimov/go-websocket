package main

import (
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
}
