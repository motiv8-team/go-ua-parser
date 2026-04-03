package uax

import "strings"

func (p *Parser) matchDevice(tokens []token, out *Result) {
	// Try ruleTable (consoles, smarttv, etc.)
	// For contains-match rules, check against the full UA assembled from tokens.
	for i := range tokens {
		tk := &tokens[i]
		if tk.name != "" {
			if r, ok := p.devices.lookup(tk.name); ok {
				out.Device.Type = r.deviceType
				out.Device.Vendor = r.deviceVendor
				if r.deviceModel != "" {
					out.Device.Model = r.deviceModel
				}
				return
			}
		}
	}
	// Also check comments for contains-match device rules
	for i := range tokens {
		if tokens[i].comment != "" {
			for j := range p.devices.nonExactIdx {
				idx := p.devices.nonExactIdx[j]
				r := &p.devices.rules[idx]
				if r.matches(tokens[i].comment) {
					out.Device.Type = r.deviceType
					out.Device.Vendor = r.deviceVendor
					return
				}
			}
		}
	}
	// Fallback: infer from OS
	if out.Device.Type == "" {
		dev := inferDeviceFromOS(out.OS, tokens)
		if dev.Type != "" {
			out.Device = dev
		}
	}
}

func inferDeviceFromOS(os OS, tokens []token) Device {
	switch os.Name {
	case "iOS":
		// Check for iPad in tokens/comments
		return Device{Type: "mobile", Vendor: "Apple", Model: "iPhone"}
	case "iPadOS":
		return Device{Type: "tablet", Vendor: "Apple", Model: "iPad"}
	case "Android":
		// Check for tablet indicators
		hasMobile := false
		for i := range tokens {
			if tokens[i].name == "Mobile" {
				hasMobile = true
				break
			}
			if tokens[i].comment != "" && strings.Contains(tokens[i].comment, "Mobile") {
				hasMobile = true
				break
			}
		}
		if hasMobile {
			return Device{Type: "mobile"}
		}
		return Device{Type: "tablet"}
	case "Windows", "macOS", "Linux", "ChromeOS":
		return Device{Type: "desktop"}
	}
	return Device{}
}
