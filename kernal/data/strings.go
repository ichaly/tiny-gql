package data

import (
	"fmt"
	"strings"
)

func (my *DBTable) String() string {
	return my.Schema + "." + my.Name
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
