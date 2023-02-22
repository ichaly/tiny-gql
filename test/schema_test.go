package test

import (
	"github.com/ichaly/tiny-go/core"
	"testing"
)

func TestIntrospection(t *testing.T) {
	in, err := core.Introspection()
	if err == nil {
		println(string(in))
	} else {
		panic(err)
	}
}
