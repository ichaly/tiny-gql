package core

type parser struct {
	lexer lexer
	err   error
}

func ParseQuery(source []byte) (*QueryDocument, error) {
	l, e := newLexer(source)
	if e != nil {
		return nil, e
	}
	p := parser{lexer: l}
	return p.parseQueryDocument(), p.err
}
