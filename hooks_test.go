package uax

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestPreParseHook(t *testing.T) {
	var called atomic.Int32
	p, _ := NewParser(WithPreParseHook(func(input Input) {
		called.Add(1)
	}))
	p.ParseString("Mozilla/5.0 Chrome/123.0")
	p.ParseString("Googlebot/2.1")
	if called.Load() != 2 {
		t.Errorf("pre-parse hook called %d times, want 2", called.Load())
	}
}

func TestPostParseHook(t *testing.T) {
	var lastDuration time.Duration
	var lastBrowser string
	p, _ := NewParser(WithPostParseHook(func(input Input, result Result, d time.Duration) {
		lastDuration = d
		lastBrowser = result.Browser.Name
	}))
	p.ParseString("Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/123.0 Safari/537.36")
	if lastBrowser != "Chrome" {
		t.Errorf("post-parse hook browser = %q, want Chrome", lastBrowser)
	}
	if lastDuration <= 0 {
		t.Error("post-parse hook duration should be > 0")
	}
}

func TestNoHooksNoCost(t *testing.T) {
	p, _ := NewParser()
	r := p.ParseString("Chrome/123.0")
	if r.Browser.Name != "Chrome" {
		t.Errorf("browser = %q, want Chrome", r.Browser.Name)
	}
}
