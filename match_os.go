package uax

import "strings"

func (p *Parser) matchOS(tokens []token, out *Result) {
	// Try ruleTable first
	for i := range tokens {
		tk := &tokens[i]
		if tk.name != "" {
			if r, ok := p.oses.lookup(tk.name); ok {
				out.OS.Name = r.osName
				out.OS.Version = tk.version
				m, mi, pa := splitVersion(tk.version)
				out.OS.Major = m
				out.OS.Minor = mi
				out.OS.Patch = pa
				return
			}
		}
	}
	// Fallback: parse from comment tokens
	for i := range tokens {
		if tokens[i].comment != "" {
			osInfo := parseOSFromComment(tokens[i].comment)
			if osInfo.Name != "" {
				out.OS = osInfo
				return
			}
		}
	}
}

func parseOSFromComment(comment string) OS {
	var os OS

	switch {
	case strings.Contains(comment, "Android"):
		os.Name = "Android"
		// Extract version: "Android 14" or "Android 14;"
		if idx := strings.Index(comment, "Android "); idx >= 0 {
			rest := comment[idx+8:]
			end := strings.IndexAny(rest, ";) ")
			if end < 0 {
				end = len(rest)
			}
			os.Version = rest[:end]
			os.Major, os.Minor, os.Patch = splitVersion(os.Version)
		}

	case strings.Contains(comment, "iPhone OS") || (strings.Contains(comment, "iPad") && strings.Contains(comment, "OS")):
		if strings.Contains(comment, "iPad") {
			os.Name = "iPadOS"
		} else {
			os.Name = "iOS"
		}
		// "CPU iPhone OS 17_4 like Mac OS X" or "CPU OS 17_4 like Mac OS X"
		ver := extractOSVersion(comment, "OS ")
		if ver != "" {
			ver = strings.ReplaceAll(ver, "_", ".")
			os.Version = ver
			os.Major, os.Minor, os.Patch = splitVersion(ver)
		}

	case strings.Contains(comment, "Windows NT"):
		os.Name = "Windows"
		if idx := strings.Index(comment, "Windows NT "); idx >= 0 {
			rest := comment[idx+11:]
			end := strings.IndexAny(rest, ";) ")
			if end < 0 {
				end = len(rest)
			}
			ntVer := rest[:end]
			os.Version = mapWindowsNT(ntVer)
			os.Major, os.Minor, os.Patch = splitVersion(os.Version)
		}

	case strings.Contains(comment, "Mac OS X"):
		os.Name = "macOS"
		if idx := strings.Index(comment, "Mac OS X "); idx >= 0 {
			rest := comment[idx+9:]
			end := strings.IndexAny(rest, ";) ")
			if end < 0 {
				end = len(rest)
			}
			ver := strings.ReplaceAll(rest[:end], "_", ".")
			os.Version = ver
			os.Major, os.Minor, os.Patch = splitVersion(ver)
		}

	case strings.Contains(comment, "CrOS"):
		os.Name = "ChromeOS"

	case strings.Contains(comment, "HarmonyOS") || strings.Contains(comment, "OpenHarmony"):
		os.Name = "HarmonyOS"

	case strings.Contains(comment, "Tizen"):
		os.Name = "Tizen"

	case strings.Contains(comment, "KaiOS"):
		os.Name = "KaiOS"

	case strings.Contains(comment, "Sailfish"):
		os.Name = "Sailfish"

	case strings.Contains(comment, "Firefox OS"):
		os.Name = "Firefox OS"

	case strings.Contains(comment, "Fuchsia"):
		os.Name = "Fuchsia"

	case strings.Contains(comment, "FreeBSD"):
		os.Name = "FreeBSD"

	case strings.Contains(comment, "OpenBSD"):
		os.Name = "OpenBSD"

	case strings.Contains(comment, "NetBSD"):
		os.Name = "NetBSD"

	case strings.Contains(comment, "Ubuntu"):
		os.Name = "Ubuntu"

	case strings.Contains(comment, "Fedora"):
		os.Name = "Fedora"

	case strings.Contains(comment, "Debian"):
		os.Name = "Debian"

	case strings.Contains(comment, "SUSE") || strings.Contains(comment, "openSUSE"):
		os.Name = "SUSE"

	case strings.Contains(comment, "Mint"):
		os.Name = "Linux Mint"

	case strings.Contains(comment, "Arch"):
		os.Name = "Arch Linux"

	case strings.Contains(comment, "Gentoo"):
		os.Name = "Gentoo"

	case strings.Contains(comment, "Manjaro"):
		os.Name = "Manjaro"

	case strings.Contains(comment, "BlackBerry") || strings.Contains(comment, "BB10"):
		os.Name = "BlackBerry"

	case strings.Contains(comment, "Symbian"):
		os.Name = "Symbian"

	case strings.Contains(comment, "Windows Phone"):
		os.Name = "Windows Phone"

	case strings.Contains(comment, "Windows CE"):
		os.Name = "Windows CE"

	case strings.Contains(comment, "Windows Mobile"):
		os.Name = "Windows Mobile"

	case strings.Contains(comment, "webOS") || strings.Contains(comment, "hpwOS"):
		os.Name = "webOS"

	case strings.Contains(comment, "Bada"):
		os.Name = "Bada"

	case strings.Contains(comment, "MeeGo"):
		os.Name = "MeeGo"

	case strings.Contains(comment, "Maemo"):
		os.Name = "Maemo"

	case strings.Contains(comment, "watchOS"):
		os.Name = "watchOS"

	case strings.Contains(comment, "tvOS"):
		os.Name = "tvOS"

	case strings.Contains(comment, "Linux"):
		os.Name = "Linux"
	}

	return os
}

func extractOSVersion(comment, marker string) string {
	idx := strings.Index(comment, marker)
	if idx < 0 {
		return ""
	}
	rest := comment[idx+len(marker):]
	end := strings.IndexAny(rest, " ;)")
	if end < 0 {
		end = len(rest)
	}
	return rest[:end]
}

func mapWindowsNT(ntVer string) string {
	switch ntVer {
	case "10.0":
		return "10"
	case "6.3":
		return "8.1"
	case "6.2":
		return "8"
	case "6.1":
		return "7"
	case "6.0":
		return "Vista"
	case "5.1", "5.2":
		return "XP"
	}
	return ntVer
}
