package uax

import (
	"fmt"
	"regexp"
)

// BrowserRule is a user-facing browser detection rule.
type BrowserRule struct {
	Token  string
	Name   string
	Family string
	Engine string
	Match  string // "exact", "contains", "prefix", or "regex" (RE2 syntax)
}

// BotRule is a user-facing bot detection rule.
type BotRule struct {
	Token    string
	Name     string
	Class    BotClass
	Vendor   string
	Verified bool
	Match    string // "exact", "contains", "prefix", or "regex" (RE2 syntax)
}

// DeviceRule is a user-facing device detection rule.
type DeviceRule struct {
	Token  string
	Type   string
	Vendor string
	Model  string
	Match  string // "exact", "contains", "prefix", or "regex" (RE2 syntax)
}

func validateMatchString(s string) error {
	switch s {
	case "exact", "contains", "prefix", "regex", "":
		return nil
	default:
		return fmt.Errorf("invalid match type %q: must be exact, contains, prefix, or regex", s)
	}
}

func parseMatchType(s string) matchType {
	switch s {
	case "contains":
		return matchContains
	case "prefix":
		return matchPrefix
	case "regex":
		return matchRegex
	default:
		return matchExact
	}
}

// compileIfRegex compiles r.pattern into r.re when r.matchType is
// matchRegex, returning an error if the pattern is not valid RE2 syntax.
// Compilation happens once here, during custom-rule conversion (which runs
// as part of NewParser's validation), so the rule that ends up in the rule
// table always carries a working compiled regexp — never a nil one that
// would silently never match.
func compileIfRegex(r *rule, pattern string) error {
	if r.matchType != matchRegex {
		return nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid regex: %w", err)
	}
	r.re = re
	return nil
}

func browserRuleToInternal(br BrowserRule) (rule, error) {
	r := rule{
		pattern:       br.Token,
		matchType:     parseMatchType(br.Match),
		browserName:   br.Name,
		browserFamily: br.Family,
		engineName:    br.Engine,
	}
	if err := compileIfRegex(&r, br.Token); err != nil {
		return rule{}, err
	}
	return r, nil
}

func botRuleToInternal(br BotRule) (rule, error) {
	r := rule{
		pattern:       br.Token,
		matchType:     parseMatchType(br.Match),
		botName:       br.Name,
		botClass:      br.Class,
		botVendor:     br.Vendor,
		botIsVerified: br.Verified,
	}
	if err := compileIfRegex(&r, br.Token); err != nil {
		return rule{}, err
	}
	return r, nil
}

func deviceRuleToInternal(dr DeviceRule) (rule, error) {
	r := rule{
		pattern:      dr.Token,
		matchType:    parseMatchType(dr.Match),
		deviceType:   dr.Type,
		deviceVendor: dr.Vendor,
		deviceModel:  dr.Model,
	}
	if err := compileIfRegex(&r, dr.Token); err != nil {
		return rule{}, err
	}
	return r, nil
}

// WithCustomBrowserRules adds custom browser rules. Precedence is not simply
// "custom before builtin": exact-match rules are looked up in a single trie
// where the first-inserted rule for a given token wins, so a custom exact
// rule overrides a builtin exact rule for the same token. Non-exact rules
// (contains/prefix/regex) are only consulted when no exact rule matched the
// token, and are tried in insertion order (custom rules first, then
// builtins). Rules match against individual UA product tokens — not the
// whole UA string.
func WithCustomBrowserRules(rules []BrowserRule) Option {
	return func(c *parserConfig) {
		c.customBrowserRules = rules
	}
}

// WithCustomBotRules adds custom bot rules. Precedence is not simply "custom
// before builtin": exact-match rules are looked up in a single trie where
// the first-inserted rule for a given token wins, so a custom exact rule
// overrides a builtin exact rule for the same token. Non-exact rules
// (contains/prefix/regex) are only consulted when no exact rule matched the
// token, and are tried in insertion order (custom rules first, then
// builtins). Rules match against individual UA product tokens and, for bot
// detection, also against semicolon-separated parts of parenthesized UA
// comments (e.g. "compatible; Googlebot/2.1; ...") — not the whole UA
// string.
func WithCustomBotRules(rules []BotRule) Option {
	return func(c *parserConfig) {
		c.customBotRules = rules
	}
}

// WithCustomDeviceRules adds custom device rules. Precedence is not simply
// "custom before builtin": exact-match rules are looked up in a single trie
// where the first-inserted rule for a given token wins, so a custom exact
// rule overrides a builtin exact rule for the same token. Non-exact rules
// (contains/prefix/regex) are only consulted when no exact rule matched the
// token, and are tried in insertion order (custom rules first, then
// builtins). Exact rules match against individual UA product tokens;
// non-exact rules are also checked against the full text of parenthesized UA
// comments — in neither case against the whole UA string.
func WithCustomDeviceRules(rules []DeviceRule) Option {
	return func(c *parserConfig) {
		c.customDeviceRules = rules
	}
}
