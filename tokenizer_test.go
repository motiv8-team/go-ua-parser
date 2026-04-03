package uax

import "testing"

func TestTokenize(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.6312.86 Safari/537.36"
	var tk tokenizer
	tk.reset(ua)
	tokens := tk.tokenize()

	found := map[string]bool{}
	for _, tok := range tokens {
		found[tok.name] = true
	}
	for _, want := range []string{"Mozilla", "AppleWebKit", "Chrome", "Safari"} {
		if !found[want] {
			t.Errorf("missing token %q in parsed tokens", want)
		}
	}
}

func TestTokenizeBot(t *testing.T) {
	ua := "Googlebot/2.1 (+http://www.google.com/bot.html)"
	var tk tokenizer
	tk.reset(ua)
	tokens := tk.tokenize()

	if len(tokens) == 0 {
		t.Fatal("expected at least one token")
	}
	if tokens[0].name != "Googlebot" {
		t.Errorf("first token = %q, want Googlebot", tokens[0].name)
	}
	if tokens[0].version != "2.1" {
		t.Errorf("version = %q, want 2.1", tokens[0].version)
	}
}

func TestTokenizeEmpty(t *testing.T) {
	var tk tokenizer
	tk.reset("")
	tokens := tk.tokenize()
	if len(tokens) != 0 {
		t.Errorf("expected 0 tokens for empty UA, got %d", len(tokens))
	}
}

func BenchmarkTokenize(b *testing.B) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.6312.86 Safari/537.36"
	var tk tokenizer
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tk.reset(ua)
		_ = tk.tokenize()
	}
}
