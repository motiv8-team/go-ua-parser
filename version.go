package uax

import "strings"

// splitVersion splits "major.minor.patch..." into its first three components.
func splitVersion(v string) (major, minor, patch string) {
	if v == "" {
		return "", "", ""
	}
	parts := strings.SplitN(v, ".", 4)
	if len(parts) >= 1 {
		major = parts[0]
	}
	if len(parts) >= 2 {
		minor = parts[1]
	}
	if len(parts) >= 3 {
		patch = parts[2]
	}
	return
}

// majorVersion returns the major component of a version string.
func majorVersion(v string) string {
	if i := strings.IndexByte(v, '.'); i >= 0 {
		return v[:i]
	}
	return v
}
