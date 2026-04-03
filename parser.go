package uax

import (
	"fmt"
	"time"
)

// ---------------------------------------------------------------------------
// Options
// ---------------------------------------------------------------------------

// Option configures a Parser.
type Option func(*parserConfig)

type parserConfig struct {
	enableBotDetection bool
	customBrowserRules []BrowserRule
	customBotRules     []BotRule
	customDeviceRules  []DeviceRule
	preParseHook       PreParseHookFunc
	postParseHook      PostParseHookFunc
}

// WithBotDetection enables or disables bot/crawler detection.
func WithBotDetection(enabled bool) Option {
	return func(c *parserConfig) { c.enableBotDetection = enabled }
}

// ---------------------------------------------------------------------------
// Parser
// ---------------------------------------------------------------------------

// Parser is a reusable, concurrency-safe user-agent parser.
type Parser struct {
	cfg      parserConfig
	browsers ruleTable
	engines  ruleTable
	oses     ruleTable
	devices  ruleTable
	cpus     ruleTable
	bots     ruleTable
	apps     ruleTable
}

// NewParser creates a Parser pre-loaded with builtin rules.
func NewParser(opts ...Option) (*Parser, error) {
	p := &Parser{
		cfg: parserConfig{
			enableBotDetection: true,
		},
	}
	for _, o := range opts {
		o(&p.cfg)
	}
	// Validate custom rules
	for _, r := range p.cfg.customBrowserRules {
		if err := validateMatchString(r.Match); err != nil {
			return nil, fmt.Errorf("browser rule %q: %w", r.Token, err)
		}
	}
	for _, r := range p.cfg.customBotRules {
		if err := validateMatchString(r.Match); err != nil {
			return nil, fmt.Errorf("bot rule %q: %w", r.Token, err)
		}
	}
	for _, r := range p.cfg.customDeviceRules {
		if err := validateMatchString(r.Match); err != nil {
			return nil, fmt.Errorf("device rule %q: %w", r.Token, err)
		}
	}
	loadBuiltinRules(p)
	return p, nil
}

// ParseString is a convenience wrapper that parses a raw UA string.
func (p *Parser) ParseString(ua string) Result {
	return p.Parse(Input{UAString: ua})
}

// Parse parses an Input and returns a Result.
func (p *Parser) Parse(input Input) Result {
	var out Result
	p.ParseInto(input, &out)
	return out
}

// ParseInto parses into a caller-provided Result, avoiding allocation.
func (p *Parser) ParseInto(input Input, out *Result) {
	*out = Result{}
	if p.cfg.preParseHook != nil {
		p.cfg.preParseHook(input)
	}

	var start time.Time
	if p.cfg.postParseHook != nil {
		start = time.Now()
	}

	out.UAString = input.UAString

	var tz tokenizer
	tz.reset(input.UAString)
	tokens := tz.tokenize()

	p.matchBot(tokens, input.UAString, out)
	if !out.IsBot {
		p.matchApp(tokens, out)
		p.matchBrowser(tokens, out)
		p.matchEngine(tokens, out)
	}
	p.matchOS(tokens, out)
	p.matchDevice(tokens, out)
	p.matchCPU(tokens, out)

	if input.HasClientHints() {
		mergeClientHints(&input.ClientHints, out)
	}

	computeDerived(out)

	if p.cfg.postParseHook != nil {
		p.cfg.postParseHook(input, *out, time.Since(start))
	}
}

// DetectBot is a quick bot-only detection path. Returns an empty Bot if
// bot detection is disabled via WithBotDetection(false).
func (p *Parser) DetectBot(ua string) Bot {
	var tz tokenizer
	tz.reset(ua)
	tokens := tz.tokenize()
	var out Result
	p.matchBot(tokens, ua, &out)
	return out.Bot
}

// ---------------------------------------------------------------------------
// Builtin rules (Task 11)
// ---------------------------------------------------------------------------

func loadBuiltinRules(p *Parser) {
	// Browsers: custom first, then builtins
	var browserRules []rule
	for _, cr := range p.cfg.customBrowserRules {
		browserRules = append(browserRules, browserRuleToInternal(cr))
	}
	browserRules = append(browserRules, builtinBrowserRules...)
	p.browsers.rules = browserRules
	p.browsers.buildTrie()

	p.engines.rules = builtinEngineRules
	p.engines.buildTrie()

	p.oses.rules = nil

	// Devices: custom first, then builtins
	var deviceRules []rule
	for _, cr := range p.cfg.customDeviceRules {
		deviceRules = append(deviceRules, deviceRuleToInternal(cr))
	}
	deviceRules = append(deviceRules, builtinDeviceRules...)
	p.devices.rules = deviceRules
	p.devices.buildTrie()

	p.cpus.rules = nil

	// Bots: custom first, then builtins
	var botRules []rule
	for _, cr := range p.cfg.customBotRules {
		botRules = append(botRules, botRuleToInternal(cr))
	}
	botRules = append(botRules, builtinBotRules...)
	p.bots.rules = botRules
	p.bots.buildTrie()

	p.apps.rules = builtinAppRules
	p.apps.buildTrie()
}
