package parser

import (
	"github.com/ichaly/tiny-go/core/lexer"
	"github.com/ichaly/tiny-go/core/schema"
	"strconv"
)

type parser struct {
	lexer lexer.Lexer
	err   error

	peeked bool
	token  lexer.Token
}

func ParseQuery(src *lexer.Input) (*QueryDocument, error) {
	l := lexer.New(src)
	p := parser{lexer: l}
	doc := QueryDocument{}
	for p.peek().Kind != lexer.EOF {
		if p.err != nil {
			break
		}
		switch p.peek().Kind {
		case lexer.Comment, lexer.BlockString:
			p.next()
			continue
		case lexer.Name:
			switch p.peek().Value {
			case "query", "mutation", "subscription":
				doc.Operations = append(doc.Operations, p.parseOperationDefinition())
			case "fragment":
				doc.Fragments = append(doc.Fragments, p.parseFragmentDefinition())
			default:
				p.unexpectedError()
			}
		case lexer.BraceL:
			doc.Operations = append(doc.Operations, p.parseOperationDefinition())
		default:
			p.unexpectedError()
		}
	}
	return &doc, p.err
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

func (p *parser) skip(kind lexer.Type) bool {
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

func (p *parser) some(start lexer.Type, end lexer.Type, cb func()) {
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

func (p *parser) expect(kind lexer.Type) lexer.Token {
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

func (p *parser) parseOperationDefinition() *OperationDefinition {
	if p.peek().Kind == lexer.BraceL {
		return &OperationDefinition{
			OperationType: Query,
			SelectionSet:  p.parseRequiredSelectionSet(),
		}
	}

	var od OperationDefinition
	od.OperationType = p.parseOperationType()

	if p.peek().Kind == lexer.Name {
		od.Name = p.next().Value
	}

	od.VariablesDefinitions = p.parseVariableDefinitions()
	od.Directives = p.parseDirectives(false)
	od.SelectionSet = p.parseRequiredSelectionSet()

	return &od
}

func (p *parser) parseOperationType() OperationType {
	tok := p.next()
	switch tok.Value {
	case "query":
		return Query
	case "mutation":
		return Mutation
	case "subscription":
		return Subscription
	}
	p.unexpectedToken(tok)
	return ""
}

func (p *parser) parseRequiredSelectionSet() []Selection {
	if p.peek().Kind != lexer.BraceL {
		p.error(p.peek(), "Expected %s, found %s", lexer.BraceL, p.peek().Kind.String())
		return nil
	}

	var selections []Selection
	p.some(lexer.BraceL, lexer.BraceR, func() {
		selections = append(selections, p.parseSelection())
	})

	return selections
}

func (p *parser) parseSelection() Selection {
	if p.peek().Kind == lexer.Spread {
		return p.parseFragment()
	}
	return p.parseField()
}

func (p *parser) parseName() string {
	token := p.expect(lexer.Name)

	return token.Value
}

func (p *parser) parseField() *Field {
	var field Field
	field.Alias = p.parseName()

	if p.skip(lexer.Colon) {
		field.Name = p.parseName()
	} else {
		field.Name = field.Alias
	}

	field.Arguments = p.parseArguments(false)
	field.Directives = p.parseDirectives(false)
	if p.peek().Kind == lexer.BraceL {
		field.SelectionSet = p.parseOptionalSelectionSet()
	}

	return &field
}

func (p *parser) parseOptionalSelectionSet() []Selection {
	var selections []Selection
	p.some(lexer.BraceL, lexer.BraceR, func() {
		selections = append(selections, p.parseSelection())
	})

	return selections
}

func (p *parser) parseArguments(isConst bool) []*Argument {
	var arguments []*Argument
	p.some(lexer.ParenL, lexer.ParenR, func() {
		arguments = append(arguments, p.parseArgument(isConst))
	})

	return arguments
}

func (p *parser) parseArgument(isConst bool) *Argument {
	arg := Argument{}
	arg.Name = p.parseName()
	p.expect(lexer.Colon)

	arg.Value = p.parseValueLiteral(isConst)
	return &arg
}

func (p *parser) parseValueLiteral(isConst bool) *Value {
	token := p.peek()

	var kind ValueKind
	switch token.Kind {
	case lexer.BracketL:
		return p.parseList(isConst)
	case lexer.BraceL:
		return p.parseObject(isConst)
	case lexer.Dollar:
		if isConst {
			p.unexpectedError()
			return nil
		}
		return &Value{Raw: p.parseVariable(), Kind: Variable}
	case lexer.Int:
		kind = IntValue
	case lexer.Float:
		kind = FloatValue
	case lexer.String:
		kind = StringValue
	case lexer.BlockString:
		kind = BlockValue
	case lexer.Name:
		switch token.Value {
		case "true", "false":
			kind = BooleanValue
		case "null":
			kind = NullValue
		default:
			kind = EnumValue
		}
	default:
		p.unexpectedError()
		return nil
	}

	p.next()

	return &Value{Raw: token.Value, Kind: kind}
}

func (p *parser) parseList(isConst bool) *Value {
	var values []*Value
	p.some(lexer.BracketL, lexer.BracketR, func() {
		values = append(values, p.parseValueLiteral(isConst))
	})
	return &Value{Children: values, Kind: ListValue}
}

func (p *parser) parseObject(isConst bool) *Value {
	var fields []*Value
	p.some(lexer.BraceL, lexer.BraceR, func() {
		fields = append(fields, p.parseObjectField(isConst))
	})

	return &Value{Children: fields, Kind: ObjectValue}
}

func (p *parser) parseObjectField(isConst bool) *Value {
	field := Value{}
	field.Name = p.parseName()

	p.expect(lexer.Colon)

	field.Children = append(field.Children, p.parseValueLiteral(isConst))
	return &field
}

func (p *parser) parseVariable() string {
	p.expect(lexer.Dollar)
	return p.parseName()
}

func (p *parser) parseFragment() Selection {
	p.expect(lexer.Spread)

	if peek := p.peek(); peek.Kind == lexer.Name && peek.Value != "on" {
		return &FragmentSpread{
			Name:       p.parseFragmentName(),
			Directives: p.parseDirectives(false),
		}
	}

	var def InlineFragment
	if p.peek().Value == "on" {
		p.next() // "on"

		def.TypeCondition = p.parseName()
	}

	def.Directives = p.parseDirectives(false)
	def.SelectionSet = p.parseRequiredSelectionSet()
	return &def
}

func (p *parser) parseFragmentName() string {
	if p.peek().Value == "on" {
		p.unexpectedError()
		return ""
	}

	return p.parseName()
}

func (p *parser) parseFragmentDefinition() *FragmentDefinition {
	var def FragmentDefinition
	p.expectKeyword("fragment")

	def.Name = p.parseFragmentName()
	//def.VariableDefinition = p.parseVariableDefinitions()

	p.expectKeyword("on")

	def.TypeCondition = p.parseName()
	def.Directives = p.parseDirectives(false)
	def.SelectionSet = p.parseRequiredSelectionSet()
	return &def
}

func (p *parser) parseVariableDefinitions() []*VariableDefinition {
	var defs []*VariableDefinition
	p.some(lexer.ParenL, lexer.ParenR, func() {
		defs = append(defs, p.parseVariableDefinition())
	})
	return defs
}

func (p *parser) parseVariableDefinition() *VariableDefinition {
	var def VariableDefinition
	def.Variable = p.parseVariable()

	p.expect(lexer.Colon)

	def.Type = p.parseTypeReference()

	if p.skip(lexer.Equals) {
		def.DefaultValue = p.parseValueLiteral(true)
	}

	def.Directives = p.parseDirectives(false)

	return &def
}

func (p *parser) parseDirectives(isConst bool) []*Directive {
	var directives []*Directive

	for p.peek().Kind == lexer.At {
		if p.err != nil {
			break
		}
		directives = append(directives, p.parseDirective(isConst))
	}
	return directives
}

func (p *parser) parseDirective(isConst bool) *Directive {
	p.expect(lexer.At)

	return &Directive{
		Name:      p.parseName(),
		Arguments: p.parseArguments(isConst),
	}
}

func (p *parser) parseTypeReference() *schema.Type {
	var typ schema.Type
	if p.skip(lexer.BracketL) {
		typ.Elem = p.parseTypeReference()
		p.expect(lexer.BracketR)
	} else {
		typ.Name = p.parseName()
	}
	if p.skip(lexer.Bang) {
		typ.NonNull = true
	}
	return &typ
}
