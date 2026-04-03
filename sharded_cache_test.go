package uax

import (
	"fmt"
	"sync"
	"testing"
)

func TestShardedCacheBasic(t *testing.T) {
	p, _ := NewParser()
	sc := NewShardedCache(p, 16, 100)

	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/123.0 Safari/537.36"
	r1 := sc.ParseString(ua)
	r2 := sc.ParseString(ua)
	if r1.Browser.Name != r2.Browser.Name || r1.Browser.Name != "Chrome" {
		t.Errorf("got %q and %q, want Chrome", r1.Browser.Name, r2.Browser.Name)
	}
}

func TestShardedCacheStats(t *testing.T) {
	p, _ := NewParser()
	sc := NewShardedCache(p, 4, 100)

	sc.ParseString("UA-1")
	sc.ParseString("UA-1") // hit
	sc.ParseString("UA-2")

	stats := sc.Stats()
	if stats.Hits != 1 {
		t.Errorf("hits = %d, want 1", stats.Hits)
	}
	if stats.Misses != 2 {
		t.Errorf("misses = %d, want 2", stats.Misses)
	}
}

func TestShardedCacheConcurrent(t *testing.T) {
	p, _ := NewParser()
	sc := NewShardedCache(p, 16, 1000)

	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			ua := fmt.Sprintf("UA-%d", idx%50)
			_ = sc.ParseString(ua)
		}(i)
	}
	wg.Wait()

	stats := sc.Stats()
	if stats.Size == 0 {
		t.Error("cache should have entries")
	}
}

func BenchmarkShardedCacheHit(b *testing.B) {
	p, _ := NewParser()
	sc := NewShardedCache(p, 16, 10000)
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/123.0 Safari/537.36"
	sc.ParseString(ua) // warm

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sc.ParseString(ua)
	}
}

func BenchmarkShardedCacheContended(b *testing.B) {
	p, _ := NewParser()
	sc := NewShardedCache(p, 16, 10000)
	uas := []string{
		"Mozilla/5.0 Chrome/123.0",
		"Mozilla/5.0 Firefox/124.0",
		"Googlebot/2.1",
		"curl/8.4.0",
	}
	for _, ua := range uas {
		sc.ParseString(ua) // warm
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			sc.ParseString(uas[i%len(uas)])
			i++
		}
	})
}
