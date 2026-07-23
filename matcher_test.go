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

func TestBuildTrieCompilesRegex(t *testing.T) {
	// A matchRegex rule declared with only a pattern (no pre-compiled re) —
	// as produced by the code generator from `match: regex` YAML — must be
	// compiled by buildTrie and become functional.
	rt := ruleTable{
		rules: []rule{
			{pattern: `(?i)PetalBot`, matchType: matchRegex, botName: "PetalBot"},
		},
	}
	rt.buildTrie()

	if rt.rules[0].re == nil {
		t.Fatal("buildTrie should compile the regex pattern into re")
	}
	r, ok := rt.lookup("petalbot")
	if !ok {
		t.Fatal("regex rule should match via lookup")
	}
	if r.botName != "PetalBot" {
		t.Errorf("botName = %q, want PetalBot", r.botName)
	}
	if _, ok := rt.lookup("Chrome"); ok {
		t.Error("regex rule should not match unrelated candidate")
	}
}

func TestBuildTrieSkipsInvalidRegex(t *testing.T) {
	// An un-compilable pattern (e.g. a backreference, unsupported by RE2)
	// must be disabled rather than panic, so a bad rule never matches.
	rt := ruleTable{
		rules: []rule{
			{pattern: `(a)\1`, matchType: matchRegex, botName: "Bad"},
		},
	}
	rt.buildTrie() // must not panic
	if rt.rules[0].re != nil {
		t.Error("invalid regex should leave re nil (disabled)")
	}
	if _, ok := rt.lookup("aa"); ok {
		t.Error("disabled regex rule must never match")
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
