package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ichaly/tiny-go/core"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
	"os"
	"path/filepath"
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
	r := chi.NewRouter()
	r.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		in, err := core.NewSchema(conf, info)
		if err != nil {
			panic(err)
		}
		_, _ = w.Write(in)
	})
	r.HandleFunc("/intro", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile("../conf/intro.json")
		if err != nil {
			return
		}
		_, _ = w.Write(file)
	})
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
