package base

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrDuplicationExists = errors.New("duplication key or value exists")

// BiMap A threadsafe generic bidirectional map structure written in Go
type BiMap[K, V comparable] struct {
	lock   sync.RWMutex
	keys   map[K]V
	values map[V]K
}

type option[K, V comparable] interface {
	apply(*BiMap[K, V])
}

type initialOption[K, V comparable] map[K]V

func (o initialOption[K, V]) apply(m *BiMap[K, V]) {
	for k, v := range map[K]V(o) {
		m.keys[k] = v
		m.values[v] = k
	}
}

// WithInitialMap returns an initialOption object that implements the option interface
func WithInitialMap[K, V comparable](m map[K]V) option[K, V] {
	return initialOption[K, V](m)
}

// NewBiMap returns a BiMap object
func NewBiMap[K, V comparable](options ...option[K, V]) *BiMap[K, V] {
	m := &BiMap[K, V]{
		keys:   make(map[K]V),
		values: make(map[V]K),
	}
	for _, opt := range options {
		opt.apply(m)
	}
	return m
}

func (my *BiMap[K, V]) String() string {
	my.lock.RLock()
	defer my.lock.RUnlock()
	pairs := make([]string, 0, len(my.keys))
	for k, v := range my.keys {
		pairs = append(pairs, fmt.Sprintf("%v:%v", k, v))
	}
	return "map[" + strings.Join(pairs, " ") + "]"
}

// Size returns size of the BiMap.
func (my *BiMap[_, _]) Size() int {
	my.lock.RLock()
	defer my.lock.RUnlock()
	return len(my.keys)
}

// Keys returns a slice of the keys in the BiMap.
func (my *BiMap[K, V]) Keys() []K {
	my.lock.RLock()
	defer my.lock.RUnlock()
	var keys []K
	for k := range my.keys {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of the values in the BiMap.
func (my *BiMap[K, V]) Values() []V {
	my.lock.RLock()
	defer my.lock.RUnlock()
	var values []V
	for v := range my.values {
		values = append(values, v)
	}
	return values
}

// Each iterate over the map for the given function
func (my *BiMap[K, V]) Each(fn func(k K, v V)) {
	my.lock.RLock()
	defer my.lock.RUnlock()
	for k, v := range my.keys {
		fn(k, v)
	}
}

// Put sets the value with corresponding key in the keys map, it will return an error if either key or value exist
func (my *BiMap[K, V]) Put(key K, val V) error {
	my.lock.Lock()
	defer my.lock.Unlock()
	if _, ok := my.keys[key]; ok {
		return ErrDuplicationExists
	}
	if _, ok := my.values[val]; ok {
		return ErrDuplicationExists
	}
	my.keys[key] = val
	my.values[val] = key
	return nil
}

// Delete deletes the key-value pair involving the given key.
func (my *BiMap[T, _]) Delete(key T) {
	my.lock.Lock()
	defer my.lock.Unlock()
	v, ok := my.keys[key]
	if !ok {
		return
	}
	delete(my.keys, key)
	delete(my.values, v)
}

// DeleteInverse deletes the key-value pair involving the given value.
func (my *BiMap[_, V]) DeleteInverse(val V) {
	my.lock.Lock()
	defer my.lock.Unlock()
	k, ok := my.values[val]
	if !ok {
		return
	}
	delete(my.keys, k)
	delete(my.values, val)
}

// Get returns the value and its existence by the given key in keys map
func (my *BiMap[K, V]) Get(key K) (V, bool) {
	my.lock.RLock()
	defer my.lock.RUnlock()
	v, ok := my.keys[key]
	return v, ok
}

// GetInverse returns the key and its existence by the given value in values map
func (my *BiMap[K, V]) GetInverse(v V) (K, bool) {
	my.lock.RLock()
	defer my.lock.RUnlock()
	k, ok := my.values[v]
	return k, ok
}
