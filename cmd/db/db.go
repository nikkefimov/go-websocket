package db

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

// Redis client.
var redisClient *redis.Client
var ctx = context.Background()

// Initialize redis.
func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis adress.
		Password: "",               // No pw.
		DB:       0,                // Default db.
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Database is not connected", err)
	}
	log.Println("Database - OK")
}
