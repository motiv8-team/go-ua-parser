package uax

import "testing"

func TestCachedParserBasic(t *testing.T) {
	p, _ := NewParser()
	cp := NewCachedParser(p, 100)

	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/123.0 Safari/537.36"
	r1 := cp.ParseString(ua)
	r2 := cp.ParseString(ua)

	if r1.Browser.Name != r2.Browser.Name {
		t.Error("cached result should be identical")
	}
	if r1.Browser.Name != "Chrome" {
		t.Errorf("browser = %q, want Chrome", r1.Browser.Name)
	}
}

func TestCachedParserEviction(t *testing.T) {
	p, _ := NewParser()
	cp := NewCachedParser(p, 2)

	cp.ParseString("UA-1")
	cp.ParseString("UA-2")
	cp.ParseString("UA-3")

	stats := cp.Stats()
	if stats.Size > 2 {
		t.Errorf("cache size = %d, want <= 2", stats.Size)
	}
}

func TestCachedParserConcurrent(t *testing.T) {
	p, _ := NewParser()
	cp := NewCachedParser(p, 1000)

	done := make(chan struct{})
	for i := 0; i < 100; i++ {
		go func() {
			_ = cp.ParseString("Mozilla/5.0 Chrome/123.0")
			done <- struct{}{}
		}()
	}
	for i := 0; i < 100; i++ {
		<-done
	}
}

func BenchmarkCachedParserHit(b *testing.B) {
	p, _ := NewParser()
	cp := NewCachedParser(p, 10000)
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/123.0.6312.86 Safari/537.36"
	cp.ParseString(ua) // warm

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cp.ParseString(ua)
	}
}
