package lexer

import (
	"bytes"
	"strconv"
	"unicode/utf8"
)

// Lexer turns graphql request and schema strings into tokens
type Lexer struct {
	*Input
	// An offset into the string in bytes
	start int
	// An offset into the string in bytes
	end int
	// the current line number
	line int
	// the current column number
	column int
}

func NewLexer(src *Input) Lexer {
	return Lexer{
		Input: src,
		line:  1,
	}
}

// take one rune from input and advance end
func (s *Lexer) peek() (rune, int) {
	return utf8.DecodeRuneInString(s.Content[s.end:])
}

func (s *Lexer) makeToken(kind Kind, values ...string) (Token, error) {
	var value string
	for _, v := range values {
		value = v
		break
	}
	if len(value) == 0 {
		switch kind {
		case EOF, Comment, Float, Int, Name, String:
			value = s.Content[s.start:s.end]
		default:
			value = kind.String()
		}
	}
	return Token{
		Kind:  kind,
		Value: value,
		Pos: Position{
			Start:  s.start,
			End:    s.end,
			Line:   s.line,
			Column: s.column,
			Src:    s.Input,
		},
	}, nil
}

func (s *Lexer) makeError(format string, args ...interface{}) (Token, error) {
	return Token{
		Kind: ERR,
		Pos: Position{
			Start:  s.start,
			End:    s.end,
			Line:   s.line,
			Column: s.column,
			Src:    s.Input,
		},
	}, ErrorLocf(s.Input.Name, s.line, s.column, format, args...)
}

// ReadToken gets the next token from the source starting at the given position.
//
// This skips over whitespace and comments until it finds the next lexable
// token, then lexes punctuators immediately or calls the appropriate helper
// function for more complicated tokens.
func (s *Lexer) ReadToken() (token Token, err error) {
	s.trim()
	s.start = s.end

	if s.end >= len(s.Content) {
		return s.makeToken(EOF)
	}
	r := s.Content[s.start]
	s.end++
	s.column++
	switch {
	case r == '!':
		return s.makeToken(Bang)
	case r == '$':
		return s.makeToken(Dollar)
	case r == '&':
		return s.makeToken(Amp)
	case r == '(':
		return s.makeToken(ParenL)
	case r == ')':
		return s.makeToken(ParenR)
	case r == '.':
		if len(s.Content) > s.start+2 && s.Content[s.start:s.start+3] == "..." {
			s.end += 2
			s.column++
			return s.makeToken(Spread)
		}
	case r == ':':
		return s.makeToken(Colon)
	case r == '=':
		return s.makeToken(Equals)
	case r == '@':
		return s.makeToken(At)
	case r == '[':
		return s.makeToken(BracketL)
	case r == ']':
		return s.makeToken(BracketR)
	case r == '{':
		return s.makeToken(BraceL)
	case r == '}':
		return s.makeToken(BraceR)
	case r == '|':
		return s.makeToken(Pipe)
	case r == '#':
		// skip no error comments
		if comment, err := s.readComment(); err != nil {
			return comment, err
		}
		return s.ReadToken()
	case r == '_' || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z'):
		return s.readName()
	case r == '-' || ('0' <= r && r <= '9'):
		return s.readNumber()
	case r == '"':
		if len(s.Content) > s.start+2 && s.Content[s.start:s.start+3] == `"""` {
			return s.readBlockString()
		}
		return s.readString()
	}

	s.end--
	s.column--
	if r < 0x0020 && r != 0x0009 && r != 0x000a && r != 0x000d {
		return s.makeError(`Cannot contain the invalid character "\u%04d"`, r)
	}

	if r == '\'' {
		return s.makeError(`Unexpected single quote character ('), did you mean to use a double quote (")?`)
	}

	return s.makeError(`Cannot parse the unexpected character "%s".`, string(r))
}

