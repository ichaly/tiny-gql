package core

import (
	"encoding/json"
	"fmt"
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

func (my *__Schema) addTablesEnum() {
	te := __Type{
		Kind:        TK_ENUM,
		Name:        "Tables" + SUFFIX_ENUM,
		Description: "All available tables",
	}

	for _, t := range my.info.Tables {
		if t.Blocked {
			continue
		}
		f := __EnumValue{Name: my.getName(t.Name)}
		if t.Comment != nil {
			f.Description = *t.Comment
		}
		te.EnumValues = append(te.EnumValues, f)
	}

	my.addType(te)
}

func (my *__Schema) addColumnsEnumType(t *DBTable) (err error) {
	tableName := my.getName(t.Name)
	ft := __Type{
		Kind:        TK_ENUM,
		Name:        tableName + "Columns" + SUFFIX_ENUM,
		Description: fmt.Sprintf("Table columns for '%s'", tableName),
	}
	for _, c := range t.Columns {
		if c.Blocked {
			continue
		}
		f := __EnumValue{Name: my.getName(c.Name)}
		if t.Comment != nil {
			f.Description = *c.Comment
		}
		ft.EnumValues = append(ft.EnumValues, f)
	}
	my.addType(ft)
	return
}

func (my *__Schema) addTypeTo(op string, ft __Type) {
	qt := my.Types[op]
	qt.Fields = append(qt.Fields, __Field{
		Name:        ft.Name,
		Description: ft.Description,
		Type:        &__Type{Name: ft.Name},
	})
	my.Types[op] = qt
}

func (my *__Schema) addTable(t *DBTable, alias string) (err error) {
	if t.Blocked || len(t.Columns) == 0 {
		return
	}

	// add table type to query and subscription
	tt, err := my.getTableType(t, alias, 0)
	if err != nil {
		return
	}
	my.addType(tt)
	my.addTypeTo("Query", tt)
	my.addTypeTo("Subscription", tt)

	return
}

func getTypeFromColumn(col DBColumn) (gqlType string) {
	if col.PrimaryKey {
		gqlType = "ID"
		return
	}
	gqlType, _ = getType(col.Type)
	return
}

func (my *__Schema) getColumnField(c DBColumn) (f __Field, err error) {
	t := __Type{Name: "String"}

	if v, ok := my.Types[getTypeFromColumn(c)]; ok {
		t.Name = v.Name
		t.Kind = v.Kind
	}

	if c.Array {
		t = __Type{Kind: TK_LIST, OfType: &__Type{
			Name: t.Name, Kind: t.Kind,
		}}
	}

	if c.NotNull {
		t = __Type{Kind: TK_NON_NULL, OfType: &__Type{
			Name: t.Name, Kind: t.Kind,
		}}
	}

	f.Name = my.getName(c.Name)
	f.Type = &t

	//f.Args = append(f.Args, __InputValue{
	//	Name: "includeIf",
	//	Type: &__Type{Name: c.Table + SUFFIX_WHERE},
	//})
	//
	//f.Args = append(f.Args, __InputValue{
	//	Name: "skipIf",
	//	Type: &__Type{Name: c.Table + SUFFIX_WHERE},
	//})
	return
}

func (my *__Schema) getTableType(t *DBTable, alias string, depth int) (ft __Type, err error) {
	ft = __Type{
		Kind: TK_OBJECT,
		//Interfaces: []__Type{},
	}

	name := t.Name
	if alias != "" {
		name = alias
	}
	ft.Name = my.getName(name)

	if t.Comment != nil {
		ft.Description = *t.Comment
	}

	if err = my.addColumnsEnumType(t); err != nil {
		return
	}

	for _, c := range t.Columns {
		if c.Blocked {
			continue
		}
		if c.FullText {
			//hasSearch = true
		}
		if c.FKRecursive {
			//hasRecursive = true
		}
		var f __Field
		f, err = my.getColumnField(c)
		if err != nil {
			return
		}
		ft.Fields = append(ft.Fields, f)
	}
	return
}

func NewSchema(conf *Config, info *DBInfo) (res json.RawMessage, err error) {
	s := __Schema{
		conf:       conf,
		info:       info,
		Types:      map[string]__Type{},
		Directives: map[string]__Directive{},
		QueryType:  __Type{Name: "Query"},
		//MutationType:     &__Type{Name: "Mutation"},
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

	s.addTablesEnum()

	for _, t := range s.info.Tables {
		if err = s.addTable(t, ""); err != nil {
			return
		}
	}

	root := map[string]interface{}{"data": map[string]interface{}{"__schema": s}}
	return sonic.Marshal(root)
}
