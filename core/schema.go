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

func (my *__Schema) getName(name string, isField ...bool) string {
	if my.conf.EnableCamelcase {
		if len(isField) > 0 && isField[0] {
			return strcase.ToLowerCamel(name)
		}
		return strcase.ToCamel(name)
	} else {
		return name
	}
}

func (my *__Schema) addType(types ...__Type) {
	for _, t := range types {
		my.Types[t.Name] = t
	}
}

func (my *__Schema) addTypeTo(op string, ot __Type, args []__InputValue) {
	t := my.Types[op]
	t.Fields = append(t.Fields, __Field{
		Name:        my.getName(ot.Name, true),
		Description: ot.Description,
		Type:        &__Type{Name: ot.Name},
		Args:        args,
	})
	my.Types[op] = t
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
	s.addExpression(v, ID, __Type{Name: ID})
	s.addExpression(v, Int, __Type{Name: Int})
	s.addExpression(v, Time, __Type{Name: Time})
	s.addExpression(v, Float, __Type{Name: Float})
	s.addExpression(v, String, __Type{Name: String})
	s.addExpression(v, Boolean, __Type{Name: Boolean})

	// ListExpression Types
	v = append(expAll, expList...)
	s.addExpression(v, "IntList", __Type{Name: Int})
	s.addExpression(v, "TimeList", __Type{Name: Time})
	s.addExpression(v, "FloatList", __Type{Name: Float})
	s.addExpression(v, "StringList", __Type{Name: String})
	s.addExpression(v, "BooleanList", __Type{Name: Boolean})

	// JsonExpression types
	v = append(expAll, expJSON...)
	s.addExpression(v, JSON, __Type{Name: String})

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
		tableName := my.getName(t.Name)
		// append tables enum value object type
		enumValues = append(enumValues, __EnumValue{Name: tableName, Description: t.Comment})

		// sort input type
		sort := __Type{
			Kind: TK_INPUT_OBJECT,
			Name: tableName + SUFFIX_SORT,
		}

		// where input type
		wn := tableName + SUFFIX_WHERE
		where := __Type{
			Kind: TK_INPUT_OBJECT,
			Name: wn,
			InputFields: []__InputValue{
				{Name: "not", Type: &__Type{Name: wn}},
				{Name: "and", Type: &__Type{Name: wn}},
				{Name: "or", Type: &__Type{Name: wn}},
			},
		}

		// upsert input type
		upsert := __Type{
			Kind: TK_INPUT_OBJECT,
			Name: tableName + SUFFIX_UPSERT,
		}
		insert := __Type{
			Kind: TK_INPUT_OBJECT,
			Name: tableName + SUFFIX_INSERT,
		}
		for _, f := range my.info.relation[t.Name] {
			name := my.getName(f) + SUFFIX_UPDATE
			insert.InputFields = append(insert.InputFields, __InputValue{
				Name: f, Type: &__Type{Name: name},
			})
		}
		update := __Type{
			Kind: TK_INPUT_OBJECT,
			Name: tableName + SUFFIX_UPDATE,
			InputFields: []__InputValue{{
				Name: "where", Type: &__Type{Name: where.Name},
				Description: fmt.Sprintf("Update rows in table '%s' that match the expression", tableName),
			}, {
				Name: "connect", Type: &__Type{Name: where.Name},
				Description: fmt.Sprintf("Connect to rows in table '%s' that match the expression", tableName),
			}, {
				Name: "disconnect", Type: &__Type{Name: where.Name},
				Description: fmt.Sprintf("Disconnect from rows in table '%s' that match the expression", tableName),
			}},
		}

		// table object type
		object := __Type{
			Kind:        TK_OBJECT,
			Name:        tableName,
			Description: t.Comment,
		}

		var hasRecursive bool

		for _, c := range t.Columns {
			if c.Blocked {
				continue
			}
			if c.FKRecursive {
				hasRecursive = true
			}

			//get column scalar type
			cn, isList := my.getColumnType(c)
			columnName := my.getName(c.Name, true)

			// append sort by input fields
			sort.InputFields = append(sort.InputFields, __InputValue{
				Name:        columnName,
				Description: c.Comment,
				Type:        &__Type{Name: Direction},
			})

			// append where input fields
			iv := __InputValue{
				Name:        columnName,
				Description: c.Comment,
			}
			if c.Array || isList {
				iv.Type = &__Type{Name: cn + SUFFIX_LISTEXP}
			} else {
				iv.Type = &__Type{Name: cn + SUFFIX_EXP}
			}
			where.InputFields = append(where.InputFields, iv)

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

			upsert.InputFields = append(upsert.InputFields, __InputValue{
				Name: columnName, Type: &ct,
			})
			insert.InputFields = append(insert.InputFields, __InputValue{
				Name: columnName, Type: &ct,
			})
			update.InputFields = append(update.InputFields, __InputValue{
				Name: columnName, Type: &ct,
			})

			object.Fields = append(object.Fields, __Field{
				Name:        columnName,
				Description: c.Comment,
				Type:        &ct,
				Args: []__InputValue{
					{Name: "includeIf", Type: &__Type{Name: where.Name}},
					{Name: "skipIf", Type: &__Type{Name: where.Name}},
				},
			})
		}

		if hasRecursive {
			object.Fields = append(object.Fields, __Field{
				Name: t.Name,
				Type: &__Type{Name: tableName},
				Args: []__InputValue{
					{Name: "includeIf", Type: &__Type{Name: where.Name}},
					{Name: "skipIf", Type: &__Type{Name: where.Name}},
					{Name: "find", Type: &__Type{Name: Recursive}},
				},
			})
		}

		// add to types
		my.addType(sort, where, upsert, insert, update, object)

		// add object Query and Subscription
		args := append(argsList,
			__InputValue{Name: "sort", Type: &__Type{Name: sort.Name}},
			__InputValue{Name: "where", Type: &__Type{Name: where.Name}},
		)
		my.addTypeTo("Query", object, args)
		my.addTypeTo("Subscription", object, args)

		// add object Mutation
		args = append(args,
			__InputValue{Name: "delete", Type: &__Type{Name: Boolean}},
			__InputValue{Name: "upsert", Type: &__Type{Name: upsert.Name}},
			__InputValue{Name: "insert", Type: &__Type{Name: insert.Name}},
			__InputValue{Name: "update", Type: &__Type{Name: update.Name}},
		)
		my.addTypeTo("Mutation", object, args)
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
		name = ID
		return
	}
	name, isList = getType(c.Type)
	return
}
