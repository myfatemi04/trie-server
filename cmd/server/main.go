package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/google/logger"
	"github.com/gorilla/websocket"
	"github.com/myfatemi04/trie/pkg/trie"
)

var port = flag.Uint("port", 3333, "port to listen on")

func main() {
	flag.Parse()

	if *port > 65535 {
		logger.Fatalf("port %d is invalid", port)
	}

	server := trie.NewServer()

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Errorf("error upgrading connection: %v", err)
			return
		}

		server.HandleClient(conn)
	})

	addr := fmt.Sprintf(":%d", *port)
	http.ListenAndServe(addr, nil)
}
