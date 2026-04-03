// Package uax is a high-performance, zero-allocation HTTP User-Agent parser
// for Go. It extracts browser, engine, OS, CPU, device, in-app browser, and
// bot information from User-Agent strings and Client Hints headers.
//
//go:generate sh -c "cd cmd/uagen && go run ."
package uax

import "sync"

var (
	defaultParser     *Parser
	defaultParserOnce sync.Once
)

func getDefaultParser() *Parser {
	defaultParserOnce.Do(func() {
		var err error
		defaultParser, err = NewParser()
		if err != nil {
			panic("uax: failed to initialize default parser: " + err.Error())
		}
	})
	return defaultParser
}

// Parse parses a raw User-Agent string using the default global parser.
// Safe for concurrent use.
func Parse(ua string) Result {
	return getDefaultParser().ParseString(ua)
}

// ShortBrowser returns a compact "Name Major" string, e.g. "Chrome 123".
func (r Result) ShortBrowser() string {
	if r.Browser.Name == "" {
		return ""
	}
	if r.Browser.Major != "" {
		return r.Browser.Name + " " + r.Browser.Major
	}
	return r.Browser.Name
}

// ShortOS returns a compact "Name Major" string, e.g. "Windows 10".
func (r Result) ShortOS() string {
	if r.OS.Name == "" {
		return ""
	}
	if r.OS.Major != "" {
		return r.OS.Name + " " + r.OS.Major
	}
	return r.OS.Name
}

// DeviceClass returns the device type or "unknown" if not detected.
func (r Result) DeviceClass() string {
	if r.Device.Type == "" {
		return "unknown"
	}
	return r.Device.Type
}
