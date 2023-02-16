package data

import (
	"database/sql"
	"fmt"
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

	hash int
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
	Comment    *string
	Type       string // json,jsonb,polymorphic,virtual
	PrimaryCol DBColumn
	Columns    map[string]DBColumn
	FullText   map[string]DBColumn
	Blocked    bool
}

type DBColumn struct {
	Comment     *string
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
	Description *string
}

func GetDBInfo(db *sql.DB, dialect string, blockList []string) (*DBInfo, error) {
	var err error
	var dbVersion int
	var dbSchema, dbName string

	// get db info
	row := db.QueryRow(postgresInfo)
	if row.Scan(&dbVersion, &dbSchema, &dbName) != nil {
		return nil, err
	}

	// get columns from db
	rows, err := db.Query(postgresColumns)
	if err != nil {
		return nil, fmt.Errorf("error fetching columns: %s", err)
	}
	defer rows.Close()

	di := &DBInfo{
		Dialect: dialect,
		Version: dbVersion,
		Schema:  dbSchema,
		Name:    dbName,
		Tables:  make(map[string]*DBTable),
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
		c.Blocked = isBlock(c.Name, blockList)
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
				//Type         string
			}
			if isBlock(c.Table, blockList) {
				t.Blocked = true
			}
			di.Tables[tk] = t
		}
		if c.FullText {
			t.FullText[ck] = c
		}
		if c.PrimaryKey {
			t.PrimaryCol = c
		} else {
			t.Columns[ck] = c
		}
	}
	return di, nil
}

func isBlock(val string, list []string) bool {
	for _, v := range list {
		regex := fmt.Sprintf("^%list$", v)
		if matched, _ := regexp.MatchString(regex, val); matched {
			return true
		}
	}
	return false
}
