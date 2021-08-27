package trie

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/google/logger"
)

/*
ThreadSafeDispatcher dispatches commands to a Trie and ensures that they are processed in the order in which they are received.
*/
type ThreadSafeDispatcher struct {
	// trie is the Trie to which commands are dispatched
	trie *Trie
	// dispatcherMutex is the Mutex to ensure that no commands are processed out of order
	dispatcherMutex sync.Mutex
}

/*
Creates a thread-safe dispatcher for the given Trie.
*/
func NewThreadSafeDispatcher(trie *Trie) *ThreadSafeDispatcher {
	if trie == nil {
		return &ThreadSafeDispatcher{trie: NewTrie()}
	}
	return &ThreadSafeDispatcher{trie: trie}
}

const (
	CMD_INSERT = iota
	CMD_DELETE
	CMD_EXISTS
	CMD_COMPLETIONS
	CMD_KEYS
)

const ASCII_0 = 48

func (s *ThreadSafeDispatcher) DispatchRaw(message []byte) ([]byte, error) {
	command := message[0] - ASCII_0

	switch command {
	case CMD_INSERT:
		logger.Infof("inserting %s", message[1:])
		result, err := s.DispatchInsert(string(message[1:]))
		if err != nil {
			return nil, err
		}
		return json.Marshal(result)
	case CMD_DELETE:
		logger.Infof("deleting %s", message[1:])
		result, err := s.DispatchDelete(string(message[1:]))
		if err != nil {
			return nil, err
		}
		return json.Marshal(result)
	case CMD_EXISTS:
		logger.Infof("checking if %s exists", message[1:])
		result, err := s.DispatchExists(string(message[1:]))
		if err != nil {
			return nil, err
		}
		return json.Marshal(result)
	case CMD_COMPLETIONS:
		logger.Infof("completing %s", message[1:])
		result, err := s.DispatchCompetions(string(message[1:]))
		if err != nil {
			return nil, err
		}
		return json.Marshal(result)
	case CMD_KEYS:
		logger.Infof("listing keys")
		result := s.DispatchKeys()
		return json.Marshal(result)
	}

	return []byte{}, errors.New("invalid command")
}

/*
Insert a key to the Trie.
*/
func (s *ThreadSafeDispatcher) DispatchInsert(key string) (bool, error) {
	s.dispatcherMutex.Lock()
	defer s.dispatcherMutex.Unlock()

	return s.trie.Add(key)
}

/*
Delete a key from the Trie.
*/
func (s *ThreadSafeDispatcher) DispatchDelete(key string) (bool, error) {
	s.dispatcherMutex.Lock()
	defer s.dispatcherMutex.Unlock()

	return s.trie.Remove(key)
}

/*
Check if a key exists in the Trie.
*/
func (s *ThreadSafeDispatcher) DispatchExists(key string) (bool, error) {
	s.dispatcherMutex.Lock()
	defer s.dispatcherMutex.Unlock()

	return s.trie.Has(key)
}

/*
Get the keys starting with a prefix in the Trie.
*/
func (s *ThreadSafeDispatcher) DispatchCompetions(prefix string) ([]string, error) {
	s.dispatcherMutex.Lock()
	defer s.dispatcherMutex.Unlock()

	return s.trie.Completions(prefix)
}

/*
List all keys, in depth-first order.
*/
func (s *ThreadSafeDispatcher) DispatchKeys() []string {
	s.dispatcherMutex.Lock()
	defer s.dispatcherMutex.Unlock()

	return s.trie.Keys()
}
