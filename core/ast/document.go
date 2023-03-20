package ast

type OperationType string

const (
	Query        OperationType = "query"
	Mutation     OperationType = "mutation"
	Subscription OperationType = "subscription"
)

type QueryDocument struct {
	Operations []*OperationDefinition
	Fragments  []*FragmentDefinition
}

type OperationDefinition struct {
	OperationType        OperationType
	Name                 string
	VariablesDefinitions []*VariableDefinition
	Directives           []*Directive
	SelectionSet         []Selection
}

type Selection interface {
	isSelection()
}

type VariableDefinition struct {
	Variable     string
	Type         *Type
	Directives   []*Directive
	DefaultValue *Value
}

type Directive struct {
	Name      string
	Arguments []*Argument
}

type Argument struct {
	Name  string
	Value *Value
}

type ValueKind int

const (
	Variable ValueKind = iota
	IntValue
	FloatValue
	BlockValue
	StringValue
	BooleanValue
	NullValue
	EnumValue
	ListValue
	ObjectValue
)

type Value struct {
	Raw      string
	Name     string
	Children []*Value
	Kind     ValueKind
}

type Field struct {
	Name         string
	Alias        string
	Arguments    []*Argument
	Directives   []*Directive
	SelectionSet []Selection
}

type FragmentSpread struct {
	Name       string
	Directives []*Directive
}

type InlineFragment struct {
	TypeCondition string
	Directives    []*Directive
	SelectionSet  []Selection
}

type FragmentDefinition struct {
	Name          string
	TypeCondition string
	Directives    []*Directive
	SelectionSet  []Selection
}

func (*Field) isSelection()          {}
func (*FragmentSpread) isSelection() {}
func (*InlineFragment) isSelection() {}

type SchemaDocument struct {
	Description string
	Schema      []*SchemaDocument
	Directives  []*Directive
	Definitions []*Definition
}
type Kind string

const (
	SCALAR       Kind = "SCALAR"
	OBJECT       Kind = "OBJECT"
	INTERFACE    Kind = "INTERFACE"
	UNION        Kind = "UNION"
	ENUM         Kind = "ENUM"
	INPUT_OBJECT Kind = "INPUT_OBJECT"
	LIST         Kind = "LIST"
	NON_NULL     Kind = "NON_NULL"
)

type Definition struct {
	Kind        Kind
	Description string
	Name        string
	Directives  []*Directive
	Interfaces  []string               // object and input object
	Fields      []*FieldDefinition     // object and input object
	Types       []string               // union
	EnumValues  []*EnumValueDefinition // enum
}

type FieldDefinition struct {
	Description  string
	Name         string
	Arguments    []*ArgumentDefinition // only for objects
	DefaultValue *Value                // only for input objects
	Type         *Type
	Directives   []*Directive
}

type ArgumentDefinition struct {
	Description  string
	Name         string
	DefaultValue *Value
	Type         *Type
	Directives   []*Directive
}

type EnumValueDefinition struct {
	Description string
	Name        string
	Directives  []*Directive
}
