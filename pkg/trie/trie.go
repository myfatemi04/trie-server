package trie

import "sync"

type Trie struct {
	Leaves      map[byte]*Trie
	IsLeaf      bool
	accessMutex sync.Mutex
}

func CreateTrie() Trie {
	return Trie{
		Leaves: make(map[byte]*Trie),
		IsLeaf: false,
	}
}

// Add adds a word to the trie, returning whether there was a change
func (t *Trie) Add(key string) bool {
	t.accessMutex.Lock()
	defer t.accessMutex.Unlock()

	if len(key) == 0 {
		previousIsLeaf := t.IsLeaf
		t.IsLeaf = true
		return t.IsLeaf != previousIsLeaf
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
	changed := t.Leaves[first].Add(key[1:])
	return changed
}

// Remove removes a word from the trie, returning whether there was a change
func (t *Trie) Remove(key string) bool {
	t.accessMutex.Lock()
	defer t.accessMutex.Unlock()

	if len(key) == 0 {
		previousIsLeaf := t.IsLeaf
		t.IsLeaf = false
		return t.IsLeaf != previousIsLeaf
	}

	first := key[0]

	if _, ok := t.Leaves[first]; !ok {
		// The key doesn't exist
		return false
	}

	// Remove the key from the leaf
	changed := t.Leaves[first].Remove(key[1:])

	if changed {
		leaf := t.Leaves[first]
		if leaf.IsEmpty() {
			// Remove the leaf if it's empty
			delete(t.Leaves, first)
		}
	}

	return changed
}

// Has returns whether the trie contains a key
func (t *Trie) Has(key string) bool {
	t.accessMutex.Lock()
	defer t.accessMutex.Unlock()

	if len(key) == 0 {
		return t.IsLeaf
	}

	first := key[0]

	if _, ok := t.Leaves[first]; !ok {
		return false
	}

	return t.Leaves[first].Has(key[1:])
}

// IsEmpty returns whether the trie is empty
func (t *Trie) IsEmpty() bool {
	t.accessMutex.Lock()
	defer t.accessMutex.Unlock()

	return len(t.Leaves) == 0 && !t.IsLeaf
}

// Return all keys that begin with `prefix`
func (t *Trie) Completions(prefix string) []string {
	t.accessMutex.Lock()
	defer t.accessMutex.Unlock()

	if len(prefix) == 0 {
		return t.Keys()
	}

	first := prefix[0]

	if _, ok := t.Leaves[first]; !ok {
		return []string{}
	}

	return t.Leaves[first].Completions(prefix[1:])
}

// Size returns the number of keys in the trie
func (t *Trie) Size() int {
	t.accessMutex.Lock()
	defer t.accessMutex.Unlock()

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
	t.accessMutex.Lock()
	defer t.accessMutex.Unlock()

	keys := make([]string, 0, t.Size())

	for first, leaf := range t.Leaves {
		for _, leafKey := range leaf.Keys() {
			keys = append(keys, string(first)+leafKey)
		}
	}

	return keys
}
