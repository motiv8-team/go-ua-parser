package uax

import "strings"

// ---------------------------------------------------------------------------
// Client Hints merge
// ---------------------------------------------------------------------------

func mergeClientHints(ch *ClientHints, out *Result) {
	if ch.Platform != "" {
		out.OS.Name = ch.Platform
	}
	if ch.PlatformVersion != "" {
		out.OS.Version = ch.PlatformVersion
		m, mi, pa := splitVersion(ch.PlatformVersion)
		out.OS.Major = m
		out.OS.Minor = mi
		out.OS.Patch = pa
	}
	if out.OS.Name == "Windows" && ch.PlatformVersion != "" {
		major, _, _ := splitVersion(ch.PlatformVersion)
		if parseIntSafe(major) >= 13 {
			out.OS.Name = "Windows 11"
		}
	}
	if ch.Arch != "" {
		out.CPU.Architecture = ch.Arch
	}
	if ch.Model != "" {
		out.Device.Model = ch.Model
	}
	if ch.FullVersion != "" {
		out.Browser.Version = ch.FullVersion
		out.Browser.Major = majorVersion(ch.FullVersion)
	}

	brands := ch.ParseBrands()
	for _, b := range brands {
		if isSkippableBrand(b.Brand) {
			continue
		}
		// Pick the first real (non-Chromium) brand
		if b.Brand != "Chromium" && b.Brand != "Google Chrome" {
			out.Browser.Name = b.Brand
			if ch.FullVersion == "" {
				out.Browser.Version = b.Version
				out.Browser.Major = majorVersion(b.Version)
			}
			lowerBrand := strings.ToLower(b.Brand)
			if strings.Contains(lowerBrand, "beta") {
				out.Browser.Channel = "beta"
			} else if strings.Contains(lowerBrand, "dev") {
				out.Browser.Channel = "dev"
			} else if strings.Contains(lowerBrand, "canary") {
				out.Browser.Channel = "canary"
			} else if strings.Contains(lowerBrand, "nightly") {
				out.Browser.Channel = "nightly"
			}
			return
		}
	}
	// If only Chromium / Google Chrome brands remain
	for _, b := range brands {
		if isSkippableBrand(b.Brand) {
			continue
		}
		if b.Brand == "Google Chrome" {
			out.Browser.Name = "Chrome"
		} else {
			out.Browser.Name = b.Brand
		}
		if ch.FullVersion == "" {
			out.Browser.Version = b.Version
			out.Browser.Major = majorVersion(b.Version)
		}
		lowerBrand := strings.ToLower(b.Brand)
		if strings.Contains(lowerBrand, "beta") {
			out.Browser.Channel = "beta"
		} else if strings.Contains(lowerBrand, "dev") {
			out.Browser.Channel = "dev"
		} else if strings.Contains(lowerBrand, "canary") {
			out.Browser.Channel = "canary"
		} else if strings.Contains(lowerBrand, "nightly") {
			out.Browser.Channel = "nightly"
		}
		return
	}
}

// isSkippableBrand returns true for GREASE brands (which contain intentionally
// odd characters like "Not;A=Brand") and for "Chromium", which is the base
// browser identity shared by many Chromium-based browsers. We skip Chromium
// so we can surface the more specific brand (e.g. "Microsoft Edge").
func isSkippableBrand(brand string) bool {
	return strings.Contains(brand, "Not") ||
		strings.Contains(brand, "Greasy") ||
		strings.Contains(brand, " Brand") ||
		brand == "" ||
		brand == "Chromium"
}

// ---------------------------------------------------------------------------
// Derived booleans
// ---------------------------------------------------------------------------

func computeDerived(out *Result) {
	out.IsBot = out.Bot.IsBot
	out.IsCrawler = out.Bot.IsBot && (out.Bot.Class == BotSearch || out.Bot.Class == BotSEO)

	dt := out.Device.Type
	switch dt {
	case "mobile":
		out.IsMobile = true
		out.Device.IsPhone = true
	case "tablet":
		out.IsTablet = true
		out.Device.IsTablet = true
	case "desktop":
		out.IsDesktop = true
		out.Device.IsDesktop = true
	case "smarttv":
		out.Device.IsTV = true
	case "console":
		// no top-level flag
	default:
		// Infer from OS if device type wasn't set
		if !out.IsBot && out.Device.Type == "" {
			if isDesktopOS(out.OS.Name) {
				out.IsDesktop = true
				out.Device.IsDesktop = true
				out.Device.Type = "desktop"
			}
		}
	}
}

func isDesktopOS(name string) bool {
	switch name {
	case "Windows", "Windows 11", "macOS", "Linux", "ChromeOS", "Ubuntu", "Fedora", "Debian":
		return true
	case "FreeBSD", "OpenBSD", "NetBSD", "SUSE", "Arch Linux", "Gentoo", "Manjaro", "Linux Mint":
		return true
	}
	return false
}

// parseIntSafe parses the leading decimal digits of s as an int.
// It returns 0 for empty or non-numeric strings and stops at the first
// non-digit character, so it never panics or errors.
func parseIntSafe(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			n = n*10 + int(s[i]-'0')
		} else {
			break
		}
	}
	return n
}
