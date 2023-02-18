package test

import (
	"github.com/ichaly/tiny-go/core"
	"path/filepath"
	"testing"
)

func TestReadInConfig(t *testing.T) {
	cfg, err := core.ReadInConfig(filepath.Join("./cfg", "prod.yml"))
	if err != nil {
		return
	}
	str, err := json.MarshalToString(cfg)
	if err != nil {
		return
	}
	println(str)
}
