package config

import (
	"fmt"
	"os"

	"gitea.com/azhai/refactor/config/dialect"
	"github.com/k0kubun/pp"
	"gopkg.in/yaml.v2"
	"xorm.io/xorm"
)

const (
	DEFAULT_DIR_MODE  = 0o755
	DEFAULT_FILE_MODE = 0o644
)

type IConnectSettings interface {
	GetConnConfigMap(keys ...string) map[string]ConnConfig
	GetConnConfig(key string) (ConnConfig, bool)
}

type IReverseSettings interface {
	GetReverseTargets() []ReverseTarget
	IConnectSettings
}

type Settings struct {
	Application    AppConfig             `json:"application" yaml:"application"`
	Connections    map[string]ConnConfig `json:"connections" yaml:"connections"`
	ReverseTargets []ReverseTarget       `json:"reverse_targets" yaml:"reverse_targets"`
}

type AppConfig struct {
	Debug       bool `json:"debug" yaml:"debug"`
	PluralTable bool `json:"plural_table" yaml:"plural_table"`
}

type PartConfig struct {
	TablePrefix   string   `json:"table_prefix" yaml:"table_prefix"`
	IncludeTables []string `json:"include_tables" yaml:"include_tables"`
	ExcludeTables []string `json:"exclude_tables" yaml:"exclude_tables"`
}

type ConnConfig struct {
	DriverName string                          `json:"driver_name" yaml:"driver_name"`
	ReadOnly   bool                            `json:"read_only" yaml:"read_only"`
	Params     dialect.ConnParams              `json:"params" yaml:"params"`
	PartConfig `json:",inline" yaml:",inline"` // 注意逗号不能少
}

type DataSource struct {
	ConnKey      string
	ImporterPath string
	Dialect      dialect.Dialect
	PartConfig
	*ReverseSource
}

func ReadSettings(fileName string) (*Settings, error) {
	cfg := new(Settings)
	rd, err := os.Open(fileName)
	if err == nil {
		err = yaml.NewDecoder(rd).Decode(&cfg)
	}
	return cfg, err
}

func SaveSettings(fileName string, cfg interface{}) error {
	wt, err := os.Open(fileName)
	if err == nil {
		err = yaml.NewEncoder(wt).Encode(cfg)
	}
	return err
}

func (cfg Settings) GetReverseTargets() []ReverseTarget {
	return cfg.ReverseTargets
}

func (cfg Settings) GetConnConfigMap(keys ...string) map[string]ConnConfig {
	if len(keys) == 0 {
		return cfg.Connections
	}
	result := make(map[string]ConnConfig)
	for _, k := range keys {
		if c, ok := cfg.Connections[k]; ok {
			result[k] = c
		}
	}
	return result
}

func (cfg Settings) GetConnConfig(key string) (ConnConfig, bool) {
	if c, ok := cfg.Connections[key]; ok {
		return c, true
	}
	return ConnConfig{}, false
}

func (c ConnConfig) Connect(verbose bool) (*xorm.Engine, error) {
	d := dialect.GetDialectByName(c.DriverName)
	drv, dsn := d.Name(), d.ParseDSN(c.Params)
	if verbose {
		pp.Printf("Connect: %s %s\n", drv, dsn)
	}
	engine, err := xorm.NewEngine(drv, dsn)
	if err == nil {
		engine.ShowSQL(verbose)
	}
	return engine, err
}

func NewDataSource(c ConnConfig, name string) *DataSource {
	d := &DataSource{ConnKey: name, PartConfig: c.PartConfig}
	d.Dialect = dialect.GetDialectByName(c.DriverName)
	if d.Dialect != nil {
		d.ImporterPath = d.Dialect.ImporterPath()
		d.ReverseSource = &ReverseSource{
			Database: d.Dialect.Name(),
			ConnStr:  d.Dialect.ParseDSN(c.Params),
		}
	}
	return d
}

func (ds *DataSource) Connect(verbose bool) (*xorm.Engine, error) {
	if ds.Database == "" || ds.ConnStr == "" {
		return nil, fmt.Errorf("the config of connection is empty")
	} else if verbose {
		pp.Println(ds.Database, ds.ConnStr)
	}
	engine, err := xorm.NewEngine(ds.Database, ds.ConnStr)
	if err == nil {
		engine.ShowSQL(verbose)
	}
	return engine, err
}
