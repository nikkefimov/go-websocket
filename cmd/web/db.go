package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

// Redis initialization.
func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis address.
		DB:   0,                // Default DB.
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")
}
