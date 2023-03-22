package lexer

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"unicode/utf8"
)

type Lexer struct {
	text  string
	start int
	end   int
	line  int
}

func NewLexer(src string) Lexer {
	return Lexer{text: src, line: 1}
}

func (my *Lexer) makeError(format string, args ...interface{}) (tok Token, err error) {
	err = errors.New(fmt.Sprintf(format, args...))
	return
}

func (my *Lexer) makeToken(kind TokenKind, values ...string) (Token, error) {
	var value string
	if len(values) > 0 {
		value = values[0]
	} else {
		switch kind {
		case EOF, INT, FLOAT, NAME, STRING, COMMENT:
			value = my.text[my.start:my.end]
		default:
			value = kind.String()
		}
	}
	return Token{Kind: kind, Value: value}, nil
}

func (my *Lexer) trim() {
	for my.end < len(my.text) {
		switch my.text[my.end] {
		case ' ', ',', '\t':
			my.end++
		case '\n':
			my.end++
			my.line++
		case '\r':
			my.end++
			if my.end < len(my.text) && my.text[my.end] == '\n' {
				my.end++
				my.line++
			}
		case 0xef:
			if my.end+2 < len(my.text) && my.text[my.end+1] == 0xbb && my.text[my.end+2] == 0xbf {
				my.end += 3
			} else {
				return
			}
		default:
			return
		}
	}
}

func (my *Lexer) peek() (rune, int) {
	return utf8.DecodeRuneInString(my.text[my.end:])
}

func (my *Lexer) ReadToken() (Token, error) {
	my.trim()
	my.start = my.end
	if my.end >= len(my.text) {
		return my.makeToken(EOF)
	}
	r := my.text[my.start]
	my.end++
	switch {
	case r == '!':
		return my.makeToken(BANG)
	case r == '$':
		return my.makeToken(DOLLAR)
	case r == '&':
		return my.makeToken(AMP)
	case r == ':':
		return my.makeToken(COLON)
	case r == '=':
		return my.makeToken(EQUALS)
	case r == '@':
		return my.makeToken(AT)
	case r == '|':
		return my.makeToken(PIPE)
	case r == '(':
		return my.makeToken(PAREN_L)
	case r == ')':
		return my.makeToken(PAREN_R)
	case r == '[':
		return my.makeToken(BRACKET_L)
	case r == ']':
		return my.makeToken(BRACKET_R)
	case r == '{':
		return my.makeToken(BRACE_L)
	case r == '}':
		return my.makeToken(BRACE_R)
	case r == '.':
		if len(my.text) > my.start+2 && my.text[my.start:my.start+3] == "..." {
			my.end += 2
			return my.makeToken(SPREAD)
		}
	case r == '#':
		// skip no error comments
		if comment, err := my.readComment(); err != nil {
			return comment, err
		}
		return my.ReadToken()
	case r == '_' || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z'):
		return my.readName()
	case r == '-' || ('0' <= r && r <= '9'):
		return my.readNumber()
	case r == '"':
		return my.readString()
	}
	return my.makeError(`Cannot parse the unexpected character "%s".`, string(r))
}

// readComment from the input
// #[\u0009\u0020-\uFFFF]*
func (my *Lexer) readComment() (Token, error) {
	for my.end < len(my.text) {
		r, w := my.peek()
		// SourceCharacter but not LineTerminator
		if r > 0x001f || r == '\t' {
			my.end += w
		} else {
			break
		}
	}
	return my.makeToken(COMMENT)
}

// readName from the input
// [_A-Za-z][_0-9A-Za-z]*
func (my *Lexer) readName() (Token, error) {
	for my.end < len(my.text) {
		r, w := my.peek()
		if (r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_' {
			my.end += w
		} else {
			break
		}
	}
	return my.makeToken(NAME)
}

// readNumber from the input, either a float or an int depending on whether a decimal point appears.
// Int:   -?(0|[1-9][0-9]*)
// Float: -?(0|[1-9][0-9]*)(\.[0-9]+)?((E|e)(+|-)?[0-9]+)?
func (my *Lexer) readNumber() (Token, error) {
	float := false
	// backup to the first digit
	my.end--
	my.acceptByte('-')

	if my.acceptByte('0') {
		if consumed := my.acceptDigits(); consumed != 0 {
			my.end -= consumed
			return my.makeError("Invalid number, unexpected digit after 0: %s.", my.describeNext())
		}
	} else {
		if consumed := my.acceptDigits(); consumed == 0 {
			return my.makeError("Invalid number, expected digit but got: %s.", my.describeNext())
		}
	}

	if my.acceptByte('.') {
		float = true
		if consumed := my.acceptDigits(); consumed == 0 {
			return my.makeError("Invalid number, expected digit but got: %s.", my.describeNext())
		}
	}

	if my.acceptByte('e', 'E') {
		float = true
		my.acceptByte('-', '+')
		if consumed := my.acceptDigits(); consumed == 0 {
			return my.makeError("Invalid number, expected digit but got: %s.", my.describeNext())
		}
	}

	if float {
		return my.makeToken(FLOAT)
	}
	return my.makeToken(INT)
}

