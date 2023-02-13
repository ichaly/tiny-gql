package test

import (
	"github.com/ichaly/tiny-go/kernal"
	"strconv"
	"testing"
)

var m *kernal.BiMap[string, int]

func TestMain(t *testing.M) {
	m = kernal.NewBiMap[string, int]()
	for i := 0; i < 10; i++ {
		_ = m.Put(strconv.Itoa(i), i)
	}
	t.Run()
}

func TestPut(t *testing.T) {
	if err := m.Put("Key_${i}", 11); err != nil {
		t.Errorf("expected nil ,but %v got", err)
	}

	if err := m.Put("1", 1); err == nil {
		t.Errorf("expected err ,but got nil")
	}
}
func TestGet(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		if v, ok := m.Get("1"); ok {
			if v != 1 {
				t.Errorf("expected 1 ,but %v got", v)
			}
		} else {
			t.Errorf("expected true ,but %v got", ok)
		}
	})

	t.Run("GetInverse", func(t *testing.T) {
		if v, ok := m.GetInverse(1); ok {
			if v != "1" {
				t.Errorf("expected \"1\" ,but %v got", v)
			}
		} else {
			t.Errorf("expected true ,but %v got", ok)
		}
	})
}

func TestDelete(t *testing.T) {
	m.Delete("1")
	m.DeleteInverse(2)
}
