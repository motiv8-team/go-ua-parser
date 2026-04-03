package uax

import "testing"

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
