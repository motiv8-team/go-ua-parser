package uax

import "sync"

// CacheStats reports the current state of a CachedParser's LRU cache.
type CacheStats struct {
	Size   int
	Hits   int64
	Misses int64
}

// lruList is an intrusive doubly-linked list used by both CachedParser and cacheShard.
type lruList struct {
	head *lruEntry
	tail *lruEntry
}

func (l *lruList) promote(e *lruEntry) {
	if l.head == e {
		return
	}
	l.remove(e)
	l.pushFront(e)
}

func (l *lruList) pushFront(e *lruEntry) {
	e.prev = nil
	e.next = l.head
	if l.head != nil {
		l.head.prev = e
	}
	l.head = e
	if l.tail == nil {
		l.tail = e
	}
}

func (l *lruList) remove(e *lruEntry) {
	if e.prev != nil {
		e.prev.next = e.next
	} else {
		l.head = e.next
	}
	if e.next != nil {
		e.next.prev = e.prev
	} else {
		l.tail = e.prev
	}
	e.prev = nil
	e.next = nil
}

func (l *lruList) evictTail() *lruEntry {
	if l.tail == nil {
		return nil
	}
	e := l.tail
	l.remove(e)
	return e
}

// CachedParser wraps a Parser with an LRU cache to avoid re-parsing identical UA strings.
type CachedParser struct {
	parser  *Parser
	mu      sync.RWMutex
	entries map[string]*lruEntry
	lru     lruList
	maxSize int
	hits    int64
	misses  int64
}

type lruEntry struct {
	key    string
	result Result
	prev   *lruEntry
	next   *lruEntry
}

// NewCachedParser creates a CachedParser backed by the given Parser with an LRU cache of maxSize entries.
// If maxSize is <= 0, it defaults to 1024.
func NewCachedParser(p *Parser, maxSize int) *CachedParser {
	if maxSize <= 0 {
		maxSize = 1024
	}
	return &CachedParser{
		parser:  p,
		entries: make(map[string]*lruEntry, maxSize),
		maxSize: maxSize,
	}
}

// ParseString parses a raw UA string, returning a cached Result when available.
func (c *CachedParser) ParseString(ua string) Result {
	c.mu.RLock()
	if e, ok := c.entries[ua]; ok {
		c.mu.RUnlock()
		c.mu.Lock()
		// Re-validate: entry may have been evicted between RUnlock and Lock
		if _, still := c.entries[ua]; still {
			c.lru.promote(e)
			c.hits++
		}
		c.mu.Unlock()
		return e.result
	}
	c.mu.RUnlock()

	result := c.parser.ParseString(ua)

	c.mu.Lock()
	if e, ok := c.entries[ua]; ok {
		c.lru.promote(e)
		c.hits++
		c.mu.Unlock()
		return e.result
	}

	c.misses++
	e := &lruEntry{key: ua, result: result}
	c.entries[ua] = e
	c.lru.pushFront(e)

	for len(c.entries) > c.maxSize {
		if tail := c.lru.evictTail(); tail != nil {
			delete(c.entries, tail.key)
		}
	}
	c.mu.Unlock()

	return result
}

// Parse parses an Input, using the cache for UA-only lookups and bypassing it when Client Hints are present.
func (c *CachedParser) Parse(input Input) Result {
	if input.HasClientHints() {
		return c.parser.Parse(input)
	}
	return c.ParseString(input.UAString)
}

// Stats returns a snapshot of the cache hit/miss counters and current entry count.
func (c *CachedParser) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return CacheStats{
		Size:   len(c.entries),
		Hits:   c.hits,
		Misses: c.misses,
	}
}
