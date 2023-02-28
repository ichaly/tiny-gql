package test

import (
	"database/sql"
	"github.com/bytedance/sonic"
	"github.com/ichaly/tiny-go/core"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	dialect string
	db      *sql.DB
)

func init() {
	var err error
	dialect = "postgres"
	url := "postgres://postgres:postgres@localhost:5432/blog_development?sslmode=disable"
	db, err = sql.Open("pgx", url)
	if err != nil {
		panic(err)
	}
}

func main() {
	info, err := core.GetDBInfo(db, dialect, nil)
	if err != nil {
		panic(err)
	}
	str, err := sonic.MarshalString(info)
	if err != nil {
		panic(err)
	}
	println(str)
}
