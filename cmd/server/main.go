package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/google/logger"
	"github.com/gorilla/websocket"
	"github.com/myfatemi04/trie/pkg/trie"
)

var port = flag.Uint("port", 80, "port to listen on")

func main() {
	output, err := os.OpenFile("output.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		// this doesn't make sense to use logger.Fatalln
		println("FATAL: opening output.log")
		return
	}

	logger.Init("trie", true, false, output)

	flag.Parse()

	if *port > 65535 {
		logger.Fatalf("port %d is invalid, please choose a port from 0 - 65535.", port)
		return
	}

	server := trie.NewServer()

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

	logger.Infof("server starting on port %d", *port)

	addr := fmt.Sprintf(":%d", *port)
	http.ListenAndServe(addr, nil)
}
