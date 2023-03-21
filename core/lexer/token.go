package lexer

import "strconv"

type TokenKind string

const (
	EOF       TokenKind = "<EOF>"
	BANG                = "!"
	DOLLAR              = "$"
	AMP                 = "&"
	SPREAD              = "..."
	COLON               = ":"
	EQUALS              = "="
	AT                  = "@"
	PAREN_L             = "("
	PAREN_R             = ")"
	BRACKET_L           = "["
	BRACKET_R           = "]"
	BRACE_L             = "{"
	BRACE_R             = "}"
	PIPE                = "|"
	NAME                = "Name"
	INT                 = "Int"
	FLOAT               = "Float"
	BLOCK               = "Block"
	STRING              = "String"
	COMMENT             = "Comment"
)

func (my TokenKind) String() string {
	return string(my)
}

type Token struct {
	Kind  TokenKind // The token type.
	Value string    // The literal value consumed.
	Pos   Position  // The position belonging to the token.
}

func (t Token) String() string {
	if t.Value != "" {
		return t.Kind.String() + " " + strconv.Quote(t.Value)
	}
	return t.Kind.String()
}

type Position struct {
	Start int // The starting position, in bytes, of this token in the input.
	End   int // The end position, in bytes, of this token in the input.
	Line  int // The line number at the start of this item.
}
