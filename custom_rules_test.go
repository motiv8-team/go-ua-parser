package uax

import (
	"strings"
	"sync"
	"testing"
)

func TestCustomBrowserRule(t *testing.T) {
	p, _ := NewParser(WithCustomBrowserRules([]BrowserRule{
		{Token: "MyBrowser", Name: "My Custom Browser", Family: "Chromium", Engine: "Blink", Match: "exact"},
	}))
	r := p.ParseString("Mozilla/5.0 MyBrowser/1.0")
	if r.Browser.Name != "My Custom Browser" {
		t.Errorf("browser = %q, want 'My Custom Browser'", r.Browser.Name)
	}
}

func TestCustomBotRule(t *testing.T) {
	p, _ := NewParser(WithCustomBotRules([]BotRule{
		{Token: "InternalBot", Name: "Our Bot", Class: BotMonitor, Vendor: "Us", Match: "exact"},
	}))
	r := p.ParseString("InternalBot/1.0")
	if !r.IsBot {
		t.Error("should detect custom bot")
	}
	if r.Bot.Name != "Our Bot" {
		t.Errorf("bot name = %q, want 'Our Bot'", r.Bot.Name)
	}
}

func TestCustomDeviceRule(t *testing.T) {
	p, _ := NewParser(WithCustomDeviceRules([]DeviceRule{
		{Token: "MyKiosk", Type: "embedded", Vendor: "Acme", Model: "Kiosk v2", Match: "exact"},
	}))
	r := p.ParseString("Mozilla/5.0 MyKiosk/3.0")
	if r.Device.Type != "embedded" {
		t.Errorf("device type = %q, want embedded", r.Device.Type)
	}
}

func TestCustomRuleOverridesBuiltin(t *testing.T) {
	p, _ := NewParser(WithCustomBrowserRules([]BrowserRule{
		{Token: "Chrome", Name: "CustomChrome", Family: "Chromium", Engine: "Blink", Match: "exact"},
	}))
	r := p.ParseString("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/123.0 Safari/537.36")
	if r.Browser.Name != "CustomChrome" {
		t.Errorf("browser = %q, want CustomChrome", r.Browser.Name)
	}
}

func TestCustomBotRuleInvalidRegexReturnsError(t *testing.T) {
	// `(a)\1` is a backreference: valid in most regex engines but not
	// supported by Go's RE2-based regexp package. NewParser must reject it
	// with a descriptive, non-nil error rather than silently constructing a
	// parser with a dead rule (the pre-fix regression this guards against).
	_, err := NewParser(WithCustomBotRules([]BotRule{
		{Token: `(a)\1`, Name: "BadBot", Class: BotOther, Match: "regex"},
	}))
	if err == nil {
		t.Fatal("expected an error for an invalid regex custom bot rule, got nil")
	}
	if !strings.Contains(err.Error(), "BadBot") && !strings.Contains(err.Error(), `\1`) {
		t.Errorf("error %q should mention the offending rule", err.Error())
	}
}

func TestCustomBotRuleValidRegexMatches(t *testing.T) {
	p, err := NewParser(WithCustomBotRules([]BotRule{
		{Token: `(?i)^somebot$`, Name: "SomeBot", Class: BotOther, Match: "regex"},
	}))
	if err != nil {
		t.Fatalf("NewParser returned unexpected error: %v", err)
	}
	r := p.ParseString("somebot/1.0")
	if !r.IsBot {
		t.Fatal("should detect custom regex bot rule")
	}
	if r.Bot.Name != "SomeBot" {
		t.Errorf("bot name = %q, want SomeBot", r.Bot.Name)
	}
}

func TestNewParserCustomRegexConcurrent(t *testing.T) {
	// Constructing many parsers concurrently, each with its own custom regex
	// bot rule, must be race-free: buildTrie must never mutate shared rule
	// state (see matcher.go), and each parser's regex must be compiled
	// independently during NewParser/loadBuiltinRules. Run with -race.
	const n = 20
	var wg sync.WaitGroup
	errs := make([]error, n)
	bots := make([]bool, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			p, err := NewParser(WithCustomBotRules([]BotRule{
				{Token: `(?i)^concurrentbot$`, Name: "ConcurrentBot", Class: BotOther, Match: "regex"},
			}))
			if err != nil {
				errs[i] = err
				return
			}
			r := p.ParseString("concurrentbot/1.0")
			bots[i] = r.IsBot && r.Bot.Name == "ConcurrentBot"
		}(i)
	}
	wg.Wait()

	for i := 0; i < n; i++ {
		if errs[i] != nil {
			t.Fatalf("goroutine %d: NewParser returned unexpected error: %v", i, errs[i])
		}
		if !bots[i] {
			t.Fatalf("goroutine %d: expected ConcurrentBot to be detected", i)
		}
	}
}
