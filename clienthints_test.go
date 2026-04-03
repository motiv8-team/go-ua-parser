package uax

import "testing"

func TestClientHintsFromMap(t *testing.T) {
	m := map[string]string{
		"Sec-CH-UA":                  `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`,
		"Sec-CH-UA-Platform":         `"Windows"`,
		"Sec-CH-UA-Platform-Version": `"15.0.0"`,
		"Sec-CH-UA-Arch":             `"x86"`,
		"Sec-CH-UA-Model":            `""`,
		"Sec-CH-UA-Full-Version":     `"124.0.6367.91"`,
	}
	ch := ClientHintsFromMap(m)
	if ch.Platform != "Windows" {
		t.Errorf("platform = %q, want Windows", ch.Platform)
	}
	if ch.PlatformVersion != "15.0.0" {
		t.Errorf("platformVersion = %q, want 15.0.0", ch.PlatformVersion)
	}
	if ch.Arch != "x86" {
		t.Errorf("arch = %q, want x86", ch.Arch)
	}
	if ch.FullVersion != "124.0.6367.91" {
		t.Errorf("fullVersion = %q, want 124.0.6367.91", ch.FullVersion)
	}
}

func TestClientHintsBrandParsing(t *testing.T) {
	ch := ClientHints{
		UA: `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`,
	}
	brands := ch.ParseBrands()
	if len(brands) != 3 {
		t.Fatalf("got %d brands, want 3", len(brands))
	}
	found := false
	for _, b := range brands {
		if b.Brand == "Google Chrome" && b.Version == "124" {
			found = true
		}
	}
	if !found {
		t.Error("expected Google Chrome brand with version 124")
	}
}

func TestInputHasClientHints(t *testing.T) {
	i := Input{UAString: "test"}
	if i.HasClientHints() {
		t.Error("empty hints should return false")
	}
	i.ClientHints.Platform = "Windows"
	if !i.HasClientHints() {
		t.Error("non-empty hints should return true")
	}
}
