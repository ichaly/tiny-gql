package test

import (
	"database/sql"
	"github.com/ichaly/tiny-go/core"
	jsoniter "github.com/json-iterator/go"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	dialect string
	db      *sql.DB
	json    = jsoniter.ConfigCompatibleWithStandardLibrary
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

func TestGetDBInfo(t *testing.T) {
	info, err := core.GetDBInfo(db, dialect, nil)
	if err != nil {
		panic(err)
	}
	str, err := json.MarshalToString(info)
	if err != nil {
		panic(err)
	}
	println(str)
}
