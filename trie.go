package main

import (
	"errors"
)

type Trie struct {
	Leaves map[byte]*Trie
	IsLeaf bool
}

func NewTrie() *Trie {
	return &Trie{
		Leaves: make(map[byte]*Trie),
		IsLeaf: false,
	}
}

const MAX_KEY_LENGTH = 256

// Add adds a word to the trie, returning whether there was a change
func (t *Trie) Add(key string) (bool, error) {
	if len(key) >= MAX_KEY_LENGTH {
		return false, errors.New("key is too long")
	}

	if len(key) == 0 {
		previousIsLeaf := t.IsLeaf
		t.IsLeaf = true
		return t.IsLeaf != previousIsLeaf, nil
	}

	first := key[0]

	if _, ok := t.Leaves[first]; !ok {
		// Add the leaf if it doesn't exist
		t.Leaves[first] = &Trie{
			Leaves: make(map[uint8]*Trie),
			IsLeaf: false,
		}
	}

	// Add the key to the leaf
	return t.Leaves[first].Add(key[1:])
}

// Remove removes a word from the trie, returning whether there was a change
func (t *Trie) Remove(key string) (bool, error) {
	if len(key) >= MAX_KEY_LENGTH {
		return false, errors.New("key is too long")
	}

	if len(key) == 0 {
		previousIsLeaf := t.IsLeaf
		t.IsLeaf = false
		return t.IsLeaf != previousIsLeaf, nil
	}

	first := key[0]

	if _, ok := t.Leaves[first]; !ok {
		// The key doesn't exist
		return false, nil
	}

	// Remove the key from the leaf
	changed, err := t.Leaves[first].Remove(key[1:])
	if err != nil {
		return changed, err
	}

	if changed {
		leaf := t.Leaves[first]
		if leaf.IsEmpty() {
			// Remove the leaf if it's empty
			delete(t.Leaves, first)
		}
	}

	return changed, err
}

// Has returns whether the trie contains a key
func (t *Trie) Has(key string) (bool, error) {
	if len(key) >= MAX_KEY_LENGTH {
		return false, errors.New("key is too long")
	}

	if len(key) == 0 {
		return t.IsLeaf, nil
	}

	first := key[0]

	if _, ok := t.Leaves[first]; !ok {
		return false, nil
	}

	return t.Leaves[first].Has(key[1:])
}

// IsEmpty returns whether the trie is empty
func (t *Trie) IsEmpty() bool {
	return len(t.Leaves) == 0 && !t.IsLeaf
}

// Return all keys that begin with `prefix`
func (t *Trie) Completions(prefix string) ([]string, error) {
	if len(prefix) == 0 {
		return t.Keys(), nil
	}

	first := prefix[0]

	// This prefix doesn't exist in the tree
	if _, ok := t.Leaves[first]; !ok {
		return []string{}, nil
	}

	leafCompletions, err := t.Leaves[first].Completions(prefix[1:])
	if err != nil {
		return nil, err
	}

	prefixedCompletions := make([]string, 0, len(leafCompletions))
	for _, leafCompletion := range leafCompletions {
		prefixedCompletions = append(prefixedCompletions, string(first)+leafCompletion)
	}
	return prefixedCompletions, nil
}

// Size returns the number of keys in the trie
func (t *Trie) Size() int {
	size := 0

	for _, leaf := range t.Leaves {
		size += leaf.Size()
	}

	if t.IsLeaf {
		size++
	}

	return size
}

// Keys returns all keys in the trie
func (t *Trie) Keys() []string {
	keys := make([]string, 0, t.Size())
	if t.IsLeaf {
		keys = append(keys, "")
	}

	for first, leaf := range t.Leaves {
		for _, leafKey := range leaf.Keys() {
			keys = append(keys, string(first)+leafKey)
		}
	}

	return keys
}
