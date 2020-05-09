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
	GetConnections(keys ...string) map[string]ConnConfig
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
	ConnKey string
	Dialect dialect.Dialect
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

func (cfg Settings) GetConnections(keys ...string) map[string]ConnConfig {
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

func NewDataSource(k string, c ConnConfig) *DataSource {
	d := &DataSource{ConnKey: k, PartConfig: c.PartConfig}
	d.Dialect = dialect.GetDialectByName(c.DriverName)
	if d.Dialect != nil {
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
