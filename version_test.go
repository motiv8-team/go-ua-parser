package uax

import "testing"

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input                      string
		major, minor, patch, full  string
	}{
		{"123.0.6312.86", "123", "0", "6312", "123.0.6312.86"},
		{"17.4", "17", "4", "", "17.4"},
		{"10", "10", "", "", "10"},
		{"", "", "", "", ""},
		{"15.0.0", "15", "0", "0", "15.0.0"},
	}
	for _, tt := range tests {
		maj, min, pat := splitVersion(tt.input)
		if maj != tt.major {
			t.Errorf("splitVersion(%q) major = %q, want %q", tt.input, maj, tt.major)
		}
		if min != tt.minor {
			t.Errorf("splitVersion(%q) minor = %q, want %q", tt.input, min, tt.minor)
		}
		if pat != tt.patch {
			t.Errorf("splitVersion(%q) patch = %q, want %q", tt.input, pat, tt.patch)
		}
	}
}

func TestMajorVersion(t *testing.T) {
	if got := majorVersion("123.0.6312.86"); got != "123" {
		t.Errorf("got %q, want 123", got)
	}
	if got := majorVersion(""); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}
