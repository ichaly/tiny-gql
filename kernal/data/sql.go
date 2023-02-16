package data

import _ "embed"

//go:embed sql/postgres_info.sql
var postgresInfo string

//go:embed sql/postgres_columns.sql
var postgresColumns string

//go:embed sql/postgres_functions.sql
var postgresFunctions string
