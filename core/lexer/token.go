package lexer

import (
	"encoding/json"
	"strconv"
)

// Kind represents a type of token. The types are predefined as constants.
type Kind int

const (
	ERR Kind = iota
	EOF
	Bang
	Dollar
	Amp
	ParenL
	ParenR
	Spread
	Colon
	Equals
	At
	BracketL
	BracketR
	BraceL
	BraceR
	Pipe
	Name
	Int
	Float
	Block
	String
	Comment
)

func (my Kind) String() string {
	switch my {
	case ERR:
		return "<Err>"
	case EOF:
		return "<EOF>"
	case Bang:
		return "!"
	case Dollar:
		return "$"
	case Amp:
		return "&"
	case ParenL:
		return "("
	case ParenR:
		return ")"
	case Spread:
		return "..."
	case Colon:
		return ":"
	case Equals:
		return "="
	case At:
		return "@"
	case BracketL:
		return "["
	case BracketR:
		return "]"
	case BraceL:
		return "{"
	case BraceR:
		return "}"
	case Pipe:
		return "|"
	case Name:
		return "Name"
	case Int:
		return "Int"
	case Float:
		return "Float"
	case Block:
		return "Block"
	case String:
		return "String"
	case Comment:
		return "Comment"
	}
	return "Unknown " + strconv.Itoa(int(my))
}

func (my Kind) MarshalJSON() ([]byte, error) {
	return json.Marshal(my.String())
}

type Token struct {
	Kind  Kind     // The token type.
	Value string   // The literal value consumed.
	Pos   Position // The file and line this token was read from
}

func (t Token) String() string {
	if t.Value != "" {
		return t.Kind.String() + " " + strconv.Quote(t.Value)
	}
	return t.Kind.String()
}
