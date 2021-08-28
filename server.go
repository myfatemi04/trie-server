package main

import (
	"io/ioutil"
	"net/http"

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
	// handles messages in a thread-safe fashion
	trieDispatcher *ThreadSafeDispatcher
	// HTTP server route handler
	httpServeMux *http.ServeMux
	// Websocket upgrader
	upgrader *websocket.Upgrader
}

/*
NewServer creates a new server with an empty Trie.
*/
func NewServer() *Server {
	// Allows us to handle messages in a thread-safe fashion.
	trieDispatcher := NewThreadSafeDispatcher(NewTrie())

	// Router for HTTP requests.
	// Allows us to add custom functions for routes.
	httpServeMux := http.NewServeMux()

	// This server accepts connections via websocket.
	// Websockets follow a similar request cycle to HTTP requests, but are
	// "upgraded" to Websocket connections. Websockets allow real-time
	// bidirectional communication.
	upgrader := &websocket.Upgrader{
		// Allow connections from any origin.
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	server := &Server{
		trieDispatcher: trieDispatcher,
		httpServeMux:   httpServeMux,
		upgrader:       upgrader,
	}

	// Add routes to the HTTP server.
	httpServeMux.HandleFunc("/http", server.HandleHTTP)
	httpServeMux.HandleFunc("/ws", server.HandleWS)

	return server
}

const WS_MESSAGE_TYPE_TEXT = 1

func (s *Server) HandleWS(w http.ResponseWriter, r *http.Request) {
	logger.Infof("received connection")

	// Upgrade the connection to a websocket connection, allowing multiway communication.
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Errorf("error upgrading connection: %v", err)
		return
	}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Errorf("Error reading message: %v", err)
			return
		}

		logger.Infof("Received message: %s", message)

		// Dispatch the message to the trie. This accepts a string
		// and returns a string as a response.
		response, err := s.Process(message)

		if err != nil {
			conn.WriteMessage(WS_MESSAGE_TYPE_TEXT, []byte("e"+err.Error()))
			logger.Errorf("Error handling message: %v", err)
		} else {
			conn.WriteMessage(WS_MESSAGE_TYPE_TEXT, []byte("s"+string(response)))
		}
	}
}

func (s *Server) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	// Verify that the request is a POST.
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read the request body.
	message, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("Error reading request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Infof("Received message: %s", message)

	// Execute the command.
	response, err := s.Process(message)
	if err != nil {
		logger.Errorf("Error handling message: %v", err)
		// Format: e{error message}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("e" + err.Error()))
	} else {
		// Format: s{successful response}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("s" + string(response)))
	}
}

func (s *Server) Process(message []byte) ([]byte, error) {
	return s.trieDispatcher.DispatchRaw(message)
}

func (s *Server) HttpServeMux() *http.ServeMux {
	return s.httpServeMux
}
