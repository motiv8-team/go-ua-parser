package uax

import "testing"

func TestWindows11FromClientHints(t *testing.T) {
	p, _ := NewParser()
	r := p.Parse(Input{
		UAString: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/123.0 Safari/537.36",
		ClientHints: ClientHints{
			Platform:        "Windows",
			PlatformVersion: "15.0.0",
		},
	})
	if r.OS.Name != "Windows 11" {
		t.Errorf("os = %q, want 'Windows 11'", r.OS.Name)
	}
}

func TestWindows10FromClientHints(t *testing.T) {
	p, _ := NewParser()
	r := p.Parse(Input{
		UAString: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/123.0 Safari/537.36",
		ClientHints: ClientHints{
			Platform:        "Windows",
			PlatformVersion: "10.0.0",
		},
	})
	if r.OS.Name != "Windows" {
		t.Errorf("os = %q, want 'Windows'", r.OS.Name)
	}
}

func TestWindows10WithoutClientHints(t *testing.T) {
	p, _ := NewParser()
	r := p.ParseString("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/123.0 Safari/537.36")
	if r.OS.Name != "Windows" {
		t.Errorf("os = %q, want 'Windows' (can't detect 11 without CH)", r.OS.Name)
	}
}

func TestWindows11Boundary(t *testing.T) {
	p, _ := NewParser()
	// Version 13.0.0 is the boundary for Windows 11
	r := p.Parse(Input{
		UAString: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/123.0",
		ClientHints: ClientHints{
			Platform:        "Windows",
			PlatformVersion: "13.0.0",
		},
	})
	if r.OS.Name != "Windows 11" {
		t.Errorf("os = %q, want 'Windows 11' (13.0.0 is Windows 11)", r.OS.Name)
	}
}
