package main

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

// enum for various command codes
const (
	CMD_INSERT = iota
	CMD_DELETE
	CMD_EXISTS
	CMD_COMPLETIONS
	CMD_KEYS
)

const ASCII_0 = 48

/*
DispatchRaw takes a command encoded as [COMMAND CODE][KEY] and dispatches it to the appropriate method.
It returns two values: the result of the command, and an error.
The result of the command is a JSON string.
If there is an error, the error message is returned as the second value.

Command codes:
	0: Insert a key
	1: Delete a key
	2: Check if a key exists
	3: Generate completions for a prefix
	4: List all keys in the trie

Because there is only one possible argument, this protocol is fairly straightforward.
If there is an argument, it is simply the remaining string after the first byte, which
is the command code.

Example Commands
	0foo: Insert "foo"
	1foo: Delete "foo"
	2foo: Check if "foo" exists
	3foo: Generate completions for "foo"
	4: List all keys in the trie

*/
func (s *ThreadSafeDispatcher) DispatchRaw(message []byte) ([]byte, error) {
	if len(message) == 0 {
		return nil, errors.New("empty message")
	}

	// First byte is the command code
	command := message[0] - ASCII_0

	switch command {
	case CMD_INSERT:
		logger.Infof("inserting %s", message[1:])

		// result: true if the key was inserted, false if it was already present
		result, err := s.DispatchInsert(string(message[1:]))
		if err != nil {
			return nil, err
		}

		return json.Marshal(result)

	case CMD_DELETE:
		logger.Infof("deleting %s", message[1:])

		// result: true if the key was deleted, false if it was not present
		result, err := s.DispatchDelete(string(message[1:]))
		if err != nil {
			return nil, err
		}

		return json.Marshal(result)

	case CMD_EXISTS:
		logger.Infof("checking if %s exists", message[1:])

		// result: true if the key exists, false if it does not
		result, err := s.DispatchExists(string(message[1:]))
		if err != nil {
			return nil, err
		}

		return json.Marshal(result)

	case CMD_COMPLETIONS:
		logger.Infof("completing %s", message[1:])

		// result: list of completions
		result, err := s.DispatchCompetions(string(message[1:]))
		if err != nil {
			return nil, err
		}

		return json.Marshal(result)

	case CMD_KEYS:
		logger.Infof("listing keys")

		if len(message) > 1 {
			return nil, errors.New("key listing command takes no arguments")
		}

		// result: list of keys
		result := s.DispatchKeys()
		return json.Marshal(result)

	}

	return []byte{}, errors.New("invalid command")
}

/*
Insert a key to the Trie, returning whether there was a change (thread-safe).
*/
func (s *ThreadSafeDispatcher) DispatchInsert(key string) (bool, error) {
	s.dispatcherMutex.Lock()
	defer s.dispatcherMutex.Unlock()

	return s.trie.Add(key)
}

/*
Delete a key from the Trie, returning whether there was a change (thread-safe).
*/
func (s *ThreadSafeDispatcher) DispatchDelete(key string) (bool, error) {
	s.dispatcherMutex.Lock()
	defer s.dispatcherMutex.Unlock()

	return s.trie.Remove(key)
}

/*
Check if a key exists in the Trie (thread-safe).
*/
func (s *ThreadSafeDispatcher) DispatchExists(key string) (bool, error) {
	s.dispatcherMutex.Lock()
	defer s.dispatcherMutex.Unlock()

	return s.trie.Has(key)
}

/*
Get the keys starting with a prefix in the Trie (thread-safe).
*/
func (s *ThreadSafeDispatcher) DispatchCompetions(prefix string) ([]string, error) {
	s.dispatcherMutex.Lock()
	defer s.dispatcherMutex.Unlock()

	return s.trie.Completions(prefix)
}

/*
List all keys, in depth-first order (thread-safe).
*/
func (s *ThreadSafeDispatcher) DispatchKeys() []string {
	s.dispatcherMutex.Lock()
	defer s.dispatcherMutex.Unlock()

	return s.trie.Keys()
}
