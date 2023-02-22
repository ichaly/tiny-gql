package test

import (
	"database/sql"
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
	conf, err := core.ReadInConfig(filepath.Join("./cfg", "prod.yml"))
	if err != nil {
		panic(err)
		return
	}
	info, err := core.GetDBInfo(db, dialect, conf.Blocklist)
	if err != nil {
		panic(err)
		return
	}
	in, err := core.NewSchema(conf, info)
	if err == nil {
		println(string(in))
	} else {
		panic(err)
	}
}
