package main

import (
	"net/http"
	"os"

	"github.com/google/logger"
	"github.com/gorilla/websocket"
)

func main() {
	// Initialize logger
	output, err := os.OpenFile("output.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		// this doesn't make sense to use logger.Fatalln
		println("FATAL: opening output.log")
		return
	}

	logger.Init("trie", true, false, output)

	port := os.Getenv("PORT")
	if port == "" {
		println("FATAL: $PORT must be set")
		return
	}

	server := NewServer()

	// This server accepts connections via websocket.
	// Websockets follow a similar request cycle to HTTP requests, but are
	// "upgraded" to Websocket connections. Websockets allow real-time
	// bidirectional communication.
	upgrader := websocket.Upgrader{
		// Allow connections from any origin.
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("received connection")

		// Upgrade the connection to a websocket connection, allowing multiway communication.
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Errorf("error upgrading connection: %v", err)
			return
		}

		server.HandleClient(conn)
	})

	// Start the server.
	addr := ":" + port
	logger.Infof("server starting on %s", addr)
	http.ListenAndServe(addr, nil)
}
