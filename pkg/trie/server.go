package trie

import (
	"github.com/google/logger"
	"github.com/gorilla/websocket"
)

type Server struct {
	trieDispatcher *ThreadSafeDispatcher
}

func NewServer() Server {
	return Server{NewThreadSafeDispatcher(NewTrie())}
}

func (s *Server) HandleClient(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Errorf("Error reading message: %v", err)
			return
		}

		s.trieDispatcher.DispatchRaw(message)
	}
}
