package core

import (
	"encoding/json"
)

//
// Types represent:
// https://spec.graphql.org/draft/#sec-Schema-Introspection
//

type __Schema struct {
	Types            []__Type      `json:"types,omitempty"`
	Directives       []__Directive `json:"directives,omitempty"`
	Description      string        `json:"description,omitempty"`
	QueryType        __Type        `json:"queryType,omitempty"`
	MutationType     *__Type       `json:"mutationType,omitempty"`
	SubscriptionType *__Type       `json:"subscriptionType,omitempty"`

	conf       Config
	info       *DBInfo
	types      map[string]bool
	directives map[string]bool
}

type __Type struct {
	Kind        __TypeKind `json:"kind,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	// must be non-null for OBJECT and INTERFACE, otherwise null.
	Fields []__Field `json:"fields,omitempty"`
	// must be non-null for OBJECT and INTERFACE, otherwise null.
	Interfaces []__Type `json:"interfaces,omitempty"`
	// must be non-null for INTERFACE and UNION, otherwise null.
	PossibleTypes []__Type `json:"possibleTypes,omitempty"`
	// must be non-null for ENUM, otherwise null.
	EnumValues []__Field `json:"enumValues,omitempty"`
	// must be non-null for INPUT_OBJECT, otherwise null.
	//InputFields func(isDeprecatedArgs) []__Field `json:"-"`
	InputFields []__Field `json:"inputFields,omitempty"`
	// must be non-null for NON_NULL and LIST, otherwise null.
	OfType *__Type `json:"ofType,omitempty"`
	// may be non-null for custom SCALAR, otherwise null.
	SpecifiedByURL string `json:"specifiedByUrl,omitempty"`
}

type __Field struct {
	Name              string    `json:"name,omitempty"`
	Description       string    `json:"description,omitempty"`
	Args              []__Field `json:"args,omitempty"`
	Type              *__Type   `json:"type,omitempty"`
	IsDeprecated      bool      `json:"isDeprecated,omitempty"`
	DeprecationReason string    `json:"deprecationReason,omitempty"`
	DefaultValue      string    `json:"defaultValue,omitempty"`
}

type __Directive struct {
	Name         string                `json:"name,omitempty"`
	Description  string                `json:"description,omitempty"`
	Locations    []__DirectiveLocation `json:"locations,omitempty"`
	Args         []__Field             `json:"args,omitempty"`
	IsRepeatable bool                  `json:"isRepeatable,omitempty"`
}

type __TypeKind string

const (
	TK_SCALAR       __TypeKind = "SCALAR"
	TK_OBJECT       __TypeKind = "OBJECT"
	TK_INTERFACE    __TypeKind = "INTERFACE"
	TK_UNION        __TypeKind = "UNION"
	TK_ENUM         __TypeKind = "ENUM"
	TK_INPUT_OBJECT __TypeKind = "INPUT_OBJECT"
	TK_LIST         __TypeKind = "LIST"
	TK_NON_NULL     __TypeKind = "NON_NULL"
)

type __DirectiveLocation string

const (
	DL_QUERY                  __DirectiveLocation = "QUERY"
	DL_MUTATION               __DirectiveLocation = "MUTATION"
	DL_SUBSCRIPTION           __DirectiveLocation = "SUBSCRIPTION"
	DL_FIELD                  __DirectiveLocation = "FIELD"
	DL_FRAGMENT_DEFINITION    __DirectiveLocation = "FRAGMENT_DEFINITION"
	DL_FRAGMENT_SPREAD        __DirectiveLocation = "FRAGMENT_SPREAD"
	DL_INLINE_FRAGMENT        __DirectiveLocation = "INLINE_FRAGMENT"
	DL_VARIABLE_DEFINITION    __DirectiveLocation = "VARIABLE_DEFINITION"
	DL_SCHEMA                 __DirectiveLocation = "SCHEMA"
	DL_SCALAR                 __DirectiveLocation = "SCALAR"
	DL_OBJECT                 __DirectiveLocation = "OBJECT"
	DL_FIELD_DEFINITION       __DirectiveLocation = "FIELD_DEFINITION"
	DL_ARGUMENT_DEFINITION    __DirectiveLocation = "ARGUMENT_DEFINITION"
	DL_INTERFACE              __DirectiveLocation = "INTERFACE"
	DL_UNION                  __DirectiveLocation = "UNION"
	DL_ENUM                   __DirectiveLocation = "ENUM"
	DL_ENUM_VALUE             __DirectiveLocation = "ENUM_VALUE"
	DL_INPUT_OBJECT           __DirectiveLocation = "INPUT_OBJECT"
	DL_INPUT_FIELD_DEFINITION __DirectiveLocation = "INPUT_FIELD_DEFINITION"
)

func (my *__Schema) addType(t __Type) {
	if _, ok := my.types[t.Name]; ok {
		return
	}
	my.Types = append(my.Types, t)
}

func (my *__Schema) addExpression(exps []__Field, name string, sub __Type) {
	name = name + SUFFIX_EXP
	if _, ok := my.types[name]; ok {
		return
	}

	t := __Type{
		Kind:        TK_INPUT_OBJECT,
		Name:        name,
		InputFields: []__Field{},
	}

	for _, ex := range exps {
		if ex.Type == nil {
			ex.Type = &sub
		}
		t.InputFields = append(t.InputFields, ex)
	}
	my.addType(t)
}

func Introspection() (res json.RawMessage, err error) {
	schema := __Schema{
		QueryType:        __Type{Name: "Query"},
		MutationType:     &__Type{Name: "Mutation"},
		SubscriptionType: &__Type{Name: "Subscription"},
	}

	for _, v := range stdTypes {
		schema.addType(v)
	}

	// Expression types
	v := append(expAll, expScalar...)
	schema.addExpression(v, "ID", __Type{Kind: TK_SCALAR, Name: "ID"})
	schema.addExpression(v, "Int", __Type{Kind: TK_SCALAR, Name: "Int"})
	schema.addExpression(v, "Float", __Type{Kind: TK_SCALAR, Name: "Float"})
	schema.addExpression(v, "String", __Type{Kind: TK_SCALAR, Name: "String"})
	schema.addExpression(v, "Boolean", __Type{Kind: TK_SCALAR, Name: "Boolean"})

	// ListExpression Types
	v = append(expAll, expList...)
	schema.addExpression(v, "IntList", __Type{Kind: TK_SCALAR, Name: "Int"})
	schema.addExpression(v, "FloatList", __Type{Kind: TK_SCALAR, Name: "Float"})
	schema.addExpression(v, "StringList", __Type{Kind: TK_SCALAR, Name: "String"})
	schema.addExpression(v, "BooleanList", __Type{Kind: TK_SCALAR, Name: "Boolean"})

	// JsonExpression types
	v = append(expAll, expJSON...)
	schema.addExpression(v, "JSON", __Type{Kind: TK_SCALAR, Name: "String"})

	return json.Marshal(schema)
}
