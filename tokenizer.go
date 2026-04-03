package uax

// token represents a product or comment token extracted from a UA string.
type token struct {
	name    string // product name, e.g. "Chrome"
	version string // product version, e.g. "123.0.6312.86"
	comment string // parenthesized comment content
}

const maxTokens = 24

// tokenizer is a zero-allocation UA string scanner.
type tokenizer struct {
	ua    string
	buf   [maxTokens]token
	count int
}

func (t *tokenizer) reset(ua string) {
	t.ua = ua
	t.count = 0
}

func (t *tokenizer) tokenize() []token {
	ua := t.ua
	i := 0
	n := len(ua)

	for i < n && t.count < maxTokens {
		for i < n && ua[i] == ' ' {
			i++
		}
		if i >= n {
			break
		}

		if ua[i] == '(' {
			j := i + 1
			depth := 1
			for j < n && depth > 0 {
				if ua[j] == '(' {
					depth++
				} else if ua[j] == ')' {
					depth--
				}
				j++
			}
			// If depth never reached 0, the comment is unterminated.
			// j == n in that case; use j as the end so the slice is valid.
			end := j - 1
			if end < i+1 {
				end = i + 1
			}
			comment := ua[i+1 : end]
			t.parseComment(comment)
			i = j
			continue
		}

		nameStart := i
		for i < n && ua[i] != '/' && ua[i] != ' ' && ua[i] != '(' {
			i++
		}
		name := ua[nameStart:i]

		var ver string
		if i < n && ua[i] == '/' {
			i++
			verStart := i
			for i < n && ua[i] != ' ' && ua[i] != '(' {
				i++
			}
			ver = ua[verStart:i]
		}

		if name != "" {
			t.buf[t.count] = token{name: name, version: ver}
			t.count++
		}
	}

	return t.buf[:t.count]
}

func (t *tokenizer) parseComment(comment string) {
	if t.count >= maxTokens || comment == "" {
		return
	}
	t.buf[t.count] = token{comment: comment}
	t.count++
}
