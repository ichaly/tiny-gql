package core

type OperationType string

const (
	Query        OperationType = "query"
	Mutation     OperationType = "mutation"
	Subscription OperationType = "subscription"
)

type OperationDefinition struct {
	OperationType       OperationType
	Name                string
	VariablesDefinition []*VariableDefinition
	Directives          []*Directive
	SelectionSet        []Selection
}

type Selection interface {
	isSelection()
}

type VariableDefinition struct {
	Variable   string
	Directives []*Directive
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
