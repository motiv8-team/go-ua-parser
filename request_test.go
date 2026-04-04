package uax

import (
	"net/http"
	"testing"
)

func TestParseRequest(t *testing.T) {
	p, _ := NewParser()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.6312.86 Safari/537.36")

	r := p.ParseRequest(req)
	if r.Browser.Name != "Chrome" {
		t.Errorf("browser = %q, want Chrome", r.Browser.Name)
	}
	if r.OS.Name != "Windows" {
		t.Errorf("os = %q, want Windows", r.OS.Name)
	}
}

func TestParseRequestWithClientHints(t *testing.T) {
	p, _ := NewParser()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/123.0 Safari/537.36")
	req.Header.Set("Sec-CH-UA", `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`)
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-CH-UA-Platform-Version", `"14.4.1"`)
	req.Header.Set("Sec-CH-UA-Arch", `"arm64"`)
	req.Header.Set("Sec-CH-UA-Full-Version", `"124.0.6367.91"`)

	r := p.ParseRequest(req)
	// Client Hints should override UA-derived values
	if r.OS.Name != "macOS" {
		t.Errorf("os = %q, want macOS (from CH)", r.OS.Name)
	}
	if r.CPU.Architecture != "arm64" {
		t.Errorf("arch = %q, want arm64 (from CH)", r.CPU.Architecture)
	}
	if r.Browser.Version != "124.0.6367.91" {
		t.Errorf("version = %q, want 124.0.6367.91 (from CH)", r.Browser.Version)
	}
}

func TestParseRequestNilRequest(t *testing.T) {
	p, _ := NewParser()
	r := p.ParseRequest(nil)
	if r.Browser.Name != "" {
		t.Error("nil request should return empty result")
	}
}

func TestParseRequestNoUA(t *testing.T) {
	p, _ := NewParser()
	req, _ := http.NewRequest("GET", "/", nil)
	r := p.ParseRequest(req)
	if r.Browser.Name != "" {
		t.Error("no UA header should return empty browser")
	}
}

func TestParseRequestInto(t *testing.T) {
	p, _ := NewParser()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.6312.86 Safari/537.36")

	var r Result
	p.ParseRequestInto(req, &r)
	if r.Browser.Name != "Chrome" {
		t.Errorf("browser = %q, want Chrome", r.Browser.Name)
	}
	if r.OS.Name != "Windows" {
		t.Errorf("os = %q, want Windows", r.OS.Name)
	}
}

func TestParseRequestIntoWithClientHints(t *testing.T) {
	p, _ := NewParser()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/123.0 Safari/537.36")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-CH-UA-Arch", `"arm64"`)

	var r Result
	p.ParseRequestInto(req, &r)
	if r.OS.Name != "macOS" {
		t.Errorf("os = %q, want macOS (from CH)", r.OS.Name)
	}
	if r.CPU.Architecture != "arm64" {
		t.Errorf("arch = %q, want arm64 (from CH)", r.CPU.Architecture)
	}
}

func TestParseRequestIntoNil(t *testing.T) {
	p, _ := NewParser()
	var r Result
	r.Browser.Name = "stale"
	p.ParseRequestInto(nil, &r)
	if r.Browser.Name != "" {
		t.Error("nil request should zero result")
	}
}
