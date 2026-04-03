package uax

// Input holds everything needed to parse a user agent: the raw string
// and optional Client Hints headers for improved accuracy.
type Input struct {
	UAString    string      `json:"uaString"`
	ClientHints ClientHints `json:"clientHints"`
}

// HasClientHints reports whether any Client Hints data is present.
func (i *Input) HasClientHints() bool {
	ch := &i.ClientHints
	return ch.UA != "" || ch.Platform != "" || ch.PlatformVersion != "" ||
		ch.Arch != "" || ch.Model != "" || ch.FullVersion != "" ||
		ch.FullVersionList != ""
}
