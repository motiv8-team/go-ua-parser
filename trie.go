package uax

// trieNode uses a [256]*trieNode array for O(1) byte-indexed child lookup.
// Each node is 2KB on 64-bit systems. With ~3000 nodes across all rule tables,
// total trie memory is ~6MB per Parser instance. This is the right tradeoff
// for a performance library — lookup is 4ns vs 11ns (compact slice) or 100ns (map).
// If memory is a concern, reuse a single Parser instance (it's concurrency-safe).
type trieNode struct {
	children [256]*trieNode
	id       int
	terminal bool
}

type trie struct {
	root trieNode
}

func newTrie() *trie {
	return &trie{}
}

func (t *trie) insert(key string, id int) {
	node := &t.root
	for i := 0; i < len(key); i++ {
		c := key[i]
		if node.children[c] == nil {
			node.children[c] = &trieNode{}
		}
		node = node.children[c]
	}
	if !node.terminal {
		node.id = id
		node.terminal = true
	}
}

func (t *trie) match(key string) (int, bool) {
	node := &t.root
	for i := 0; i < len(key); i++ {
		next := node.children[key[i]]
		if next == nil {
			return 0, false
		}
		node = next
	}
	if node.terminal {
		return node.id, true
	}
	return 0, false
}
