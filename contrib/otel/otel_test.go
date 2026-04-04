package uaxotel

import (
	"testing"

	uax "github.com/motiv8-team/go-ua-parser"
)

func TestAttributes(t *testing.T) {
	r := uax.Result{
		Browser: uax.Browser{Name: "Chrome", Version: "123.0", Major: "123", Family: "Chromium"},
		Engine:  uax.Engine{Name: "Blink", Version: "123.0"},
		OS:      uax.OS{Name: "Windows", Version: "10.0", Major: "10"},
		Device:  uax.Device{Type: "desktop", Vendor: "Dell", Model: "XPS"},
		CPU:     uax.CPU{Architecture: "amd64"},
	}

	attrs := Attributes(r)
	if len(attrs) == 0 {
		t.Fatal("expected non-empty attributes")
	}

	got := make(map[string]string)
	for _, a := range attrs {
		got[string(a.Key)] = a.Value.Emit()
	}

	checks := map[string]string{
		"http.user_agent.browser.name":    "Chrome",
		"http.user_agent.browser.version": "123.0",
		"http.user_agent.browser.family":  "Chromium",
		"http.user_agent.engine.name":     "Blink",
		"http.user_agent.os.name":         "Windows",
		"http.user_agent.device.type":     "desktop",
		"http.user_agent.cpu.architecture": "amd64",
	}
	for k, want := range checks {
		if got[k] != want {
			t.Errorf("attr %s = %q, want %q", k, got[k], want)
		}
	}
}

func TestAttributesBot(t *testing.T) {
	r := uax.Result{
		IsBot: true,
		Bot:   uax.Bot{Name: "Googlebot", Class: uax.BotSearch},
	}

	attrs := Attributes(r)
	got := make(map[string]string)
	for _, a := range attrs {
		got[string(a.Key)] = a.Value.Emit()
	}

	if got["http.user_agent.is_bot"] != "true" {
		t.Error("expected is_bot=true")
	}
	if got["http.user_agent.bot.name"] != "Googlebot" {
		t.Errorf("bot.name = %q, want Googlebot", got["http.user_agent.bot.name"])
	}
	if got["http.user_agent.bot.class"] != "search" {
		t.Errorf("bot.class = %q, want search", got["http.user_agent.bot.class"])
	}
}

func TestAttributesEmpty(t *testing.T) {
	attrs := Attributes(uax.Result{})
	if len(attrs) != 0 {
		t.Errorf("empty result should produce 0 attrs, got %d", len(attrs))
	}
}
