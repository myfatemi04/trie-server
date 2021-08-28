package main

import (
	"net/http"
	"os"

	"github.com/google/logger"
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
		logger.Fatalln("$PORT must be set")
		return
	}

	server := NewServer()

	logger.Infof("server starting on :%s", port)
	http.ListenAndServe(":"+port, server.HttpServeMux())
}
