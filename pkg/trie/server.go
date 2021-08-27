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

const WS_MESSAGE_TYPE_TEXT = 1

func (s *Server) HandleClient(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Errorf("Error reading message: %v", err)
			return
		}

		logger.Infof("Received message: %s", message)

		response, err := s.trieDispatcher.DispatchRaw(message)

		if err != nil {
			conn.WriteMessage(WS_MESSAGE_TYPE_TEXT, []byte("e"+err.Error()))
			logger.Errorf("Error dispatching message: %v", err)
		} else {
			conn.WriteMessage(WS_MESSAGE_TYPE_TEXT, []byte("s"+string(response)))
		}
	}
}
