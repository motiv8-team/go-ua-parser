package uax

import (
	"regexp"
	"strings"
)

type matchType int

const (
	matchExact    matchType = iota
	matchContains
	matchPrefix
	matchRegex
)

type rule struct {
	pattern   string
	matchType matchType
	re        *regexp.Regexp

	browserName    string
	browserFamily  string
	browserChannel string
	engineName     string
	osName         string
	deviceType     string
	deviceVendor   string
	deviceModel    string
	cpuArch        string
	cpuBits        int
	appName        string
	appKind        string
	botName        string
	botClass       BotClass
	botVendor      string
	botIsVerified  bool

	versionFrom int
}

func (r *rule) matches(candidate string) bool {
	switch r.matchType {
	case matchExact:
		return candidate == r.pattern
	case matchContains:
		return strings.Contains(candidate, r.pattern)
	case matchPrefix:
		return strings.HasPrefix(candidate, r.pattern)
	case matchRegex:
		if r.re != nil {
			return r.re.MatchString(candidate)
		}
		return false
	}
	return false
}

type ruleTable struct {
	rules    []rule
	fastTrie *trie
	nonExactIdx []int
}

func (rt *ruleTable) buildTrie() {
	rt.fastTrie = newTrie()
	rt.nonExactIdx = nil
	for i := range rt.rules {
		r := &rt.rules[i]
		switch r.matchType {
		case matchExact:
			rt.fastTrie.insert(r.pattern, i+1)
		default:
			rt.nonExactIdx = append(rt.nonExactIdx, i)
		}
	}
}

func (rt *ruleTable) lookup(candidate string) (*rule, bool) {
	if rt.fastTrie != nil {
		if id, ok := rt.fastTrie.match(candidate); ok {
			return &rt.rules[id-1], true
		}
	}
	for _, idx := range rt.nonExactIdx {
		if rt.rules[idx].matches(candidate) {
			return &rt.rules[idx], true
		}
	}
	return nil, false
}
