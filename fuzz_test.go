package uax

import "testing"

func FuzzParseString(f *testing.F) {
	seeds := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.6312.86 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Mobile/15E148 Safari/604.1",
		"Googlebot/2.1 (+http://www.google.com/bot.html)",
		"curl/8.4.0",
		"",
		"()()()()()",
		string(make([]byte, 10000)),
	}
	for _, s := range seeds {
		f.Add(s)
	}

	p, _ := NewParser()
	f.Fuzz(func(t *testing.T, ua string) {
		r := p.ParseString(ua)
		if r.IsBot && !r.Bot.IsBot {
			t.Errorf("Result.IsBot=true but Bot.IsBot=false for UA %q", ua)
		}
		if r.IsMobile && r.IsDesktop {
			t.Errorf("both IsMobile and IsDesktop for UA %q", ua)
		}
	})
}

func FuzzDetectBot(f *testing.F) {
	f.Add("Googlebot/2.1")
	f.Add("normal browser UA")
	f.Add("")
	f.Add("bot bot bot bot")

	p, _ := NewParser()
	f.Fuzz(func(t *testing.T, ua string) {
		b := p.DetectBot(ua)
		if b.Confidence < 0 || b.Confidence > 1 {
			t.Errorf("confidence out of range: %f", b.Confidence)
		}
	})
}
