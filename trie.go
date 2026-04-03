package uax

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
		c := key[i]
		next := node.children[c]
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

