package parser

import (
	"github.com/ichaly/tiny-go/core/lexer"
	"strconv"
)

type parser struct {
	lexer lexer.Lexer
	err   error

	peeked bool
	token  lexer.Token
}

func (p *parser) unexpectedError() {
	p.unexpectedToken(p.peek())
}

func (p *parser) unexpectedToken(tok lexer.Token) {
	p.error(tok, "Unexpected %s", tok.String())
}

func (p *parser) error(tok lexer.Token, format string, args ...interface{}) {
	if p.err != nil {
		return
	}
	p.err = lexer.ErrorLocf(tok.Pos.Src.Name, tok.Pos.Line, tok.Pos.Column, format, args...)
}

func (p *parser) peek() lexer.Token {
	if p.err != nil {
		return p.token
	}
	if !p.peeked {
		p.token, p.err = p.lexer.ReadToken()
		p.peeked = true
	}
	return p.token
}

func (p *parser) next() lexer.Token {
	if p.err != nil {
		return p.token
	}
	if p.peeked {
		p.peeked = false
	} else {
		p.token, p.err = p.lexer.ReadToken()
	}
	return p.token
}

func (p *parser) skip(kind lexer.Kind) bool {
	if p.err != nil {
		return false
	}

	tok := p.peek()

	if tok.Kind != kind {
		return false
	}
	p.next()
	return true
}

func (p *parser) some(start lexer.Kind, end lexer.Kind, cb func()) {
	hasDef := p.skip(start)
	if !hasDef {
		return
	}

	called := false
	for p.peek().Kind != end && p.err == nil {
		called = true
		cb()
	}

	if !called {
		p.error(p.peek(), "expected at least one definition, found %s", p.peek().Kind.String())
		return
	}

	p.next()
}

func (p *parser) expect(kind lexer.Kind) lexer.Token {
	tok := p.peek()
	if tok.Kind == kind {
		return p.next()
	}

	p.error(tok, "Expected %s, found %s", kind, tok.Kind.String())
	return tok
}

func (p *parser) expectKeyword(value string) lexer.Token {
	tok := p.peek()
	if tok.Kind == lexer.Name && tok.Value == value {
		return p.next()
	}

	p.error(tok, "Expected %s, found %s", strconv.Quote(value), tok.String())
	return tok
}
