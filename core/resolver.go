package core

type ResolverConfig struct {
	Name      string
	Type      string
	Schema    string
	Table     string
	Column    string
	StripPath string        `mapstructure:"strip_path" json:"strip_path" yaml:"strip_path"`
	Props     ResolverProps `mapstructure:",remain"`
}

type ResolverProps map[string]interface{}
