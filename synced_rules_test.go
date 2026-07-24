package uax

import "testing"

// This file pins behavior confirmed during a code review of the upstream
// rule sync (rules/*.yaml: bots 314->821, browsers 124->385, apps 59->368,
// engines 13->16). Each test documents a specific fixed point so that future
// rule edits cannot silently regress it.

// TestReclassificationGuards pins bot-class assignments for User-Agents that
// the review corrected: an AI-agent crawler that must not be classified as a
// generic bot, and a monitoring health-check probe that must not be
// classified as an AI agent.
func TestReclassificationGuards(t *testing.T) {
	p, _ := NewParser()
	tests := []struct {
		ua        string
		wantClass BotClass
	}{
		{"Mozilla/5.0 (compatible; GoogleAgent-Mariner/1.0; +http://www.google.com/bot.html)", BotAI},
		{"Cloudflare-Healthchecks/1.0", BotMonitor},
	}
	for _, tt := range tests {
		t.Run(tt.ua, func(t *testing.T) {
			r := p.ParseString(tt.ua)
			if !r.IsBot {
				t.Fatalf("IsBot = false, want true\n  UA: %s", tt.ua)
			}
			if r.Bot.Class != tt.wantClass {
				t.Errorf("Bot.Class = %q, want %q\n  UA: %s", r.Bot.Class, tt.wantClass, tt.ua)
			}
		})
	}
}

// TestHijackGuards pins User-Agents that superficially resemble a rule
// pattern (an HTTP library that contains "Boto", a HarmonyOS in-app webview
// whose comment block must not swallow the real Chrome/Blink token, a
// synthetic UA crafted to collide with an app-name substring rule, and a
// feature-phone UA whose Java/Obigo tokens must not be mistaken for a
// browser) but must NOT trigger detection in the wrong category.
func TestHijackGuards(t *testing.T) {
	p, _ := NewParser()

	t.Run("boto3 http library is not a browser or bot", func(t *testing.T) {
		ua := "Boto3/1.9.0 Python/3.6.5 Linux/4.14 Botocore/1.12.0"
		r := p.ParseString(ua)
		if r.Browser.Name != "" {
			t.Errorf("Browser.Name = %q, want empty\n  UA: %s", r.Browser.Name, ua)
		}
		if r.IsBot {
			t.Errorf("IsBot = true, want false\n  UA: %s", ua)
		}
	})

	t.Run("OpenHarmony ArkWeb webview still resolves Chrome", func(t *testing.T) {
		ua := "Mozilla/5.0 (Phone; OpenHarmony 5.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 ArkWeb/4.1.6.1 Mobile"
		r := p.ParseString(ua)
		if r.Browser.Name != "Chrome" {
			t.Errorf("Browser.Name = %q, want %q\n  UA: %s", r.Browser.Name, "Chrome", ua)
		}
	})

	t.Run("app-name substring collision does not produce a false app match", func(t *testing.T) {
		ua := "SomeApp AppName/HideX app_version/1.2"
		r := p.ParseString(ua)
		if r.App.Name != "" {
			t.Errorf("App.Name = %q, want empty\n  UA: %s", r.App.Name, ua)
		}
		if r.IsInApp {
			t.Errorf("IsInApp = true, want false\n  UA: %s", ua)
		}
	})

	t.Run("feature phone Obigo/Java tokens do not resolve a browser", func(t *testing.T) {
		ua := "LG-CU920 Obigo/Q05.1 MMS/LG-CU920 Java/ASVM/1.1"
		r := p.ParseString(ua)
		if r.Browser.Name != "" {
			t.Errorf("Browser.Name = %q, want empty\n  UA: %s", r.Browser.Name, ua)
		}
	})
}

// TestSyncedBotAndAppDetections pins a sample of new detections introduced by
// the upstream sync: an AI crawler, an in-app browser wrapper, matched by
// name, class and (where applicable) vendor.
func TestSyncedBotAndAppDetections(t *testing.T) {
	p, _ := NewParser()

	botTests := []struct {
		name       string
		ua         string
		wantName   string
		wantClass  BotClass
		wantVendor string
	}{
		{
			name:       "ClaudeBot",
			ua:         "Mozilla/5.0 (compatible; ClaudeBot/1.0; +claudebot@anthropic.com)",
			wantName:   "ClaudeBot",
			wantClass:  BotAI,
			wantVendor: "Anthropic",
		},
		{
			name:       "PerplexityBot",
			ua:         "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; PerplexityBot/1.0; +https://perplexity.ai/perplexitybot)",
			wantName:   "PerplexityBot",
			wantClass:  BotAI,
			wantVendor: "Perplexity",
		},
		{
			name:       "Bytespider",
			ua:         "Mozilla/5.0 (compatible; Bytespider; spider-feedback@bytedance.com)",
			wantName:   "Bytespider",
			wantClass:  BotAI,
			wantVendor: "ByteDance",
		},
	}
	for _, tt := range botTests {
		t.Run(tt.name, func(t *testing.T) {
			r := p.ParseString(tt.ua)
			if !r.IsBot {
				t.Fatalf("IsBot = false, want true\n  UA: %s", tt.ua)
			}
			if r.Bot.Name != tt.wantName {
				t.Errorf("Bot.Name = %q, want %q", r.Bot.Name, tt.wantName)
			}
			if r.Bot.Class != tt.wantClass {
				t.Errorf("Bot.Class = %q, want %q", r.Bot.Class, tt.wantClass)
			}
			if r.Bot.Vendor != tt.wantVendor {
				t.Errorf("Bot.Vendor = %q, want %q", r.Bot.Vendor, tt.wantVendor)
			}
		})
	}

	t.Run("ChatGPT in-app browser", func(t *testing.T) {
		ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) ChatGPT/1.2024.100"
		r := p.ParseString(ua)
		if r.App.Name != "ChatGPT" {
			t.Errorf("App.Name = %q, want %q\n  UA: %s", r.App.Name, "ChatGPT", ua)
		}
		if r.App.Version != "1.2024.100" {
			t.Errorf("App.Version = %q, want %q\n  UA: %s", r.App.Version, "1.2024.100", ua)
		}
		if !r.IsInApp {
			t.Errorf("IsInApp = false, want true\n  UA: %s", ua)
		}
		if r.IsBot {
			t.Errorf("IsBot = true, want false (app, not a bot)\n  UA: %s", ua)
		}
	})
}
