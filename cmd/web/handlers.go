package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// PushMessageToRedis pushes the message to redis with users and messages IDs.
func PushMessageToRedis(userID, messageID, message string) error {
	ctx := context.Background()
	// Store the message in Redis under a list 'chat:messages' with format "userID:messageID:message".
	err := redisClient.LPush(ctx, "chat:messages", fmt.Sprintf("%s:%s:%s", userID, messageID, message)).Err()
	if err != nil {
		log.Println("Error pushing message to Redis:", err)
		return err
	}
	return nil
}

// HandleWebSocket handles the websocket connection, reads incoming messages, and broadcasts them.
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a websocket connection.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade WebSocket:", err)
		return
	}
	defer conn.Close()

	// Generate unique IDs for each user.
	userID := fmt.Sprintf("%d", time.Now().UnixNano()) // A simple unique user ID based on time.
	log.Printf("User connected with ID: %s", userID)

	// Add user connection to clients map, it allows to watch self messages which user sends.
	clients[conn] = true
	defer delete(clients, conn)

	// Read messages from the websocket and store them in database.
	for {
		// Read the websocket message. This will return three values:
		// 1. Message type (will ignore it by using _ in loop),
		// 2. The message content (which is the byte slice),
		// 3. Error (if any error exists).
		_, message, err := conn.ReadMessage() // This reads the message content into the "message" variable.
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Convert the byte slice(message) to a string.
		messageStr := string(message)

		// Log the message with the user ID.
		log.Printf("User %s sent message: %s", userID, messageStr)

		// Generate a unique message ID, by using timestamp and random ID.
		messageID := fmt.Sprintf("%d", time.Now().UnixNano())
		log.Printf("Message ID: %s", messageID)

		// Push the message to database (implement PushMessageToRedis).
		if err := PushMessageToRedis(userID, messageID, messageStr); err != nil {
			log.Println("Error pushing message to Redis:", err)
		}

		// Send message to the broadcast channel
		broadcast <- Message{
			UserID:    userID,
			MessageID: messageID,
			Message:   messageStr,
		}
	}
}

// Func runs own goruutine and broadcast to all connected users.
func handleBroadcast() {
	for {
		// Receive a message from the broadcast channel.
		message := <-broadcast

		// Send the message to all connected users.
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(message.Message))
			if err != nil {
				log.Printf("Error sending message to client %s: %v", message.UserID, err)
				client.Close()
				delete(clients, client) // Remove disconnected client from the list.
			}
		}
	}
}
