package uax

import "strings"

func (p *Parser) matchCPU(tokens []token, out *Result) {
	for i := range tokens {
		if tokens[i].comment != "" {
			cpu := parseCPUFromComment(tokens[i].comment)
			if cpu.Architecture != "" {
				out.CPU = cpu
				return
			}
		}
	}
}

func parseCPUFromComment(comment string) CPU {
	lower := strings.ToLower(comment)
	switch {
	case strings.Contains(lower, "x86_64") || strings.Contains(lower, "x86-64") ||
		strings.Contains(lower, "x64") || strings.Contains(lower, "win64") ||
		strings.Contains(lower, "amd64"):
		return CPU{Architecture: "x86_64", Bits: 64}
	case strings.Contains(lower, "aarch64") || strings.Contains(lower, "arm64"):
		return CPU{Architecture: "arm64", Bits: 64}
	case strings.Contains(lower, "i686") || strings.Contains(lower, "i386") || strings.Contains(lower, "x86"):
		return CPU{Architecture: "x86", Bits: 32}
	case strings.Contains(lower, "armv7") || strings.Contains(lower, "armv6") || strings.Contains(lower, "arm"):
		return CPU{Architecture: "arm", Bits: 32}
	case strings.Contains(lower, "android"):
		// Most Android devices are arm64 these days
		return CPU{Architecture: "arm64", Bits: 64}
	}
	return CPU{}
}
