package data

type DBInfo struct {
	Dialect string
	Version int
	Schema  string
	Name    string

	Tables    []DBTable
	Functions []DBFunction
	VTables   []VirtualTable `json:"-"` // for polymorphic relationships
	colMap    map[string]int
	tableMap  map[string]int
	hash      int
}

type VirtualTable struct {
	Name       string
	IDColumn   string
	TypeColumn string
	FKeyColumn string
}

type DBTable struct {
	Comment      string
	Schema       string
	Name         string
	Type         string
	Columns      []DBColumn
	PrimaryCol   DBColumn
	SecondaryCol DBColumn
	FullText     []DBColumn
	Blocked      bool
	Func         DBFunction
	colMap       map[string]int
}

type DBColumn struct {
	Comment     string
	ID          int32
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

type DBFunction struct {
	Comment string
	Schema  string
	Name    string
	Type    string
	Agg     bool
	Inputs  []DBFuncParam
	Outputs []DBFuncParam
}

type DBFuncParam struct {
	ID    int
	Name  string
	Type  string
	Array bool
}
