package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		//	ReadBufferSize:  4096,
		//	WriteBufferSize: 4096,
		// CheckOrigin is used to allow connections from any origin, its importatnt for testing,
		// in production should validate the origin to prevent cross-origin attacks.
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
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
