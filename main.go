package main

import (
	"net/http"
	"os"

	"github.com/google/logger"
	"github.com/gorilla/websocket"
)

func main() {
	output, err := os.OpenFile("output.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		// this doesn't make sense to use logger.Fatalln
		println("FATAL: opening output.log")
		return
	}

	logger.Init("trie", true, false, output)

	// port := os.Getenv("PORT")
	// if port == "" {
	// 	println("FATAL: $PORT must be set")
	// 	return
	// }

	port := "3333"

	server := NewServer()

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("received connection")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Errorf("error upgrading connection: %v", err)
			return
		}

		server.HandleClient(conn)
	})

	addr := ":" + port
	logger.Infof("server starting on %s", addr)
	http.ListenAndServe(addr, nil)
}
