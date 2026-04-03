package uax

func (p *Parser) matchApp(tokens []token, out *Result) {
	for i := len(tokens) - 1; i >= 0; i-- {
		tk := &tokens[i]
		if tk.name == "" {
			continue
		}
		if r, ok := p.apps.lookup(tk.name); ok {
			out.App.Name = r.appName
			out.App.Version = tk.version
			out.App.Kind = r.appKind
			out.IsInApp = true
			return
		}
	}
}
