package uax

import "strings"

func (p *Parser) matchBot(tokens []token, ua string, out *Result) {
	if !p.cfg.enableBotDetection {
		return
	}
	// Try ruleTable first (exact name match against tokens)
	for i := range tokens {
		tk := &tokens[i]
		if tk.name != "" {
			if r, ok := p.bots.lookup(tk.name); ok {
				out.Bot = Bot{
					IsBot:      true,
					Name:       r.botName,
					Class:      r.botClass,
					Vendor:     r.botVendor,
					IsVerified: r.botIsVerified,
					Version:    tk.version,
				}
				out.IsBot = true
				return
			}
		}
	}
	// Also check comments for bot names (e.g. "compatible; Googlebot/2.1; ...")
	for i := range tokens {
		if tokens[i].comment != "" {
			parts := strings.Split(tokens[i].comment, ";")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				// Extract name/version from part
				name := part
				ver := ""
				if slashIdx := strings.IndexByte(part, '/'); slashIdx >= 0 {
					name = part[:slashIdx]
					ver = part[slashIdx+1:]
				}
				if r, ok := p.bots.lookup(name); ok {
					out.Bot = Bot{
						IsBot:      true,
						Name:       r.botName,
						Class:      r.botClass,
						Vendor:     r.botVendor,
						IsVerified: r.botIsVerified,
						Version:    ver,
					}
					out.IsBot = true
					return
				}
			}
		}
	}
	// Heuristic fallback
	bot := detectBotHeuristic(ua)
	if bot.IsBot {
		out.Bot = bot
		out.IsBot = true
	}
}

// containsFold checks if s contains substr (case-insensitive) without allocating.
func containsFold(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if hasPrefixFold(s[i:], substr) {
			return true
		}
	}
	return false
}

func hasPrefixFold(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		a, b := s[i], prefix[i]
		if a >= 'A' && a <= 'Z' {
			a += 'a' - 'A'
		}
		if b >= 'A' && b <= 'Z' {
			b += 'a' - 'A'
		}
		if a != b {
			return false
		}
	}
	return true
}

func detectBotHeuristic(ua string) Bot {
	// Common bot indicators (case-insensitive, zero-alloc)
	botPatterns := []string{
		"bot/", "bot;", "bot ", "crawler", "spider", "scraper",
		"archiver", "transcoder", "fetcher", "validator",
		"checker", "extractor", "monitoring", "analyzer",
	}
	for _, pat := range botPatterns {
		if containsFold(ua, pat) {
			return Bot{IsBot: true, Class: BotOther, Confidence: 0.7}
		}
	}

	// URL in UA is a strong bot signal
	if strings.Contains(ua, "+http://") || strings.Contains(ua, "+https://") {
		return Bot{IsBot: true, Class: BotOther, Confidence: 0.8}
	}

	// Known CLI tools (case-insensitive prefix check)
	cliTools := []string{"curl/", "wget/", "python-requests/", "httpie/", "java/", "go-http-client/", "node-fetch/", "axios/", "libwww-perl/", "php/"}
	for _, tool := range cliTools {
		if hasPrefixFold(ua, tool) {
			name := ua[:strings.IndexByte(ua, '/')]
			return Bot{IsBot: true, Class: BotOther, Name: name, Confidence: 0.9}
		}
	}

	return Bot{}
}
