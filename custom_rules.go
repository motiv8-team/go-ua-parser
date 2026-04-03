package uax

import "fmt"

// BrowserRule is a user-facing browser detection rule.
type BrowserRule struct {
	Token  string
	Name   string
	Family string
	Engine string
	Match  string // "exact", "contains", "prefix"
}

// BotRule is a user-facing bot detection rule.
type BotRule struct {
	Token    string
	Name     string
	Class    BotClass
	Vendor   string
	Verified bool
	Match    string
}

// DeviceRule is a user-facing device detection rule.
type DeviceRule struct {
	Token  string
	Type   string
	Vendor string
	Model  string
	Match  string
}

func validateMatchString(s string) error {
	switch s {
	case "exact", "contains", "prefix", "":
		return nil
	default:
		return fmt.Errorf("invalid match type %q: must be exact, contains, or prefix", s)
	}
}

func parseMatchType(s string) matchType {
	switch s {
	case "contains":
		return matchContains
	case "prefix":
		return matchPrefix
	default:
		return matchExact
	}
}

func browserRuleToInternal(br BrowserRule) rule {
	return rule{
		pattern:       br.Token,
		matchType:     parseMatchType(br.Match),
		browserName:   br.Name,
		browserFamily: br.Family,
		engineName:    br.Engine,
	}
}

func botRuleToInternal(br BotRule) rule {
	return rule{
		pattern:       br.Token,
		matchType:     parseMatchType(br.Match),
		botName:       br.Name,
		botClass:      br.Class,
		botVendor:     br.Vendor,
		botIsVerified: br.Verified,
	}
}

func deviceRuleToInternal(dr DeviceRule) rule {
	return rule{
		pattern:      dr.Token,
		matchType:    parseMatchType(dr.Match),
		deviceType:   dr.Type,
		deviceVendor: dr.Vendor,
		deviceModel:  dr.Model,
	}
}

// WithCustomBrowserRules adds custom browser rules checked before builtins.
func WithCustomBrowserRules(rules []BrowserRule) Option {
	return func(c *parserConfig) {
		c.customBrowserRules = rules
	}
}

// WithCustomBotRules adds custom bot rules checked before builtins.
func WithCustomBotRules(rules []BotRule) Option {
	return func(c *parserConfig) {
		c.customBotRules = rules
	}
}

// WithCustomDeviceRules adds custom device rules checked before builtins.
func WithCustomDeviceRules(rules []DeviceRule) Option {
	return func(c *parserConfig) {
		c.customDeviceRules = rules
	}
}
