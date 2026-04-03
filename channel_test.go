package uax

import "testing"

func TestBrowserChannel(t *testing.T) {
	p, _ := NewParser()
	tests := []struct {
		ua      string
		channel string
	}{
		{"Mozilla/5.0 (X11; Linux x86_64; rv:126.0a1) Gecko/20100101 Firefox/126.0a1", "nightly"},
		{"Mozilla/5.0 (X11; Linux x86_64; rv:125.0b9) Gecko/20100101 Firefox/125.0b9", "beta"},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.6312.86 Safari/537.36", ""},
	}
	for _, tt := range tests {
		r := p.ParseString(tt.ua)
		if r.Browser.Channel != tt.channel {
			t.Errorf("channel for %q = %q, want %q", tt.ua[:40], r.Browser.Channel, tt.channel)
		}
	}
}
