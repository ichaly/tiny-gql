package core

import (
	"database/sql"
	"fmt"
	"github.com/ichaly/tiny-go/core/internal"
	"github.com/ichaly/tiny-go/core/internal/data"
	"regexp"
	"strings"
)

type DBInfo struct {
	Dialect string
	Version int
	Schema  string
	Name    string
	Tables  map[string]*DBTable
	VTables []VirtualTable `json:"-"` // for polymorphic relationships

	relation data.BiDict
	hash     int
}

func (my *DBInfo) Hash() int {
	return my.hash
}

type VirtualTable struct {
	Name       string
	IDColumn   string
	TypeColumn string
	FKeyColumn string
}

type DBTable struct {
	Name       string
	Schema     string
	Comment    string
	Type       string // json,jsonb,polymorphic,virtual
	PrimaryCol DBColumn
	Columns    map[string]DBColumn
	FullText   map[string]DBColumn
	Blocked    bool
}

func (my *DBTable) String() string {
	return my.Schema + "." + my.Name
}

type DBColumn struct {
	Comment     string
	Name        string
	Type        string
	Array       bool
	NotNull     bool
	PrimaryKey  bool
	UniqueKey   bool
	FullText    bool
	FKRecursive bool
	FKeySchema  string
	FKeyTable   string
	FKeyCol     string
	Blocked     bool

	Schema      string
	Table       string
	Description string // table comment
}

func (my DBColumn) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s.%s.%s", my.Schema, my.Table, my.Name))
	sb.WriteString(fmt.Sprintf(
		" [type:%v, array:%t, notNull:%t, fullText:%t]",
		my.Type, my.Array, my.NotNull, my.FullText,
	))

	if my.FKeyCol != "" {
		sb.WriteString(fmt.Sprintf(" -> %s.%s.%s", my.FKeySchema, my.FKeyTable, my.FKeyCol))
	}
	return sb.String()
}

func GetDBInfo(db *sql.DB, dialect string, blockList []string) (*DBInfo, error) {
	var err error
	var dbVersion int
	var dbSchema, dbName string

	// get db info
	row := db.QueryRow(internal.PostgresInfo)
	if row.Scan(&dbVersion, &dbSchema, &dbName) != nil {
		return nil, err
	}

	// get columns from db
	rows, err := db.Query(internal.PostgresColumns)
	if err != nil {
		return nil, fmt.Errorf("error fetching columns: %s", err)
	}
	defer rows.Close()

	di := &DBInfo{
		Dialect:  dialect,
		Version:  dbVersion,
		Schema:   dbSchema,
		Name:     dbName,
		Tables:   make(map[string]*DBTable),
		relation: make(map[string][]string),
	}

	// we have to rescan and update columns to overcome
	// weird bugs in mysql like joins with information_schema
	// don't work in 8.0.22 etc.
	for rows.Next() {
		var c DBColumn
		err = rows.Scan(
			&c.Schema, &c.Table, &c.Description, &c.Name, &c.Comment, &c.Type, &c.NotNull,
			&c.PrimaryKey, &c.UniqueKey, &c.Array, &c.FullText, &c.FKeySchema, &c.FKeyTable, &c.FKeyCol,
		)

		// safety verification
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(c.Table, "_gj_") {
			continue
		}

		// data perfection
		c.Blocked = isBlocked(c.Name, blockList)
		if c.PrimaryKey {
			c.UniqueKey = true
		}
		if c.FKeySchema == c.Schema && c.FKeyTable == c.Table {
			c.FKRecursive = true
		}

		// fill data
		tk := fmt.Sprintf("%s:%s", c.Schema, c.Table)
		ck := fmt.Sprintf("%s:%s:%s", c.Schema, c.Table, c.Name)
		t, ok := di.Tables[tk]
		if !ok {
			t = &DBTable{
				Name:     c.Table,
				Schema:   c.Schema,
				Comment:  c.Description,
				Columns:  make(map[string]DBColumn),
				FullText: make(map[string]DBColumn),
			}
			if isBlocked(c.Table, blockList) {
				t.Blocked = true
			}
			di.Tables[tk] = t
		}
		if c.FullText {
			t.FullText[ck] = c
		}
		if c.PrimaryKey {
			t.PrimaryCol = c
		}
		if c.FKeyTable != "" {
			di.relation.Put(t.Name, c.FKeyTable)
		}
		t.Columns[ck] = c
	}
	return di, nil
}

func isBlocked(val string, list []string) bool {
	for _, v := range list {
		regex := fmt.Sprintf("^%s$", v)
		if matched, _ := regexp.MatchString(regex, val); matched {
			return true
		}
	}
	return false
}
