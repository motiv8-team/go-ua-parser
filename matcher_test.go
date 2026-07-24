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

func TestBuildTrieIndexesPrecompiledRegex(t *testing.T) {
	// buildTrie no longer compiles regex patterns itself — that now happens
	// at codegen time for builtins (cmd/uagen) and during custom-rule
	// conversion for user-supplied rules (custom_rules.go). buildTrie must
	// still read-only index a matchRegex rule that already carries a
	// compiled re, and lookup must use it.
	rt := ruleTable{
		rules: []rule{
			{pattern: `(?i)PetalBot`, matchType: matchRegex, re: regexp.MustCompile(`(?i)PetalBot`), botName: "PetalBot"},
		},
	}
	rt.buildTrie()

	if rt.rules[0].re == nil {
		t.Fatal("buildTrie must not clear a pre-compiled re")
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

func TestBuildTrieNilRegexNeverMatches(t *testing.T) {
	// buildTrie must not attempt to compile regex patterns — compilation is
	// the responsibility of the code generator (builtins) and custom-rule
	// validation in NewParser (custom rules), both of which reject invalid
	// patterns before a rule ever reaches a ruleTable. If a matchRegex rule
	// somehow arrives with a nil re, buildTrie must index it without
	// panicking, and it must never match.
	rt := ruleTable{
		rules: []rule{
			{pattern: `(a)\1`, matchType: matchRegex, botName: "Bad"},
		},
	}
	rt.buildTrie() // must not panic
	if rt.rules[0].re != nil {
		t.Error("buildTrie must not compile regex patterns; re should remain nil")
	}
	if _, ok := rt.lookup("aa"); ok {
		t.Error("a matchRegex rule with nil re must never match")
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
