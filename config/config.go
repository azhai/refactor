package config

import (
	"fmt"
	"os"

	"gitea.com/azhai/refactor/config/dialect"
	"gopkg.in/yaml.v2"
)

var cfg *Settings

type Settings struct {
	isEmpty        bool
	Application    AppConfig             `json:"application" yaml:"application"`
	Connections    map[string]ConnConfig `json:"connections" yaml:"connections"`
	ReverseTargets []ReverseTarget       `json:"reverse_targets" yaml:"reverse_targets"`
}

type AppConfig struct {
	Debug       bool `json:"debug" yaml:"debug"`
	PluralTable bool `json:"plural_table" yaml:"plural_table"`
}

type ConnConfig struct {
	DriverName  string `json:"driver_name" yaml:"driver_name"`
	TablePrefix string `json:"table_prefix" yaml:"table_prefix"`
	ReadOnly    string `json:"read_only" yaml:"read_only"`
	Params      dialect.ConnParams
}

type DataSource struct {
	ConnKey     string
	TablePrefix string
	Dialect     dialect.Dialect
	*ReverseSource
}

func GetSettings() *Settings {
	if cfg == nil {
		cfg = new(Settings)
		cfg.isEmpty = true
	}
	return cfg
}

func ReadSettings(file string) (*Settings, error) {
	cfg = new(Settings)
	rd, err := os.Open(file)
	if err == nil {
		err = yaml.NewDecoder(rd).Decode(&cfg)
		if err == nil {
			cfg.isEmpty = false
		}
	}
	return cfg, err
}

func SaveSettings(file string) error {
	if cfg = GetSettings(); cfg.isEmpty {
		return fmt.Errorf("the settings is not exists")
	}
	wt, err := os.Open(file)
	if err == nil {
		err = yaml.NewEncoder(wt).Encode(cfg)
	}
	return err
}

func (cfg *Settings) GetDataSources(names []string) (ds []*DataSource) {
	if len(names) == 0 {
		for name, c := range cfg.Connections {
			ds = append(ds, NewDataSource(name, c))
		}
		return
	} else {
		for _, name := range names {
			c, ok := cfg.Connections[name]
			if !ok {
				continue
			}
			ds = append(ds, NewDataSource(name, c))
		}
		return
	}
}

func NewDataSource(k string, c ConnConfig) *DataSource {
	d := &DataSource{ConnKey: k, TablePrefix: c.TablePrefix}
	d.Dialect = dialect.GetDialectByName(c.DriverName)
	if d.Dialect != nil {
		d.ReverseSource = &ReverseSource{
			Database: d.Dialect.Name(),
			ConnStr:  d.Dialect.GetDSN(c.Params),
		}
	}
	return d
}
