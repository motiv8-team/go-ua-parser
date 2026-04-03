package uax

import (
	"hash/maphash"
	"sync"
)

// ShardedCache distributes cached entries across multiple independent LRU shards
// to reduce lock contention under high concurrency.
type ShardedCache struct {
	parser     *Parser
	shards     []cacheShard
	shardCount uint64
	seed       maphash.Seed
}

type cacheShard struct {
	mu      sync.RWMutex
	entries map[string]*lruEntry
	lru     lruList
	maxSize int
	hits    int64
	misses  int64
}

// NewShardedCache creates a sharded cache with the given number of shards
// and per-shard capacity. Total capacity = shards * perShardSize.
func NewShardedCache(p *Parser, shards, perShardSize int) *ShardedCache {
	if shards <= 0 {
		shards = 16
	}
	if perShardSize <= 0 {
		perShardSize = 256
	}
	sc := &ShardedCache{
		parser:     p,
		shards:     make([]cacheShard, shards),
		shardCount: uint64(shards),
		seed:       maphash.MakeSeed(),
	}
	for i := range sc.shards {
		sc.shards[i].entries = make(map[string]*lruEntry, perShardSize)
		sc.shards[i].maxSize = perShardSize
	}
	return sc
}

// ParseString parses a UA string with sharded caching.
func (sc *ShardedCache) ParseString(ua string) Result {
	shard := &sc.shards[sc.shardIndex(ua)]
	return shard.getOrParse(ua, sc.parser)
}

// Parse parses an Input. Only caches if no Client Hints (CH makes each request unique).
func (sc *ShardedCache) Parse(input Input) Result {
	if input.HasClientHints() {
		return sc.parser.Parse(input)
	}
	return sc.ParseString(input.UAString)
}

// Stats returns aggregated cache statistics across all shards.
func (sc *ShardedCache) Stats() CacheStats {
	var s CacheStats
	for i := range sc.shards {
		sc.shards[i].mu.RLock()
		s.Size += len(sc.shards[i].entries)
		s.Hits += sc.shards[i].hits
		s.Misses += sc.shards[i].misses
		sc.shards[i].mu.RUnlock()
	}
	return s
}

func (sc *ShardedCache) shardIndex(key string) uint64 {
	var h maphash.Hash
	h.SetSeed(sc.seed)
	h.WriteString(key)
	return h.Sum64() % sc.shardCount
}

func (s *cacheShard) getOrParse(ua string, p *Parser) Result {
	// Fast path: read lock
	s.mu.RLock()
	if e, ok := s.entries[ua]; ok {
		s.mu.RUnlock()
		s.mu.Lock()
		// Re-validate: entry may have been evicted between RUnlock and Lock
		if _, still := s.entries[ua]; still {
			s.lru.promote(e)
			s.hits++
		}
		s.mu.Unlock()
		return e.result
	}
	s.mu.RUnlock()

	// Miss: parse
	result := p.ParseString(ua)

	s.mu.Lock()
	// Double-check
	if e, ok := s.entries[ua]; ok {
		s.lru.promote(e)
		s.hits++
		s.mu.Unlock()
		return e.result
	}

	s.misses++
	e := &lruEntry{key: ua, result: result}
	s.entries[ua] = e
	s.lru.pushFront(e)

	for len(s.entries) > s.maxSize {
		if tail := s.lru.evictTail(); tail != nil {
			delete(s.entries, tail.key)
		}
	}
	s.mu.Unlock()
	return result
}
