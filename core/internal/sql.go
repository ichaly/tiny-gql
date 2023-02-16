package internal

import _ "embed"

//go:embed sql/postgres_info.sql
var PostgresInfo string

//go:embed sql/postgres_columns.sql
var PostgresColumns string

//go:embed sql/postgres_functions.sql
var PostgresFunctions string
