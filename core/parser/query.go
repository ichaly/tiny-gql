package parser

import (
	"github.com/ichaly/tiny-go/core/ast"
	"github.com/ichaly/tiny-go/core/lexer"
)

func ParseQuery(src *lexer.Input) (*ast.QueryDocument, error) {
	l := lexer.NewLexer(src)
	p := parser{lexer: l}
	doc := ast.QueryDocument{}
	for p.peek().Kind != lexer.EOF {
		if p.err != nil {
			break
		}
		t := p.peek()
		switch {
		case t.Value == "fragment":
			doc.Fragments = append(doc.Fragments, p.parseFragmentDefinition())
		case t.Kind == lexer.BraceL || t.Value == "query" || t.Value == "mutation" || t.Value == "subscription":
			doc.Operations = append(doc.Operations, p.parseOperationDefinition())
		default:
			p.unexpectedError()
		}
	}
	return &doc, p.err
}

func (p *parser) parseOperationDefinition() *ast.OperationDefinition {
	if p.peek().Kind == lexer.BraceL {
		return &ast.OperationDefinition{
			OperationType: ast.Query,
			SelectionSet:  p.parseRequiredSelectionSet(),
		}
	}

	var od ast.OperationDefinition
	od.OperationType = p.parseOperationType()

	if p.peek().Kind == lexer.Name {
		od.Name = p.next().Value
	}

	od.VariablesDefinitions = p.parseVariableDefinitions()
	od.Directives = p.parseDirectives(false)
	od.SelectionSet = p.parseRequiredSelectionSet()

	return &od
}

func (p *parser) parseOperationType() ast.OperationType {
	tok := p.next()
	switch tok.Value {
	case "query":
		return ast.Query
	case "mutation":
		return ast.Mutation
	case "subscription":
		return ast.Subscription
	}
	p.unexpectedToken(tok)
	return ""
}

func (p *parser) parseRequiredSelectionSet() []ast.Selection {
	if p.peek().Kind != lexer.BraceL {
		p.error(p.peek(), "Expected %s, found %s", lexer.BraceL, p.peek().Kind.String())
		return nil
	}

	var selections []ast.Selection
	p.some(lexer.BraceL, lexer.BraceR, func() {
		selections = append(selections, p.parseSelection())
	})

	return selections
}

func (p *parser) parseSelection() ast.Selection {
	if p.peek().Kind == lexer.Spread {
		return p.parseFragment()
	}
	return p.parseField()
}

func (p *parser) parseName() string {
	token := p.expect(lexer.Name)

	return token.Value
}

func (p *parser) parseField() *ast.Field {
	var field ast.Field
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

func (p *parser) parseOptionalSelectionSet() []ast.Selection {
	var selections []ast.Selection
	p.some(lexer.BraceL, lexer.BraceR, func() {
		selections = append(selections, p.parseSelection())
	})

	return selections
}

func (p *parser) parseArguments(isConst bool) []*ast.Argument {
	var arguments []*ast.Argument
	p.some(lexer.ParenL, lexer.ParenR, func() {
		arguments = append(arguments, p.parseArgument(isConst))
	})

	return arguments
}

func (p *parser) parseArgument(isConst bool) *ast.Argument {
	arg := ast.Argument{}
	arg.Name = p.parseName()
	p.expect(lexer.Colon)

	arg.Value = p.parseValueLiteral(isConst)
	return &arg
}

func (p *parser) parseValueLiteral(isConst bool) *ast.Value {
	token := p.peek()

	var kind ast.ValueKind
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
		return &ast.Value{Raw: p.parseVariable(), Kind: ast.Variable}
	case lexer.Int:
		kind = ast.IntValue
	case lexer.Float:
		kind = ast.FloatValue
	case lexer.Block:
		kind = ast.BlockValue
	case lexer.String:
		kind = ast.StringValue
	case lexer.Name:
		switch token.Value {
		case "true", "false":
			kind = ast.BooleanValue
		case "null":
			kind = ast.NullValue
		default:
			kind = ast.EnumValue
		}
	default:
		p.unexpectedError()
		return nil
	}

	p.next()

	return &ast.Value{Raw: token.Value, Kind: kind}
}

func (p *parser) parseList(isConst bool) *ast.Value {
	var values []*ast.Value
	p.some(lexer.BracketL, lexer.BracketR, func() {
		values = append(values, p.parseValueLiteral(isConst))
	})
	return &ast.Value{Children: values, Kind: ast.ListValue}
}

func (p *parser) parseObject(isConst bool) *ast.Value {
	var fields []*ast.Value
	p.some(lexer.BraceL, lexer.BraceR, func() {
		fields = append(fields, p.parseObjectField(isConst))
	})

	return &ast.Value{Children: fields, Kind: ast.ObjectValue}
}

func (p *parser) parseObjectField(isConst bool) *ast.Value {
	field := ast.Value{}
	field.Name = p.parseName()

	p.expect(lexer.Colon)

	field.Children = append(field.Children, p.parseValueLiteral(isConst))
	return &field
}

func (p *parser) parseVariable() string {
	p.expect(lexer.Dollar)
	return p.parseName()
}

func (p *parser) parseFragment() ast.Selection {
	p.expect(lexer.Spread)

	if peek := p.peek(); peek.Kind == lexer.Name && peek.Value != "on" {
		return &ast.FragmentSpread{
			Name:       p.parseFragmentName(),
			Directives: p.parseDirectives(false),
		}
	}

	var def ast.InlineFragment
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

func (p *parser) parseFragmentDefinition() *ast.FragmentDefinition {
	var def ast.FragmentDefinition
	p.expectKeyword("fragment")

	def.Name = p.parseFragmentName()

	p.expectKeyword("on")

	def.TypeCondition = p.parseName()
	def.Directives = p.parseDirectives(false)
	def.SelectionSet = p.parseRequiredSelectionSet()
	return &def
}

func (p *parser) parseVariableDefinitions() []*ast.VariableDefinition {
	var defs []*ast.VariableDefinition
	p.some(lexer.ParenL, lexer.ParenR, func() {
		defs = append(defs, p.parseVariableDefinition())
	})
	return defs
}

func (p *parser) parseVariableDefinition() *ast.VariableDefinition {
	var def ast.VariableDefinition
	def.Variable = p.parseVariable()

	p.expect(lexer.Colon)

	def.Type = p.parseTypeReference()

	if p.skip(lexer.Equals) {
		def.DefaultValue = p.parseValueLiteral(true)
	}

	def.Directives = p.parseDirectives(false)

	return &def
}

func (p *parser) parseDirectives(isConst bool) []*ast.Directive {
	var directives []*ast.Directive

	for p.peek().Kind == lexer.At {
		if p.err != nil {
			break
		}
		directives = append(directives, p.parseDirective(isConst))
	}
	return directives
}

func (p *parser) parseDirective(isConst bool) *ast.Directive {
	p.expect(lexer.At)

	return &ast.Directive{
		Name:      p.parseName(),
		Arguments: p.parseArguments(isConst),
	}
}

func (p *parser) parseTypeReference() *ast.Type {
	var typ ast.Type
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
