package test

import (
	"github.com/bytedance/sonic"
	"testing"
)

type Object struct {
	Fields []Field `json:"fields,omitempty"`
}
type Field struct {
}

func TestSonic(t *testing.T) {
	data := Object{Fields: []Field{}}
	output, err := sonic.Marshal(&data)
	if err != nil {
		return
	}
	println(string(output))
}
