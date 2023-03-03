package core

import "strings"

var dbTypes = map[string]string{
	"timestamp without time zone": "Time",
	"character varying":           "String",
	"text":                        "String",
	"smallint":                    "Int",
	"integer":                     "Int",
	"bigint":                      "Int",
	"smallserial":                 "Int",
	"serial":                      "Int",
	"bigserial":                   "Int",
	"decimal":                     "Float",
	"numeric":                     "Float",
	"real":                        "Float",
	"double precision":            "Float",
	"money":                       "Float",
	"boolean":                     "Boolean",
}

const (
	SUFFIX_EXP     = "Expression"
	SUFFIX_LISTEXP = "ListExpression"

	SUFFIX_ENUM   = "Enum"
	SUFFIX_ARGS   = "ArgsInput"
	SUFFIX_SORT   = "SortInput"
	SUFFIX_WHERE  = "WhereInput"
	SUFFIX_UPSERT = "UpsertInput"
	SUFFIX_INSERT = "InsertInput"
	SUFFIX_UPDATE = "UpdateInput"
)

var stdTypes = []__Type{
	{
		Kind:        TK_SCALAR,
		Name:        "ID",
		Description: "The `ID` scalar type represents a unique identifier, often used to refetch an object or as key for a cache.\nThe ID type appears in a JSON response as a String; however, it is not intended to be human-readable.\nWhen expected as an input type, any string (such as `\"4\"`) or integer (such as `4`) input value will be accepted\nas an ID.\n",
	}, {
		Kind:        TK_SCALAR,
		Name:        "Boolean",
		Description: "The `Boolean` scalar type represents `true` or `false`",
	}, {
		Kind:        TK_SCALAR,
		Name:        "Float",
		Description: "The `Float` scalar type represents signed double-precision fractional values as specified by\n[IEEE 754](http://en.wikipedia.org/wiki/IEEE_floating_point).",
	}, {
		Kind:        TK_SCALAR,
		Name:        "Int",
		Description: "The `Int` scalar type represents non-fractional signed whole numeric values. Int can represent\nvalues between -(2^31) and 2^31 - 1.\n",
	}, {
		Kind:        TK_SCALAR,
		Name:        "String",
		Description: "The `String` scalar type represents textual data, represented as UTF-8 character sequences.\nThe String type is most often used by GraphQL to represent free-form human-readable text.\n",
	}, {
		Kind:        TK_SCALAR,
		Name:        "JSON",
		Description: "The `JSON` scalar type represents json data",
	}, {
		Kind:        TK_SCALAR,
		Name:        "Cursor",
		Description: "A cursor is an encoded string use for pagination",
	}, {
		Kind:           TK_SCALAR,
		Name:           "File",
		Description:    "The `File` scalar type references to a multipart file, often used to upload files to the server. Expects a string with the form file field name",
		SpecifiedByURL: "https://github.com/mjarkk/yarql#file-upload",
	}, {
		Kind:           TK_SCALAR,
		Name:           "Time",
		Description:    "The `Time` scalar type references to a ISO 8601 date+time, often used to insert and/or view dates. Expects a string with the ISO 8601 format",
		SpecifiedByURL: "https://en.wikipedia.org/wiki/ISO_8601",
	}, {
		Kind: TK_OBJECT,
		Name: "Query",
	}, {
		Kind: TK_OBJECT,
		Name: "Mutation",
	}, {
		Kind: TK_OBJECT,
		Name: "Subscription",
	}, {
		Kind:        TK_ENUM,
		Name:        "Recursive",
		Description: "Recursive relation types",
		EnumValues: []__EnumValue{
			{
				Name:        "children",
				Description: "Children of parent row",
			}, {
				Name:        "parents",
				Description: "Parents of current row",
			},
		},
	}, {
		Kind:        TK_ENUM,
		Name:        "Direction",
		Description: "Result ordering types",
		EnumValues: []__EnumValue{
			{
				Name:        "asc",
				Description: "Ascending order",
			}, {
				Name:        "desc",
				Description: "Descending order",
			}, {
				Name:        "asc_nulls_first",
				Description: "Ascending nulls first order",
			}, {
				Name:        "desc_nulls_first",
				Description: "Descending nulls first order",
			}, {
				Name:        "asc_nulls_last",
				Description: "Ascending nulls last order",
			}, {
				Name:        "desc_nulls_last",
				Description: "Descending nulls last order",
			},
		},
	},
}

