package core

// AMap is a map fill array value
type AMap map[string][]string

func (my AMap) Put(key string, val string) {
	my[key] = append(my[key], val)
	my[val] = append(my[val], key)
}

func (my AMap) Get(key string) ([]string, bool) {
	v, ok := my[key]
	return v, ok
}