// acceptByte if it matches any of given bytes, returning true if it found anything
func (my *Lexer) acceptByte(bytes ...uint8) bool {
	if my.end >= len(my.text) {
		return false
	}
	for _, accepted := range bytes {
		if my.text[my.end] == accepted {
			my.end++
			return true
		}
	}
	return false
}

// acceptDigits from the input, returning the number of digits it found
func (my *Lexer) acceptDigits() int {
	consumed := 0
	for my.end < len(my.text) && my.text[my.end] >= '0' && my.text[my.end] <= '9' {
		my.end++
		consumed++
	}
	return consumed
}

// describeNext peeks at the input and returns a human-readable string. This should will alloc
// and should only be used in errors
func (my *Lexer) describeNext() string {
	if my.end < len(my.text) {
		return strconv.Quote(string(my.text[my.end]))
	}
	return EOF.String()
}

// readString from the input
// "([^"\\\u000A\u000D]|(\\(u[0-9a-fA-F]{4}|["\\/bfnrt])))*"
func (my *Lexer) readString() (Token, error) {
	size := len(my.text)
	block := false
	if size > my.start+2 && my.text[my.start:my.start+3] == `"""` {
		block = true
		// skip the opening quote
		my.start += 3
		my.end += 2
	} else {
		// skip the opening quote
		my.start++
	}
	// this buffer is lazily created only if there are escape characters.
	var buf *bytes.Buffer

	for my.end < size {
		r := my.text[my.end]
		switch {
		default:
			var char = rune(r)
			var w = 1
			// skip unicode overhead if we are in the ascii range
			if r >= 127 {
				char, w = my.peek()
			}
			my.end += w

			if buf != nil {
				buf.WriteRune(char)
			}
		case r < ' ' && r != '\t' && (!block || (r != '\n' && r != '\r')): // SourceCharacter but not LineTerminator or Space
			return my.makeError(`Invalid character within String: "\u%04d".`, r)
		case r == '"':
			if !block {
				t, err := my.makeToken(STRING)
				// the token should not include the quotes in its value, but should cover them in its position
				t.Pos.Start--
				t.Pos.End++
				if buf != nil {
					t.Value = buf.String()
				}
				// skip the close quote
				my.end++
				return t, err
			}
			if block && my.end+2 < size && my.text[my.end:my.end+3] == `"""` {
				t, err := my.makeToken(BLOCK, buf.String())

				// the token should not include the quotes in its value, but should cover them in its position
				t.Pos.Start -= 3
				t.Pos.End += 3

				// skip the close quote
				my.end += 3
				return t, err
			}
			return my.makeError(`block string cannot contain quotes.`)
		case r == '\\':
			if my.end+1 >= size {
				my.end++
				return my.makeError(`Invalid character escape sequence.`)
			}

			if buf == nil {
				buf = bytes.NewBufferString(my.text[my.start:my.end])
			}

			escape := my.text[my.end+1]

			if escape == 'u' {
				if my.end+6 >= size {
					my.end++
					return my.makeError("Invalid character escape sequence: \\%s.", my.text[my.end:])
				}

				r, ok := unHex(my.text[my.end+2 : my.end+6])
				if !ok {
					my.end++
					return my.makeError("Invalid character escape sequence: \\%s.", my.text[my.end:my.end+5])
				}
				buf.WriteRune(r)
				my.end += 6
			} else {
				switch escape {
				case '"', '/', '\\':
					buf.WriteByte(escape)
				case 'b':
					buf.WriteByte('\b')
				case 'f':
					buf.WriteByte('\f')
				case 'n':
					buf.WriteByte('\n')
				case 'r':
					buf.WriteByte('\r')
				case 't':
					buf.WriteByte('\t')
				default:
					my.end++
					return my.makeError("Invalid character escape sequence: \\%s.", string(escape))
				}
				my.end += 2
			}
		}
	}
	return my.makeError("Unterminated string.")
}

func unHex(b string) (v rune, ok bool) {
	for _, c := range b {
		v <<= 4
		switch {
		case '0' <= c && c <= '9':
			v |= c - '0'
		case 'a' <= c && c <= 'f':
			v |= c - 'a' + 10
		case 'A' <= c && c <= 'F':
			v |= c - 'A' + 10
		default:
			return 0, false
		}
	}
	return v, true
}
