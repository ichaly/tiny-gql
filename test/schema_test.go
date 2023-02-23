package test

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ichaly/tiny-go/core"
	"net/http"
	"os"
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
	conf, err := core.ReadInConfig(filepath.Join("./cfg", "prod.yml"))
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

	r := chi.NewRouter()
	//r.Post("/graphql", func(w http.ResponseWriter, r *http.Request) {
	//	_, _ = w.Write(in)
	//})
	r.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(in)
	})
	r.HandleFunc("/intro", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile("./intro.json")
		if err != nil {
			return
		}
		_, _ = w.Write(file)
	})
	_ = http.ListenAndServe(":3000", r)
}