// trim reads from body starting at startPosition until it finds a non-whitespace
// or commented character, and updates the token end to include all whitespace
func (s *Lexer) trim() {
	for s.end < len(s.Content) {
		switch s.Content[s.end] {
		case '\t', ' ', ',':
			s.end++
			s.column++
		case '\n':
			s.end++
			s.line++
			s.column = 0
		case '\r':
			s.end++
			s.line++
			s.column++
			// skip the following newline if its there
			if s.end < len(s.Content) && s.Content[s.end] == '\n' {
				s.end++
				s.column = 0
			}
			// byte order mark, given trim is hot path we aren't relying on the unicode package here.
		case 0xef:
			if s.end+2 < len(s.Content) && s.Content[s.end+1] == 0xBB && s.Content[s.end+2] == 0xBF {
				s.end += 3
				s.column++
			} else {
				return
			}
		default:
			return
		}
	}
}

// readComment from the input
// #[\u0009\u0020-\uFFFF]*
func (s *Lexer) readComment() (Token, error) {
	for s.end < len(s.Content) {
		r, w := s.peek()
		// SourceCharacter but not LineTerminator
		if r > 0x001f || r == '\t' {
			s.end += w
			s.column++
		} else {
			break
		}
	}
	return s.makeToken(Comment)
}

// readNumber from the input, either a float
// or an int depending on whether a decimal point appears.
//
// Int:   -?(0|[1-9][0-9]*)
// Float: -?(0|[1-9][0-9]*)(\.[0-9]+)?((E|e)(+|-)?[0-9]+)?
func (s *Lexer) readNumber() (Token, error) {
	float := false

	// backup to the first digit
	s.end--
	s.column--
	s.acceptByte('-')

	if s.acceptByte('0') {
		if consumed := s.acceptDigits(); consumed != 0 {
			s.end -= consumed
			s.column -= consumed
			return s.makeError("Invalid number, unexpected digit after 0: %s.", s.describeNext())
		}
	} else {
		if consumed := s.acceptDigits(); consumed == 0 {
			return s.makeError("Invalid number, expected digit but got: %s.", s.describeNext())
		}
	}

	if s.acceptByte('.') {
		float = true

		if consumed := s.acceptDigits(); consumed == 0 {
			return s.makeError("Invalid number, expected digit but got: %s.", s.describeNext())
		}
	}

	if s.acceptByte('e', 'E') {
		float = true

		s.acceptByte('-', '+')

		if consumed := s.acceptDigits(); consumed == 0 {
			return s.makeError("Invalid number, expected digit but got: %s.", s.describeNext())
		}
	}

	if float {
		return s.makeToken(Float)
	}
	return s.makeToken(Int)

}

// acceptByte if it matches any of given bytes, returning true if it found anything
func (s *Lexer) acceptByte(bytes ...uint8) bool {
	if s.end >= len(s.Content) {
		return false
	}

	for _, accepted := range bytes {
		if s.Content[s.end] == accepted {
			s.end++
			s.column++
			return true
		}
	}
	return false
}

// acceptDigits from the input, returning the number of digits it found
func (s *Lexer) acceptDigits() int {
	consumed := 0
	for s.end < len(s.Content) && s.Content[s.end] >= '0' && s.Content[s.end] <= '9' {
		s.end++
		s.column++
		consumed++
	}

	return consumed
}

// describeNext peeks at the input and returns a human-readable string. This should will alloc
// and should only be used in errors
func (s *Lexer) describeNext() string {
	if s.end < len(s.Content) {
		return strconv.Quote(string(s.Content[s.end]))
	}
	return "<EOF>"
}

