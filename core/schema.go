package core

import (
	"encoding/json"
	"github.com/bytedance/sonic"
	"github.com/iancoleman/strcase"
)

//
// Types represent:
// https://spec.graphql.org/draft/#sec-Schema-Introspection
//

type mapArray[V any] map[string]V

func (my mapArray[V]) MarshalJSON() ([]byte, error) {
	var list []V
	for _, t := range my {
		list = append(list, t)
	}
	return sonic.Marshal(list)
}

type __Schema struct {
	Types            mapArray[__Type]      `json:"types"`
	QueryType        __Type                `json:"queryType"`
	Directives       mapArray[__Directive] `json:"directives,omitempty"`
	Description      string                `json:"description,omitempty"`
	MutationType     *__Type               `json:"mutationType,omitempty"`
	SubscriptionType *__Type               `json:"subscriptionType,omitempty"`

	conf *Config
	info *DBInfo
}

type __Type struct {
	Kind        __TypeKind `json:"kind"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	// must be non-null for OBJECT and INTERFACE, otherwise null.
	Fields []__Field `json:"fields"`
	// must be non-null for OBJECT and INTERFACE, otherwise null.
	Interfaces []__Type `json:"interfaces"`
	// must be non-null for INTERFACE and UNION, otherwise null.
	PossibleTypes []__Type `json:"possibleTypes"`
	// must be non-null for ENUM, otherwise null.
	EnumValues []__EnumValue `json:"enumValues"`
	// must be non-null for INPUT_OBJECT, otherwise null.
	InputFields []__InputValue `json:"inputFields"`
	// must be non-null for NON_NULL and LIST, otherwise null.
	OfType *__Type `json:"ofType"`
	// may be non-null for custom SCALAR, otherwise null.
	SpecifiedByURL string `json:"specifiedByUrl"`
}

func (my __Type) MarshalJSON() ([]byte, error) {
	obj := make(map[string]interface{})
	if my.Kind != "" {
		obj["kind"] = my.Kind
	}
	if my.Name != "" {
		obj["name"] = my.Name
	}
	if my.Description != "" {
		obj["description"] = my.Description
	}
	if my.Kind == TK_OBJECT || my.Kind == TK_INTERFACE {
		obj["fields"] = append([]__Field{}, my.Fields...)
		obj["interfaces"] = append([]__Type{}, my.Interfaces...)
	}
	if my.Kind == TK_INTERFACE || my.Kind == TK_UNION {
		obj["possibleTypes"] = append([]__Type{}, my.PossibleTypes...)
	}
	if my.Kind == TK_ENUM {
		obj["enumValues"] = append([]__EnumValue{}, my.EnumValues...)
	}
	if my.Kind == TK_INPUT_OBJECT {
		obj["inputFields"] = append([]__InputValue{}, my.InputFields...)
	}
	if my.OfType != nil {
		obj["ofType"] = my.OfType
	}
	if my.SpecifiedByURL != "" {
		obj["specifiedByUrl"] = my.SpecifiedByURL
	}
	return sonic.Marshal(obj)
}

type __Field struct {
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	Args              []__InputValue `json:"args"`
	Type              *__Type        `json:"type"`
	IsDeprecated      bool           `json:"isDeprecated"`
	DeprecationReason string         `json:"deprecationReason"`
}

func (my __Field) MarshalJSON() ([]byte, error) {
	obj := make(map[string]interface{})
	if my.Name != "" {
		obj["name"] = my.Name
	}
	if my.Description != "" {
		obj["description"] = my.Description
	}
	obj["args"] = append([]__InputValue{}, my.Args...)
	if my.Type != nil {
		obj["type"] = my.Type
	}
	obj["isDeprecated"] = my.IsDeprecated
	if my.DeprecationReason != "" {
		obj["deprecationReason"] = my.DeprecationReason
	}
	return sonic.Marshal(obj)
}

type __InputValue struct {
	Name              string  `json:"name,omitempty"`
	Description       string  `json:"description,omitempty"`
	Type              *__Type `json:"type,omitempty"`
	IsDeprecated      bool    `json:"isDeprecated"`
	DeprecationReason string  `json:"deprecationReason,omitempty"`
	DefaultValue      string  `json:"defaultValue,omitempty"`
}

type __EnumValue struct {
	Name              string `json:"name,omitempty"`
	Description       string `json:"description,omitempty"`
	IsDeprecated      bool   `json:"isDeprecated"`
	DeprecationReason string `json:"deprecationReason,omitempty"`
}

type __Directive struct {
	Name         string                `json:"name,omitempty"`
	Description  string                `json:"description,omitempty"`
	Locations    []__DirectiveLocation `json:"locations"`
	Args         []__InputValue        `json:"args"`
	IsRepeatable bool                  `json:"isRepeatable"`
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

func (my *__Schema) getName(name string) string {
	if my.conf.EnableCamelcase {
		return strcase.ToCamel(name)
	} else {
		return name
	}
}

func (my *__Schema) addType(t __Type) {
	my.Types[t.Name] = t
}

func (my *__Schema) addTypeTo(op string, ft __Type) {
	qt := my.Types[op]
	qt.Fields = append(qt.Fields, __Field{
		Name:        strcase.ToLowerCamel(ft.Name),
		Description: ft.Description,
		Type:        &__Type{Name: ft.Name},
		Args:        ft.InputFields,
	})
	my.Types[op] = qt
}

func (my *__Schema) addExpression(exps []__InputValue, name string, sub __Type) {
	t := __Type{
		Kind:        TK_INPUT_OBJECT,
		Name:        name + SUFFIX_EXP,
		InputFields: []__InputValue{},
	}

	for _, ex := range exps {
		if ex.Type == nil {
			ex.Type = &sub
		}
		t.InputFields = append(t.InputFields, ex)
	}
	my.addType(t)
}

func NewSchema(conf *Config, info *DBInfo) (res json.RawMessage, err error) {
	s := __Schema{
		conf:             conf,
		info:             info,
		Types:            map[string]__Type{},
		Directives:       map[string]__Directive{},
		QueryType:        __Type{Name: "Query"},
		MutationType:     &__Type{Name: "Mutation"},
		SubscriptionType: &__Type{Name: "Subscription"},
	}

	for _, v := range stdTypes {
		s.addType(v)
	}

	// Expression types
	v := append(expAll, expScalar...)
	s.addExpression(v, "ID", __Type{Kind: TK_SCALAR, Name: "ID"})
	s.addExpression(v, "Int", __Type{Kind: TK_SCALAR, Name: "Int"})
	s.addExpression(v, "Float", __Type{Kind: TK_SCALAR, Name: "Float"})
	s.addExpression(v, "String", __Type{Kind: TK_SCALAR, Name: "String"})
	s.addExpression(v, "Boolean", __Type{Kind: TK_SCALAR, Name: "Boolean"})

	// ListExpression Types
	v = append(expAll, expList...)
	s.addExpression(v, "IntList", __Type{Kind: TK_SCALAR, Name: "Int"})
	s.addExpression(v, "FloatList", __Type{Kind: TK_SCALAR, Name: "Float"})
	s.addExpression(v, "StringList", __Type{Kind: TK_SCALAR, Name: "String"})
	s.addExpression(v, "BooleanList", __Type{Kind: TK_SCALAR, Name: "Boolean"})

	// JsonExpression types
	v = append(expAll, expJSON...)
	s.addExpression(v, "JSON", __Type{Kind: TK_SCALAR, Name: "String"})

	s.addTablesType()

	root := map[string]interface{}{"data": map[string]interface{}{"__schema": s}}
	return sonic.Marshal(root)
}

func (my *__Schema) addTablesType() {
	var enumValues []__EnumValue

	for _, t := range my.info.Tables {
		if t.Blocked {
			continue
		}
		// append tables enum value to type
		enumValues = append(enumValues, __EnumValue{Name: my.getName(t.Name), Description: t.Comment})

		// where input type
		wn := strcase.ToCamel(t.Name) + SUFFIX_WHERE
		wi := __Type{
			Kind: TK_INPUT_OBJECT,
			Name: wn,
			InputFields: []__InputValue{
				{Name: "not", Type: &__Type{Name: wn}},
				{Name: "and", Type: &__Type{Name: wn}},
				{Name: "or", Type: &__Type{Name: wn}},
			},
		}

		// order by input type
		oi := __Type{
			Kind: TK_INPUT_OBJECT,
			Name: strcase.ToCamel(t.Name) + SUFFIX_ORDER_BY,
		}

		// table object type
		to := __Type{
			Kind:        TK_OBJECT,
			Name:        my.getName(t.Name),
			Description: t.Comment,
			InputFields: argsList,
		}

		for _, c := range t.Columns {
			if c.Blocked {
				continue
			}

			//get column scalar type
			cn, isList := my.getColumnType(c)

			// append order by input fields
			oi.InputFields = append(oi.InputFields, __InputValue{
				Name:        strcase.ToCamel(c.Name),
				Description: c.Comment,
				Type:        &__Type{Name: "Direction"},
			})

			// append where input fields
			iv := __InputValue{
				Name:        strcase.ToLowerCamel(c.Name),
				Description: c.Comment,
			}
			if c.Array || isList {
				iv.Type = &__Type{Name: cn + SUFFIX_LISTEXP}
			} else {
				iv.Type = &__Type{Name: cn + SUFFIX_EXP}
			}
			wi.InputFields = append(wi.InputFields, iv)

			// append table object field
			ct := __Type{Name: cn}
			if c.Array || isList {
				ct = __Type{Kind: TK_LIST, OfType: &__Type{
					Name: cn,
				}}
			}
			if c.NotNull {
				ct = __Type{Kind: TK_NON_NULL, OfType: &__Type{
					Name: cn,
				}}
			}
			to.Fields = append(to.Fields, __Field{
				Name:        strcase.ToLowerCamel(c.Name),
				Description: c.Comment,
				Type:        &ct,
				Args: []__InputValue{
					{Name: "includeIf", Type: &__Type{Name: wi.Name}},
					{Name: "skipIf", Type: &__Type{Name: wi.Name}},
				},
			})
		}

		// add order by input to types
		my.addType(oi)

		// add where input to types
		my.addType(wi)

		// add table object to types
		my.addType(to)

		// add to Query and Subscription
		to.InputFields = append(to.InputFields, __InputValue{
			Name: "orderBy", Type: &__Type{Name: oi.Name},
		})
		to.InputFields = append(to.InputFields, __InputValue{
			Name: "where", Type: &__Type{Name: wi.Name},
		})
		my.addTypeTo("Query", to)
		my.addTypeTo("Subscription", to)

		// add to Mutation
		to.InputFields = append(to.InputFields, __InputValue{Name: "delete", Type: &__Type{Name: "Boolean"}})
		my.addTypeTo("Mutation", to)
	}

	// add tables enum to types
	my.addType(__Type{
		Kind:        TK_ENUM,
		Name:        "Tables" + SUFFIX_ENUM,
		Description: "All available tables",
		EnumValues:  enumValues,
	})
}

func (my *__Schema) getColumnType(c DBColumn) (name string, isList bool) {
	if c.PrimaryKey {
		name = "ID"
		return
	}
	name, isList = getType(c.Type)
	return
}
