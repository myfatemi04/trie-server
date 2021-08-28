package main

import (
	"github.com/google/logger"
	"github.com/gorilla/websocket"
)

/*
A Server is just a wrapper around a TrieDispatcher.
This TrieDispatcher is responsible for ensuring that
the state of the trie is atomic, and that all messages
are handled individually and in order.
*/
type Server struct {
	trieDispatcher *ThreadSafeDispatcher
}

/*
NewServer creates a new server with an empty Trie.
*/
func NewServer() Server {
	return Server{NewThreadSafeDispatcher(NewTrie())}
}

const WS_MESSAGE_TYPE_TEXT = 1

/*
HandleClient handles a client connection.
This simply reads messages from the client and dispatches them to the trie.
*/
func (s *Server) HandleClient(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Errorf("Error reading message: %v", err)
			return
		}

		logger.Infof("Received message: %s", message)

		// Dispatch the message to the trie. This accepts a string
		// and returns a string as a response.
		response, err := s.trieDispatcher.DispatchRaw(message)

		if err != nil {
			conn.WriteMessage(WS_MESSAGE_TYPE_TEXT, []byte("e"+err.Error()))
			logger.Errorf("Error dispatching message: %v", err)
		} else {
			conn.WriteMessage(WS_MESSAGE_TYPE_TEXT, []byte("s"+string(response)))
		}
	}
}
