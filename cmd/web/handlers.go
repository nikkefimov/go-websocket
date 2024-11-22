package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type ChatServer struct {
	// Map with users's active websocket connections,
	// each key is a websocket connection
	// and the value is a boolean indicating whether the connection is active.
	clients map[*websocket.Conn]bool
	// Broadcast channel is uset to send messages to all connected users.
	broadcast chan string
	// "mu" is a mutex used to synchronize access to the users map (thread-safety).
	mu sync.Mutex
}

// NewChatServer creates and returns a new instance of ChatServer.
func NewChatServer() *ChatServer {
	return &ChatServer{
		// Initialize the map to keep track of active users.
		clients: make(map[*websocket.Conn]bool),
		// Creating a new channel for broadcasting messages.
		broadcast: make(chan string),
	}
}

// Handler for listening messages.
// handleMessages listens for incoming messages on the broadcast channel
// and send them to all connected users.
func (cs *ChatServer) handleMessages() {
	// Loop indefinitely to handle the messages.
	for {
		// Get the next message from the broadcast channel.
		msg := <-cs.broadcast
		// Lock the users map to prevent race conditions while message is sending.
		cs.mu.Lock()
		// Loop over all users and send the message to them.
		for client := range cs.clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			// If an error here while sending the message, log it and close the connection.
			if err != nil {
				log.Printf("Error sending message: %v", err)
				client.Close()
				// Remove the user from the users map.
				delete(cs.clients, client)
			}
		}
		// Unlock the users map after sending messages.
		cs.mu.Lock()
	}
}

// Handler for incoming websocket connection from a user.
func (cs *ChatServer) handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP request to a websocket connection.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Lock the users map before modifying it.
	cs.mu.Lock()
	// Add the new conncetion to the users's map.
	cs.clients[conn] = true
	// Unlock users's map after modification.
	cs.mu.Unlock()

	// Ensure the websocket connection is closed when the function exits.
	defer func() {
		cs.mu.Lock()
		delete(cs.clients, conn)
		cs.mu.Unlock()
		conn.Close()
	}()

	// Listen messages from the user.
	for {
		// Read a message from the websocket connection.
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("User disconnected:", err)
			// If error or user disconnected than break the loop.
			break
		}

		// Broadcast the message to all connected users.
		cs.broadcast <- string(msg)

		// Save the message to redis, ensure all messages are available in the system.
		go func() {
			// Publish the message to the Redis channel "chat:messages".
			err := redisClient.Publish(context.Background(), "chat:messages", msg).Err()
			if err != nil {
				log.Printf("Error database: %v", err)
			}
		}()
	}
}
