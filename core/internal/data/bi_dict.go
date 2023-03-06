package data

import "golang.org/x/exp/slices"

// BiDict is a map fill array value
type BiDict map[string][]string

func (my BiDict) Put(key string, val string) {
	if v, ok := my[key]; !ok || !slices.Contains(v, val) {
		my[key] = append(v, val)
	}

	if v, ok := my[val]; !ok || !slices.Contains(v, key) {
		my[val] = append(v, key)
	}
}

func (my BiDict) Get(key string) ([]string, bool) {
	v, ok := my[key]
	return v, ok
}
