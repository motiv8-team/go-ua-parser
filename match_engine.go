package uax

func (p *Parser) matchEngine(tokens []token, out *Result) {
	if out.Engine.Name != "" {
		// Already set by browser match — just fill version.
		for i := range tokens {
			if tokens[i].name == "AppleWebKit" && out.Engine.Name == "WebKit" {
				out.Engine.Version = tokens[i].version
				return
			}
			if tokens[i].name == "AppleWebKit" && out.Engine.Name == "Blink" {
				out.Engine.Version = tokens[i].version
				return
			}
		}
		return
	}
	for i := range tokens {
		tk := &tokens[i]
		if tk.name == "" {
			continue
		}
		if r, ok := p.engines.lookup(tk.name); ok {
			out.Engine.Name = r.engineName
			out.Engine.Version = tk.version
			return
		}
	}
}
