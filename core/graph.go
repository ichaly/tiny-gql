package core

type Graph map[string][]string

func (my Graph) Put(key string, val string) {
	my[key] = append(my[key], val)
	my[val] = append(my[val], key)
}

func (my Graph) Get(key string) ([]string, bool) {
	v, ok := my[key]
	return v, ok
}
