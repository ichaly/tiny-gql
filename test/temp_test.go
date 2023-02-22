package test

import (
	"strings"
	"testing"
)

func TestIndexRune(t *testing.T) {
	str := "adfasdfsa(dfasd)"
	if i := strings.IndexRune(str, '('); i != -1 {
		str = str[:i]
	}
	println(str)
}
