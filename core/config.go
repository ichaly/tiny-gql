package core

import (
	"fmt"
	"github.com/ichaly/tiny-go/core/internal/util"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Debug     bool             `jsonschema:"title=Debug,default=false"`
	Tables    []Table          `jsonschema:"title=Tables"`
	Resolvers []ResolverConfig `jsonschema:"-"`

	ConfigPath string `mapstructure:"config_path" jsonschema:"title=Config Path"`
}

type Table struct {
	Schema    string
	Table     string // Inherits Table
	Name      string
	Type      string
	Columns   []Column
	Blocklist []string
	// Permitted order by options
	OrderBy map[string][]string `mapstructure:"order_by" json:"order_by" yaml:"order_by" jsonschema:"title=Order By Options,example=created_at desc"`

	Query  *Query
	Insert *Insert
	Update *Update
	Upsert *Upsert
	Delete *Delete
}

type Column struct {
	Name       string
	Type       string `jsonschema:"example=integer,example=text"`
	Array      bool
	Primary    bool
	ForeignKey string `mapstructure:"related_to" json:"related_to" yaml:"related_to" jsonschema:"title=Related To,example=other_table.id_column,example=users.id"`
}

type Query struct {
	Limit int
	// Use filters to enforce table wide things like { disabled: false } where you never want disabled users to be shown.
	Filters          []string
	Columns          []string
	DisableFunctions bool `mapstructure:"disable_functions" json:"disable_functions" yaml:"disable_functions"`
	Block            bool
}

type Insert struct {
	Filters []string
	Columns []string
	Presets map[string]string
	Block   bool
}

type Update struct {
	Filters []string
	Columns []string
	Presets map[string]string
	Block   bool
}

type Upsert struct {
	Filters []string
	Columns []string
	Presets map[string]string
	Block   bool
}

type Delete struct {
	Filters []string
	Columns []string
	Block   bool
}

func ReadInConfig(configFile string) (*Config, error) {
	return readInConfig(configFile, nil)
}

func readInConfig(configFile string, fs afero.Fs) (*Config, error) {
	cp := filepath.Dir(configFile)
	vi := newViper(cp, filepath.Base(configFile))

	if fs != nil {
		vi.SetFs(fs)
	}

	if err := vi.ReadInConfig(); err != nil {
		return nil, err
	}

	if pcf := vi.GetString("inherits"); pcf != "" {
		cf := vi.ConfigFileUsed()
		vi = newViper(cp, pcf)
		if fs != nil {
			vi.SetFs(fs)
		}

		if err := vi.ReadInConfig(); err != nil {
			return nil, err
		}

		if v := vi.GetString("inherits"); v != "" {
			return nil, fmt.Errorf("inherited config '%s' cannot itself inherit '%s'", pcf, v)
		}

		vi.SetConfigFile(cf)

		if err := vi.MergeInConfig(); err != nil {
			return nil, err
		}
	}

	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "GJ_") || strings.HasPrefix(e, "SJ_") {
			kv := strings.SplitN(e, "=", 2)
			util.SetKeyValue(vi, kv[0], kv[1])
		}
	}

	c := &Config{}
	c.ConfigPath = cp

	if err := vi.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("failed to decode config, %v", err)
	}

	return c, nil
}

func newViper(configPath, configFile string) *viper.Viper {
	vi := newViperWithDefaults()
	vi.SetConfigName(strings.TrimSuffix(configFile, filepath.Ext(configFile)))

	if configPath == "" {
		vi.AddConfigPath("./config")
	} else {
		vi.AddConfigPath(configPath)
	}

	return vi
}

func newViperWithDefaults() *viper.Viper {
	vi := viper.New()

	vi.SetDefault("host_port", "0.0.0.0:8080")
	vi.SetDefault("web_ui", false)
	vi.SetDefault("enable_tracing", false)
	vi.SetDefault("auth_fail_block", false)
	vi.SetDefault("seed_file", "seed.js")

	vi.SetDefault("log_level", "info")
	vi.SetDefault("log_format", "json")

	vi.SetDefault("default_block", true)

	vi.SetDefault("database.type", "postgres")
	vi.SetDefault("database.host", "localhost")
	vi.SetDefault("database.port", 5432)
	vi.SetDefault("database.user", "postgres")
	vi.SetDefault("database.password", "")
	vi.SetDefault("database.schema", "public")
	vi.SetDefault("database.pool_size", 10)

	vi.SetDefault("env", "development")

	vi.BindEnv("env", "GO_ENV") //nolint:errcheck
	vi.BindEnv("host", "HOST")  //nolint:errcheck
	vi.BindEnv("port", "PORT")  //nolint:errcheck

	vi.SetDefault("auth.rails.max_idle", 80)
	vi.SetDefault("auth.rails.max_active", 12000)
	vi.SetDefault("auth.subs_creds_in_vars", false)

	return vi
}
