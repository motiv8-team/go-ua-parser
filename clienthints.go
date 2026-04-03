package uax

import "strings"

// ClientHints holds parsed Sec-CH-UA-* header values.
type ClientHints struct {
	UA              string `json:"ua,omitempty"`
	Platform        string `json:"platform,omitempty"`
	PlatformVersion string `json:"platformVersion,omitempty"`
	Arch            string `json:"arch,omitempty"`
	Model           string `json:"model,omitempty"`
	FullVersion     string `json:"fullVersion,omitempty"`
	FullVersionList string `json:"fullVersionList,omitempty"`
}

// Brand is a single entry from the Sec-CH-UA or Sec-CH-UA-Full-Version-List header.
type Brand struct {
	Brand   string
	Version string
}

// headerMap maps standard HTTP header names to ClientHints fields.
var headerMap = map[string]func(*ClientHints, string){
	"Sec-CH-UA":                   func(ch *ClientHints, v string) { ch.UA = v },
	"Sec-CH-UA-Platform":          func(ch *ClientHints, v string) { ch.Platform = unquote(v) },
	"Sec-CH-UA-Platform-Version":  func(ch *ClientHints, v string) { ch.PlatformVersion = unquote(v) },
	"Sec-CH-UA-Arch":              func(ch *ClientHints, v string) { ch.Arch = unquote(v) },
	"Sec-CH-UA-Model":             func(ch *ClientHints, v string) { ch.Model = unquote(v) },
	"Sec-CH-UA-Full-Version":      func(ch *ClientHints, v string) { ch.FullVersion = unquote(v) },
	"Sec-CH-UA-Full-Version-List": func(ch *ClientHints, v string) { ch.FullVersionList = v },
}

// ClientHintsFromMap creates ClientHints from a map of HTTP header name→value.
func ClientHintsFromMap(m map[string]string) ClientHints {
	var ch ClientHints
	for k, v := range m {
		if setter, ok := headerMap[k]; ok {
			setter(&ch, v)
		}
	}
	return ch
}

// ParseBrands parses the Sec-CH-UA style header into brand entries.
// Format: "BrandA";v="1", "BrandB";v="2"
func (ch *ClientHints) ParseBrands() []Brand {
	raw := ch.UA
	if ch.FullVersionList != "" {
		raw = ch.FullVersionList
	}
	if raw == "" {
		return nil
	}
	return parseBrandList(raw)
}

func parseBrandList(raw string) []Brand {
	var brands []Brand
	for _, entry := range strings.Split(raw, ",") {
		entry = strings.TrimSpace(entry)
		parts := strings.SplitN(entry, ";", 2)
		if len(parts) < 2 {
			continue
		}
		brand := unquote(strings.TrimSpace(parts[0]))
		verPart := strings.TrimSpace(parts[1])
		if idx := strings.Index(verPart, "="); idx >= 0 {
			ver := unquote(strings.TrimSpace(verPart[idx+1:]))
			brands = append(brands, Brand{Brand: brand, Version: ver})
		}
	}
	return brands
}

// unquote removes surrounding double quotes from a string.
func unquote(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}
