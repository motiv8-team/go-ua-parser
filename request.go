package uax

import "net/http"

// clientHintsHeaders lists the standard Client Hints header names.
var clientHintsHeaders = []string{
	"Sec-CH-UA",
	"Sec-CH-UA-Platform",
	"Sec-CH-UA-Platform-Version",
	"Sec-CH-UA-Arch",
	"Sec-CH-UA-Model",
	"Sec-CH-UA-Full-Version",
	"Sec-CH-UA-Full-Version-List",
}

// ParseRequest parses a User-Agent and Client Hints from an HTTP request.
// Extracts the User-Agent header and all Sec-CH-UA-* headers automatically.
// Returns an empty Result if r is nil.
func (p *Parser) ParseRequest(r *http.Request) Result {
	if r == nil {
		return Result{}
	}

	input := Input{
		UAString: r.Header.Get("User-Agent"),
	}

	// Extract Client Hints headers if present
	for _, h := range clientHintsHeaders {
		if v := r.Header.Get(h); v != "" {
			if setter, ok := headerMap[h]; ok {
				setter(&input.ClientHints, v)
			}
		}
	}

	return p.Parse(input)
}