// readString from the input
//
// "([^"\\\u000A\u000D]|(\\(u[0-9a-fA-F]{4}|["\\/bfnrt])))*"
func (s *Lexer) readString() (Token, error) {
	size := len(s.Content)

	// this buffer is lazily created only if there are escape characters.
	var buf *bytes.Buffer

	// skip the opening quote
	s.start++
	s.column++

	for s.end < size {
		r := s.Content[s.end]
		if r == '\n' || r == '\r' {
			break
		}
		if r < 0x0020 && r != '\t' {
			return s.makeError(`Invalid character within String: "\u%04d".`, r)
		}
		switch r {
		default:
			var char = rune(r)
			var w = 1

			// skip unicode overhead if we are in the ascii range
			if r >= 127 {
				char, w = utf8.DecodeRuneInString(s.Content[s.end:])
			}
			s.end += w
			s.column++

			if buf != nil {
				buf.WriteRune(char)
			}

		case '"':
			t, err := s.makeToken(String)
			// the token should not include the quotes in its value, but should cover them in its position
			t.Pos.Start--
			t.Pos.End++

			if buf != nil {
				t.Value = buf.String()
			}

			// skip the close quote
			s.end++
			s.column++

			return t, err

		case '\\':
			if s.end+1 >= size {
				s.end++
				s.column++
				return s.makeError(`Invalid character escape sequence.`)
			}

			if buf == nil {
				buf = bytes.NewBufferString(s.Content[s.start:s.end])
			}

			escape := s.Content[s.end+1]

			if escape == 'u' {
				if s.end+6 >= size {
					s.end++
					s.column++
					return s.makeError("Invalid character escape sequence: \\%s.", s.Content[s.end:])
				}

				r, ok := unhex(s.Content[s.end+2 : s.end+6])
				if !ok {
					s.end++
					s.column++
					return s.makeError("Invalid character escape sequence: \\%s.", s.Content[s.end:s.end+5])
				}
				buf.WriteRune(r)
				s.end += 6
				s.column++
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
					s.end++
					s.column++
					return s.makeError("Invalid character escape sequence: \\%s.", string(escape))
				}
				s.end += 2
				s.column++
			}
		}
	}

	return s.makeError("Unterminated string.")
}

// readBlockString from the input
//
// """("?"?(\\"""|\\(?!=""")|[^"\\]))*"""
func (s *Lexer) readBlockString() (Token, error) {
	size := len(s.Content)

	var buf bytes.Buffer

	// skip the opening quote
	s.start += 3
	s.end += 2
	s.column += 3

	for s.end < size {
		r := s.Content[s.end]

		// Closing triple quote (""")
		if r == '"' && s.end+3 <= size && s.Content[s.end:s.end+3] == `"""` {
			t, err := s.makeToken(Block, blockStringValue(buf.String()))

			// the token should not include the quotes in its value, but should cover them in its position
			t.Pos.Start -= 3
			t.Pos.End += 3

			// skip the close quote
			s.end += 3
			s.column += 3
			return t, err
		}

		// SourceCharacter
		if r < 0x0020 && r != '\t' && r != '\n' && r != '\r' {
			return s.makeError(`Invalid character within String: "\u%04d".`, r)
		}

		if r == '\\' && s.end+4 <= size && s.Content[s.end:s.end+4] == `\"""` {
			buf.WriteString(`"""`)
			s.end += 4
		} else if r == '\r' {
			if s.end+1 < size && s.Content[s.end+1] == '\n' {
				s.end++
			}

			buf.WriteByte('\n')
			s.end++
			s.line++
			s.column = 0
		} else {
			var char = rune(r)
			var w = 1

			// skip unicode overhead if we are in the ascii range
			if r >= 127 {
				char, w = utf8.DecodeRuneInString(s.Content[s.end:])
			}
			s.end += w
			buf.WriteRune(char)
			if r == '\n' {
				s.line++
				s.column = 0
			}
		}
	}

	return s.makeError("Unterminated string.")
}

func unhex(b string) (v rune, ok bool) {
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

// readName from the input
//
// [_A-Za-z][_0-9A-Za-z]*
func (s *Lexer) readName() (Token, error) {
	for s.end < len(s.Content) {
		r, w := s.peek()
		if (r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_' {
			s.end += w
			s.column++
		} else {
			break
		}
	}
	return s.makeToken(Name)
}
