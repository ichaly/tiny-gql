package schema

func NamedType(named string) *Type {
	return &Type{Name: named, NonNull: false}
}

func NonNull(named string) *Type {
	return &Type{Name: named, NonNull: true}
}

func ListType(elem *Type) *Type {
	return &Type{Elem: elem, NonNull: false}
}

type Type struct {
	Name    string
	Elem    *Type
	NonNull bool
}

func (t *Type) Strip() string {
	if t.Name != "" {
		return t.Name
	}

	return t.Elem.Strip()
}

func (t *Type) String() string {
	s := ""
	if t.NonNull {
		s = "!"
	}
	if t.Name != "" {
		return t.Name + s
	}
	return "[" + t.Elem.String() + "]" + s
}

func (t *Type) IsCompatible(other *Type) bool {
	if t.Name != other.Name {
		return false
	}

	if t.Elem != nil && other.Elem == nil {
		return false
	}

	if t.Elem != nil && !t.Elem.IsCompatible(other.Elem) {
		return false
	}

	if other.NonNull {
		return t.NonNull
	}

	return true
}
