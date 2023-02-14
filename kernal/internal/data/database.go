package data

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

type tableKey struct {
	schema string
	table  string
}

type columnKey struct {
	schema string
	table  string
	column string
}

type DBInfo struct {
	Dialect string
	Version string
	Schema  string
	Name    string
	Tables  map[tableKey]DBTable
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
	Comment    string
	Schema     string
	Name       string
	Type       string // json,jsonb,polymorphic,virtual
	Columns    map[columnKey]DBColumn
	PrimaryCol DBColumn
	FullText   map[columnKey]DBColumn
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
	Table       string
	Schema      string
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
	var dbVersion, dbSchema, dbName string

	// get db info
	var sqlSchema string
	switch dialect {
	case "mysql":
		sqlSchema = mysqlInfo
	default:
		sqlSchema = postgresInfo
	}
	row := db.QueryRow(sqlSchema)
	if row.Scan(&dbVersion, &dbSchema, &dbName) != nil {
		return nil, err
	}

	// get columns from db
	var sqlColumn string
	switch dialect {
	case "mysql":
		sqlColumn = mysqlColumns
	default:
		sqlColumn = postgresColumns
	}
	rows, err := db.Query(sqlColumn)
	if err != nil {
		return nil, fmt.Errorf("error fetching columns: %s", err)
	}
	defer rows.Close()

	di := &DBInfo{
		Dialect: dialect,
		Version: dbVersion,
		Schema:  dbSchema,
		Name:    dbName,
		Tables:  make(map[tableKey]DBTable),
	}

	// we have to rescan and update columns to overcome
	// weird bugs in mysql like joins with information_schema
	// don't work in 8.0.22 etc.
	for rows.Next() {
		var c DBColumn
		err = rows.Scan(
			&c.Schema, &c.Table, &c.Name, &c.Type, &c.NotNull, &c.PrimaryKey,
			&c.UniqueKey, &c.Array, &c.FullText, &c.FKeySchema, &c.FKeyTable, &c.FKeyCol,
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
		tk := tableKey{c.Schema, c.Table}
		t, ok := di.Tables[tk]
		if !ok {
			t = DBTable{
				Schema:   c.Schema,
				Name:     c.Table,
				Columns:  make(map[columnKey]DBColumn),
				FullText: make(map[columnKey]DBColumn),
				//Comment      string
				//Type         string
			}
		}
		if c.PrimaryKey {
			t.PrimaryCol = c
		}
		if isBlock(c.Table, blockList) {
			t.Blocked = true
		}
		ck := columnKey{c.Schema, c.Table, c.Name}
		t.Columns[ck] = c
		if c.FullText {
			t.FullText[ck] = c
		}
		di.Tables[tk] = t
	}
	return di, nil
}

func isBlock(val string, s []string) bool {
	for _, v := range s {
		regex := fmt.Sprintf("^%s$", v)
		if matched, _ := regexp.MatchString(regex, val); matched {
			return true
		}
	}
	return false
}
