package uax

import "strings"

func (p *Parser) matchBrowser(tokens []token, out *Result) {
	// Iterate tokens in reverse so that later, more-specific tokens
	// (e.g. "Edg" after "Chrome") win.
	for i := len(tokens) - 1; i >= 0; i-- {
		tk := &tokens[i]
		if tk.name == "" || tk.name == "Mozilla" || tk.name == "Safari" || tk.name == "Mobile" {
			continue
		}
		if r, ok := p.browsers.lookup(tk.name); ok {
			out.Browser.Name = r.browserName
			out.Browser.Family = r.browserFamily
			out.Browser.Version = tk.version
			out.Browser.Major = majorVersion(tk.version)
			if r.engineName != "" {
				out.Engine.Name = r.engineName
			}
			detectChannel(out)
			return
		}
	}
	// Fallback: if Safari token present and "Version" token exists
	for i := range tokens {
		if tokens[i].name == "Safari" {
			out.Browser.Name = "Safari"
			out.Browser.Family = "WebKit"
			out.Browser.Version = tokens[i].version
			out.Browser.Major = majorVersion(tokens[i].version)
			out.Engine.Name = "WebKit"
			// Look for Version/ token for the real version
			for j := range tokens {
				if tokens[j].name == "Version" {
					out.Browser.Version = tokens[j].version
					out.Browser.Major = majorVersion(tokens[j].version)
					break
				}
			}
			detectChannel(out)
			return
		}
	}
}

func detectChannel(out *Result) {
	ver := out.Browser.Version
	if ver == "" {
		return
	}
	switch {
	case strings.HasSuffix(ver, "a1") || strings.HasSuffix(ver, "a2"):
		out.Browser.Channel = "nightly"
	case containsBetaMarker(ver, out.Browser.Family):
		out.Browser.Channel = "beta"
	}
}

func containsBetaMarker(ver, family string) bool {
	if family != "Firefox" {
		return false
	}
	for i := 0; i < len(ver); i++ {
		if ver[i] == 'b' && i > 0 && i < len(ver)-1 && ver[i-1] >= '0' && ver[i-1] <= '9' && ver[i+1] >= '0' && ver[i+1] <= '9' {
			return true
		}
	}
	return false
}
