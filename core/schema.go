package core

//
// Types represent:
// https://spec.graphql.org/draft/#sec-Schema-Introspection
//

type __Schema struct {
	Types            func() []__Type `json:"-"`
	Description      string          `json:"description"`
	QueryType        __Type          `json:"queryType"`
	MutationType     *__Type         `json:"mutationType"`
	SubscriptionType *__Type         `json:"subscriptionType"`
	Directives       []__Directive   `json:"directives"`

	JSONTypes []__Type `json:"types"`
}

type __Type struct {
	Kind        __TypeKind `json:"-"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	// must be non-null for OBJECT and INTERFACE, otherwise null.
	Fields func(isDeprecatedArgs) []__Field `json:"-"`
	// must be non-null for OBJECT and INTERFACE, otherwise null.
	Interfaces []__Type `json:"interfaces"`
	// must be non-null for INTERFACE and UNION, otherwise null.
	PossibleTypes func() []__Type `json:"possibleTypes"`
	// must be non-null for ENUM, otherwise null.
	EnumValues func(isDeprecatedArgs) []__EnumValue `json:"-"`
	// must be non-null for INPUT_OBJECT, otherwise null.
	InputFields func(isDeprecatedArgs) []__InputValue `json:"-"`
	// must be non-null for NON_NULL and LIST, otherwise null.
	OfType *__Type `json:"ofType"`
	// may be non-null for custom SCALAR, otherwise null.
	SpecifiedByURL *string `json:"specifiedByUrl"`

	JSONKind        string    `json:"kind"`
	JSONFields      []__Field `json:"fields"`
	JSONEnumValues  []__Field `json:"enumValues"`
	JSONInputFields []__Field `json:"inputFields"`
}

type __Field struct {
	Name              string                                `json:"name"`
	Description       *string                               `json:"description"`
	Args              func(isDeprecatedArgs) []__InputValue `json:"-"`
	Type              __Type                                `json:"type"`
	IsDeprecated      bool                                  `json:"isDeprecated"`
	DeprecationReason *string                               `json:"deprecationReason"`

	JSONArgs []__Field `json:"args"`
}

type __InputValue struct {
	Name              string  `json:"name"`
	Description       *string `json:"description"`
	Type              __Type  `json:"type"`
	DefaultValue      *string `json:"defaultValue"`
	IsDeprecated      bool    `json:"isDeprecated"`
	DeprecationReason *string `json:"deprecationReason"`
}

type __EnumValue struct {
	Name              string  `json:"name"`
	Description       *string `json:"description"`
	IsDeprecated      bool    `json:"isDeprecated"`
	DeprecationReason *string `json:"deprecationReason"`
}

type __Directive struct {
	Name         string                                `json:"name"`
	Description  *string                               `json:"description"`
	Locations    []__DirectiveLocation                 `json:"-"`
	Args         func(isDeprecatedArgs) []__InputValue `json:"-"`
	IsRepeatable bool                                  `json:"isRepeatable"`

	JSONLocations []string  `json:"locations"`
	JSONArgs      []__Field `json:"args"`
}

type isDeprecatedArgs struct {
	IncludeDeprecated bool `json:"includeDeprecated"`
}

type __TypeKind uint8

const (
	typeKindScalar __TypeKind = iota
	typeKindObject
	typeKindInterface
	typeKindUnion
	typeKindEnum
	typeKindInputObject
	typeKindList
	typeKindNonNull
)

var typeKindEnumMap = NewBiMap[string, __TypeKind](WithInitialMap(map[string]__TypeKind{
	"SCALAR":       typeKindScalar,
	"OBJECT":       typeKindObject,
	"INTERFACE":    typeKindInterface,
	"UNION":        typeKindUnion,
	"ENUM":         typeKindEnum,
	"INPUT_OBJECT": typeKindInputObject,
	"LIST":         typeKindList,
	"NON_NULL":     typeKindNonNull,
}))

func (my __TypeKind) String() string {
	if v, ok := typeKindEnumMap.GetInverse(my); ok {
		return v
	}
	return ""
}

type __DirectiveLocation uint8

const (
	directiveLocationQuery __DirectiveLocation = iota
	directiveLocationMutation
	directiveLocationSubscription
	directiveLocationField
	directiveLocationFragmentDefinition
	directiveLocationFragmentSpread
	directiveLocationInlineFragment
	directiveLocationVariableDefinition
	directiveLocationSchema
	directiveLocationScalar
	directiveLocationObject
	directiveLocationFieldDefinition
	directiveLocationArgumentDefinition
	directiveLocationInterface
	directiveLocationUnion
	directiveLocationEnum
	directiveLocationEnumValue
	directiveLocationInputObject
	directiveLocationInputFieldDefinition
)

var directiveLocationMap = NewBiMap[string, __DirectiveLocation](WithInitialMap(map[string]__DirectiveLocation{
	"QUERY":                  directiveLocationQuery,
	"MUTATION":               directiveLocationMutation,
	"SUBSCRIPTION":           directiveLocationSubscription,
	"FIELD":                  directiveLocationField,
	"FRAGMENT_DEFINITION":    directiveLocationFragmentDefinition,
	"FRAGMENT_SPREAD":        directiveLocationFragmentSpread,
	"INLINE_FRAGMENT":        directiveLocationInlineFragment,
	"VARIABLE_DEFINITION":    directiveLocationVariableDefinition,
	"SCHEMA":                 directiveLocationSchema,
	"SCALAR":                 directiveLocationScalar,
	"OBJECT":                 directiveLocationObject,
	"FIELD_DEFINITION":       directiveLocationFieldDefinition,
	"ARGUMENT_DEFINITION":    directiveLocationArgumentDefinition,
	"INTERFACE":              directiveLocationInterface,
	"UNION":                  directiveLocationUnion,
	"ENUM":                   directiveLocationEnum,
	"ENUM_VALUE":             directiveLocationEnumValue,
	"INPUT_OBJECT":           directiveLocationInputObject,
	"INPUT_FIELD_DEFINITION": directiveLocationInputFieldDefinition,
}))

func (my __DirectiveLocation) String() string {
	if v, ok := directiveLocationMap.GetInverse(my); ok {
		return v
	}
	return ""
}
