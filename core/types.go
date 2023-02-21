package core

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
		Kind:       TK_OBJECT,
		Name:       "Query",
		Interfaces: []__Type{},
	}, {
		Kind:       TK_OBJECT,
		Name:       "Subscription",
		Interfaces: []__Type{},
	}, {
		Kind:       TK_OBJECT,
		Name:       "Mutation",
		Interfaces: []__Type{},
	}, {
		Kind: TK_ENUM,
		Name: "SearchInput",
		EnumValues: func(args isDeprecatedArgs) []__Field {
			return []__Field{
				{
					Name:        "children",
					Description: "Children of parent row",
				}, {
					Name:        "parents",
					Description: "Parents of current row",
				},
			}
		},
	}, {
		Kind:        TK_ENUM,
		Name:        "OrderDirection",
		Description: "Result ordering types",
		EnumValues: func(args isDeprecatedArgs) []__Field {
			return []__Field{
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
			}
		},
	},
}
