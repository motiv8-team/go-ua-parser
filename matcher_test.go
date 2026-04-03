package uax

import (
	"regexp"
	"testing"
)

func TestRuleMatchExact(t *testing.T) {
	r := rule{
		pattern:   "Googlebot",
		matchType: matchExact,
	}
	if !r.matches("Googlebot") {
		t.Error("should match exact")
	}
	if r.matches("googlebot") {
		t.Error("should not match case-insensitive by default")
	}
}

func TestRuleMatchRegex(t *testing.T) {
	r := rule{
		pattern:   `(?i)^Googlebot`,
		matchType: matchRegex,
		re:        regexp.MustCompile(`(?i)^Googlebot`),
	}
	if !r.matches("Googlebot/2.1") {
		t.Error("should match regex")
	}
	if r.matches("NotGooglebot") {
		t.Error("should not match")
	}
}

func TestRuleMatchContains(t *testing.T) {
	r := rule{
		pattern:   "bot",
		matchType: matchContains,
	}
	if !r.matches("Googlebot/2.1") {
		t.Error("should match contains")
	}
	if r.matches("Chrome/123") {
		t.Error("should not match")
	}
}

func TestRuleTableLookup(t *testing.T) {
	rt := ruleTable{
		rules: []rule{
			{pattern: "Chrome", matchType: matchExact, browserName: "Chrome", engineName: "Blink"},
			{pattern: "Firefox", matchType: matchExact, browserName: "Firefox", engineName: "Gecko"},
		},
	}
	rt.buildTrie()

	r, ok := rt.lookup("Chrome")
	if !ok {
		t.Fatal("Chrome should match")
	}
	if r.browserName != "Chrome" {
		t.Errorf("browserName = %q, want Chrome", r.browserName)
	}
}
