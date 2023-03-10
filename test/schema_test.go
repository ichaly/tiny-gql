package test

import (
	"database/sql"
	"fmt"
	"github.com/ichaly/tiny-go/core"
	"path/filepath"
	"testing"
)

func init() {
	var err error
	dialect = "postgres"
	url := "postgres://postgres:postgres@localhost:5432/test_development?sslmode=disable"
	db, err = sql.Open("pgx", url)
	if err != nil {
		panic(err)
	}
}

func TestNewSchema(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	conf, err := core.ReadInConfig(filepath.Join("../conf", "prod.yml"))
	if err != nil {
		panic(err)
	}
	info, err := core.GetDBInfo(db, dialect, conf.Blocklist)
	if err != nil {
		panic(err)
	}
	in, err := core.NewSchema(conf, info)
	if err != nil {
		panic(err)
	}
	println(string(in))
}
