package main

import (
	"errors"
)

/*
The Trie data structure stores a set of strings.
Each Trie object can act as a node in any part of the trie.
To make prefix lookups efficient, each node simply stores a map of (next character) --> (sub-trie).
For example, a Trie containing the words "foo", "bar", "baz", and "fo" would have the following structure:

HEAD
	- f
		- o
			- Issubtrie: true (end of "fo")
			- o
				- Issubtrie: true (end of "foo")
	- b
		- a
			- r
				- Issubtrie: true (end of "bar")
			- z
				- Issubtrie: true (end of "baz")

This way, to find all words with a given prefix, we just follow the subtries corresponding
to the prefix characters.

If we wanted to find all words beginning with "ba" in the tree, we start from HEAD,
and follow the subtries corresponding to the characters "b" and "a".

This yields a trie with the following structure:

HEAD
 - r: subtrie
 - z: subtrie

Now, we can just list all the keys for the subtries.
If the current trie were also a subtrie, then we would add the prefix to the list.

This method is implemented recursively.
The base case is that the prefix is empty, in which case we return the keys of the current trie.
Otherwise, we follow the subtrie corresponding to the first character of the prefix,
and recursively call the function on the subtrie with the remaining characters of the prefix.
(We must remember to prepend each of these results with the character of the subtrie that led to them)
Then, if the current trie is the end of a word, we add the current prefix to the results.
*/
type Trie struct {
	// The map of (next character) --> (sub-trie)
	Subtries map[byte]*Trie
	// If this is set to true, a word ends at this node
	IsEndOfWord bool
}

// Creates an empty trie
func NewTrie() *Trie {
	return &Trie{
		Subtries:    make(map[byte]*Trie),
		IsEndOfWord: false,
	}
}

// Keys can be a maximum length of 256 characters
const MAX_KEY_LENGTH = 256

// Add adds a word to the trie, returning whether there was a change
func (t *Trie) Add(key string) (bool, error) {
	// Verify that the key is not too long
	if len(key) >= MAX_KEY_LENGTH {
		return false, errors.New("key is too long")
	}

	// If the key is empty, then this node is the end of a word
	if len(key) == 0 {
		// Return whether there was a change
		previousIssubtrie := t.IsEndOfWord
		t.IsEndOfWord = true
		return t.IsEndOfWord != previousIssubtrie, nil
	}

	first := key[0]

	if _, ok := t.Subtries[first]; !ok {
		// Add a subtrie if a subtrie for the next character doesn't exist
		t.Subtries[first] = &Trie{
			Subtries:    make(map[uint8]*Trie),
			IsEndOfWord: false,
		}
	}

	// Add the key to the subtrie
	return t.Subtries[first].Add(key[1:])
}

// Remove removes a word from the trie, returning whether there was a change
func (t *Trie) Remove(key string) (bool, error) {
	// Verify that the key is not too long
	if len(key) >= MAX_KEY_LENGTH {
		return false, errors.New("key is too long")
	}

	if len(key) == 0 {
		// Return whether there was a change
		previousIssubtrie := t.IsEndOfWord
		t.IsEndOfWord = false
		return t.IsEndOfWord != previousIssubtrie, nil
	}

	// Recurse. Follow the path of subtries, and remove the node corresponding to the key.
	// If the node has no other children, then it can be removed.
	first := key[0]

	if _, ok := t.Subtries[first]; !ok {
		// The key doesn't exist
		return false, nil
	}

	// Remove the key from the subtrie
	changed, err := t.Subtries[first].Remove(key[1:])
	if err != nil {
		return changed, err
	}

	if changed {
		subtrie := t.Subtries[first]
		if subtrie.IsEmpty() {
			// Remove the subtrie if it's empty
			delete(t.Subtries, first)
		}
	}

	return changed, err
}

// Has returns whether the trie contains a key
func (t *Trie) Has(key string) (bool, error) {
	// Verify that the key is not too long
	if len(key) >= MAX_KEY_LENGTH {
		return false, errors.New("key is too long")
	}

	// Base case: if the key is blank, and this trie marks the end of a word,
	// then the key is in the trie
	if len(key) == 0 {
		return t.IsEndOfWord, nil
	}

	// Recurse. Follow the path of subtries, and check if the key is in the subtrie
	first := key[0]

	if _, ok := t.Subtries[first]; !ok {
		return false, nil
	}

	return t.Subtries[first].Has(key[1:])
}

// IsEmpty returns whether the trie is empty
func (t *Trie) IsEmpty() bool {
	return len(t.Subtries) == 0 && !t.IsEndOfWord
}

// Return all keys that begin with `prefix`
func (t *Trie) Completions(prefix string) ([]string, error) {
	// Verify that the prefix is not too long
	if len(prefix) >= MAX_KEY_LENGTH {
		return nil, errors.New("prefix is too long")
	}

	// If we reach the end of the prefix, simply return all characters in the trie
	if len(prefix) == 0 {
		return t.Keys(), nil
	}

	// Recurse. Follow the path of subtries, and check if the prefix is in the subtrie
	first := prefix[0]

	if _, ok := t.Subtries[first]; !ok {
		// This prefix doesn't exist in the tree
		return []string{}, nil
	}

	// Find completions from the corresponding subtrie
	subtrieCompletions, err := t.Subtries[first].Completions(prefix[1:])
	if err != nil {
		return nil, err
	}

	prefixedCompletions := make([]string, 0, len(subtrieCompletions))
	for _, subtrieCompletion := range subtrieCompletions {
		prefixedCompletions = append(prefixedCompletions, string(first)+subtrieCompletion)
	}
	return prefixedCompletions, nil
}

// Size returns the number of keys in the trie
func (t *Trie) Size() int {
	size := 0

	// Size = number of keys in subtries + whether or not this node is the end of a word
	for _, subtrie := range t.Subtries {
		size += subtrie.Size()
	}

	if t.IsEndOfWord {
		size++
	}

	return size
}

// Keys returns all keys in the trie
func (t *Trie) Keys() []string {
	keys := make([]string, 0, t.Size())
	// End of a word: add an empty string.
	if t.IsEndOfWord {
		keys = append(keys, "")
	}

	for characterThatLedToSubtrie, subtrie := range t.Subtries {
		for _, subtrieKey := range subtrie.Keys() {
			// Add keys from subtries, prefixed with character that led to them
			keys = append(keys, string(characterThatLedToSubtrie)+subtrieKey)
		}
	}

	return keys
}
