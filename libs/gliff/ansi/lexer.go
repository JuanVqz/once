package ansi

import "strings"

// Token represents a lexical token in an ANSI stream.
type Token struct {
	Type   TokenType
	Text   string // The raw text of the token
	Params string // CSI: parameter bytes between ESC[ and final byte
	Final  byte   // CSI: command byte; ESC: second byte
}

// TokenType identifies the kind of token.
type TokenType int

const (
	TextToken TokenType = iota
	CSIToken
	ESCToken
	EOFToken
)

// Lexer tokenizes a string containing ANSI escape sequences.
type Lexer struct {
	input string
	pos   int
}

// NewLexer creates a new lexer for the given input.
func NewLexer(input string) Lexer {
	return Lexer{input: input}
}

// Next returns the next token from the input.
func (l *Lexer) Next() Token {
	if l.pos >= len(l.input) {
		return Token{Type: EOFToken}
	}

	ch := l.input[l.pos]

	if ch == '\x1b' {
		return l.readEscape()
	}

	return l.readText()
}

// ParseCSIParams parses semicolon/colon-separated integers from a CSI parameter string.
// If buf is provided, it is reused to avoid allocation; results may exceed buf's capacity
// via append if needed.
func ParseCSIParams(s string, buf []int) []int {
	params := buf[:0]

	if s == "" {
		return append(params, 0)
	}

	var current int
	hasCurrent := false

	for i := range len(s) {
		b := s[i]
		if b >= '0' && b <= '9' {
			current = current*10 + int(b-'0')
			hasCurrent = true
		} else if b == ';' || b == ':' {
			if hasCurrent {
				params = append(params, current)
			} else {
				params = append(params, 0)
			}
			current = 0
			hasCurrent = false
		}
	}

	if hasCurrent {
		params = append(params, current)
	} else if len(params) == 0 {
		params = append(params, 0)
	}

	return params
}

// ParseCSIParam parses a single integer from a CSI parameter string.
func ParseCSIParam(s string) int {
	n := 0
	for i := range len(s) {
		b := s[i]
		if b >= '0' && b <= '9' {
			n = n*10 + int(b-'0')
		}
	}
	return n
}

// Private

func (l *Lexer) readText() Token {
	start := l.pos
	if i := strings.IndexByte(l.input[l.pos:], '\x1b'); i >= 0 {
		l.pos += i
	} else {
		l.pos = len(l.input)
	}
	return Token{
		Type: TextToken,
		Text: l.input[start:l.pos],
	}
}

func (l *Lexer) readEscape() Token {
	start := l.pos
	l.pos++ // consume ESC

	if l.pos >= len(l.input) {
		return Token{
			Type: TextToken,
			Text: l.input[start:l.pos],
		}
	}

	if l.input[l.pos] == '[' {
		return l.readCSI(start)
	}

	// Other escape sequence (ESC followed by single char)
	final := l.input[l.pos]
	l.pos++
	return Token{
		Type:  ESCToken,
		Text:  l.input[start:l.pos],
		Final: final,
	}
}

func (l *Lexer) readCSI(start int) Token {
	l.pos++ // consume '['

	paramStart := l.pos

	// Read parameter bytes (0x30-0x3F) and intermediate bytes (0x20-0x2F)
	for l.pos < len(l.input) {
		b := l.input[l.pos]
		if (b >= 0x30 && b <= 0x3F) || (b >= 0x20 && b <= 0x2F) {
			l.pos++
		} else {
			break
		}
	}

	paramEnd := l.pos

	// Read final byte (0x40-0x7E)
	var final byte
	if l.pos < len(l.input) {
		b := l.input[l.pos]
		if b >= 0x40 && b <= 0x7E {
			final = b
			l.pos++
		}
	}

	return Token{
		Type:   CSIToken,
		Text:   l.input[start:l.pos],
		Params: l.input[paramStart:paramEnd],
		Final:  final,
	}
}