var expAll = []__InputValue{
	{
		Name:        "isNull",
		Type:        &__Type{Kind: TK_SCALAR, Name: "Boolean"},
		Description: "Is value null (true) or not null (false)",
	},
}

var expScalar = []__InputValue{
	{Name: "equals", Description: "Equals value"},
	{Name: "notEquals", Description: "Does not equal value"},
	{Name: "greaterThan", Description: "Is greater than value"},
	{Name: "lesserThan", Description: "Is lesser than value"},
	{Name: "greaterOrEquals", Description: "Is greater than or equals value"},
	{Name: "lesserOrEquals", Description: "Is lesser than or equals value"},
	{Name: "like", Description: "Value matching (case-insensitive) pattern where '%' represents zero or more characters and '_' represents a single character. Eg. '_r%' finds values having 'r' in second position"},
	{Name: "notLike", Description: "Value not matching pattern where '%' represents zero or more characters and '_' represents a single character. Eg. '_r%' finds values not having 'r' in second position"},
	{Name: "iLike", Description: "Value matching (case-insensitive) pattern where '%' represents zero or more characters and '_' represents a single character. Eg. '_r%' finds values having 'r' in second position"},
	{Name: "notILike", Description: "Value not matching (case-insensitive) pattern where '%' represents zero or more characters and '_' represents a single character. Eg. '_r%' finds values not having 'r' in second position"},
	{Name: "similar", Description: "Value matching regex pattern. Similar to the 'like' operator but with support for regex. Pattern must match entire value."},
	{Name: "notSimilar", Description: "Value not matching regex pattern. Similar to the 'like' operator but with support for regex. Pattern must not match entire value."},
	{Name: "regex", Description: "Value matches regex pattern"},
	{Name: "notRegex", Description: "Value not matching regex pattern"},
	{Name: "iRegex", Description: "Value matches (case-insensitive) regex pattern"},
	{Name: "notIRegex", Description: "Value not matching (case-insensitive) regex pattern"},
}

var expList = []__InputValue{
	{Name: "in", Description: "Is in list of values"},
	{Name: "notIn", Description: "Is not in list of values"},
}

var expJSON = []__InputValue{
	{Name: "hasKey", Description: "JSON value contains this key"},
	{Name: "hasKeyAny", Description: "JSON value contains any of these keys"},
	{Name: "hasKeyAll", Description: "JSON value contains all of these keys"},
	{Name: "contains", Description: "JSON value matches any of they key/value pairs"},
	{Name: "containedIn", Description: "JSON value contains all of they key/value pairs"},
}

var argsList = []__InputValue{
	{Name: "id", Type: &__Type{Name: "ID"}},
	{Name: "limit", Type: &__Type{Name: "Int"}},
	{Name: "offset", Type: &__Type{Name: "Int"}},
	{Name: "first", Type: &__Type{Name: "Int"}},
	{Name: "last", Type: &__Type{Name: "Int"}},
	{Name: "before", Type: &__Type{Name: "Cursor"}},
	{Name: "after", Type: &__Type{Name: "Cursor"}},
	{Name: "distinctOn", Type: &__Type{Kind: TK_LIST, OfType: &__Type{Name: "String"}}},
}

func getType(t string) (gqlType string, list bool) {
	if i := strings.IndexRune(t, '('); i != -1 {
		t = t[:i]
	}
	if i := strings.IndexRune(t, '['); i != -1 {
		list = true
		t = t[:i]
	}
	if v, ok := dbTypes[t]; ok {
		gqlType = v
	} else if t == "json" || t == "jsonb" {
		gqlType = "JSON"
	} else {
		gqlType = "String"
	}
	return
}
