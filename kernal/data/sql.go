package data

import _ "embed"

//go:embed sql/postgres_info.sql
var postgresInfo string

//go:embed sql/postgres_columns.sql
var postgresColumns string

//go:embed sql/postgres_functions.sql
var postgresFunctions string

//go:embed sql/mysql_info.sql
var mysqlInfo string

//go:embed sql/mysql_columns.sql
var mysqlColumns string

//go:embed sql/mysql_functions.sql
var mysqlFunctions string
