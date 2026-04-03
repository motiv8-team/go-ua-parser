package uax

import "testing"

var benchUAs = map[string]string{
	"ChromeDesktop":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.6312.86 Safari/537.36",
	"SafariMobile":    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Mobile/15E148 Safari/604.1",
	"FirefoxDesktop":  "Mozilla/5.0 (X11; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0",
	"EdgeDesktop":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 Edg/123.0.2420.65",
	"Googlebot":       "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	"GPTBot":          "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; GPTBot/1.2; +https://openai.com/gptbot)",
	"Curl":            "curl/8.4.0",
	"AndroidChrome":   "Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.6312.86 Mobile Safari/537.36",
	"SamsungTV":       "Mozilla/5.0 (SMART-TV; LINUX; Tizen 7.0) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/5.0 Chrome/85.0.4183.93 TV Safari/537.36",
	"FacebookInApp":   "Mozilla/5.0 (iPhone; CPU iPhone OS 17_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 [FBAN/FBIOS;FBAV/453.0.0.44.104]",
}

func BenchmarkParseString(b *testing.B) {
	p, _ := NewParser()
	for name, ua := range benchUAs {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = p.ParseString(ua)
			}
		})
	}
}

func BenchmarkParseInto(b *testing.B) {
	p, _ := NewParser()
	for name, ua := range benchUAs {
		b.Run(name, func(b *testing.B) {
			var r Result
			input := Input{UAString: ua}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				p.ParseInto(input, &r)
			}
		})
	}
}

func BenchmarkDetectBot(b *testing.B) {
	p, _ := NewParser()
	uas := []string{
		"Googlebot/2.1 (+http://www.google.com/bot.html)",
		"Mozilla/5.0 (compatible; GPTBot/1.2; +https://openai.com/gptbot)",
		"curl/8.4.0",
		"Mozilla/5.0 Chrome/123.0 Safari/537.36",
	}
	for _, ua := range uas {
		short := ua
		if len(short) > 20 {
			short = short[:20]
		}
		b.Run(short, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = p.DetectBot(ua)
			}
		})
	}
}

func BenchmarkGlobalParse(b *testing.B) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.6312.86 Safari/537.36"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Parse(ua)
	}
}
