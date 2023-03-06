package data

import (
	"github.com/bytedance/sonic"
	"testing"
)

var d BiDict

func init() {
	d = map[string][]string{}
}

func TestBiDict(t *testing.T) {
	d.Put("key1", "val1")
	d.Put("key1", "val2")
	d.Put("key1", "val1")
	v, ok := d.Get("key1")
	if ok && len(v) != 2 {
		t.Errorf("The values of %v is not %v\n", len(v), 2)
	}
	output, err := sonic.Marshal(&d)
	if err != nil {
		return
	}
	println(string(output))
}
