package core

import (
	"encoding/json"
)

//
// Types represent:
// https://spec.graphql.org/draft/#sec-Schema-Introspection
//

type __Schema struct {
	Types            func() []__Type `json:"-"`
	Description      string          `json:"description,omitempty"`
	QueryType        __Type          `json:"queryType"`
	MutationType     *__Type         `json:"mutationType,omitempty"`
	SubscriptionType *__Type         `json:"subscriptionType,omitempty"`
	Directives       []__Directive   `json:"directives"`

	JSONTypes []__Type `json:"types"`
}

type __Type struct {
	Kind        __TypeKind `json:"-"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	// must be non-null for OBJECT and INTERFACE, otherwise null.
	Fields func(isDeprecatedArgs) []__Field `json:"-"`
	// must be non-null for OBJECT and INTERFACE, otherwise null.
	Interfaces []__Type `json:"interfaces"`
	// must be non-null for INTERFACE and UNION, otherwise null.
	PossibleTypes func() []__Type `json:"possibleTypes"`
	// must be non-null for ENUM, otherwise null.
	EnumValues func(isDeprecatedArgs) []__Field `json:"-"`
	// must be non-null for INPUT_OBJECT, otherwise null.
	InputFields func(isDeprecatedArgs) []__Field `json:"-"`
	// must be non-null for NON_NULL and LIST, otherwise null.
	OfType *__Type `json:"ofType"`
	// may be non-null for custom SCALAR, otherwise null.
	SpecifiedByURL string `json:"specifiedByUrl"`

	JSONKind        string    `json:"kind"`
	JSONFields      []__Field `json:"fields"`
	JSONEnumValues  []__Field `json:"enumValues"`
	JSONInputFields []__Field `json:"inputFields"`
}

type __Field struct {
	Name              string                           `json:"name"`
	Description       string                           `json:"description,omitempty"`
	Args              func(isDeprecatedArgs) []__Field `json:"-"`
	Type              __Type                           `json:"type"`
	IsDeprecated      bool                             `json:"isDeprecated,omitempty"`
	DeprecationReason string                           `json:"deprecationReason"`
	DefaultValue      string                           `json:"defaultValue,omitempty"`

	JSONArgs []__Field `json:"args"`
}

type __Directive struct {
	Name         string                           `json:"name"`
	Description  string                           `json:"description,omitempty"`
	Locations    []__DirectiveLocation            `json:"-"`
	Args         func(isDeprecatedArgs) []__Field `json:"-"`
	IsRepeatable bool                             `json:"isRepeatable"`

	JSONLocations []string  `json:"locations"`
	JSONArgs      []__Field `json:"args"`
}

type isDeprecatedArgs struct {
	IncludeDeprecated bool `json:"includeDeprecated"`
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

func introspection() (res json.RawMessage, err error) {
	schema := __Schema{
		QueryType:        __Type{Name: "Query"},
		MutationType:     &__Type{Name: "Mutation"},
		SubscriptionType: &__Type{Name: "Subscription"},
	}
	res, err = json.Marshal(schema)
	return
}
