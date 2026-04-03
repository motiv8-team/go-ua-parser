package uax

import "time"

// PreParseHookFunc is called before parsing begins.
type PreParseHookFunc func(input Input)

// PostParseHookFunc is called after parsing completes with timing info.
type PostParseHookFunc func(input Input, result Result, duration time.Duration)

// WithPreParseHook registers a function called before each parse.
func WithPreParseHook(fn PreParseHookFunc) Option {
	return func(c *parserConfig) {
		c.preParseHook = fn
	}
}

// WithPostParseHook registers a function called after each parse.
func WithPostParseHook(fn PostParseHookFunc) Option {
	return func(c *parserConfig) {
		c.postParseHook = fn
	}
}
