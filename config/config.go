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

func (cfg Settings) GetSource(name string) (ReverseSource, ConnConfig) {
	var ok bool
	src, c := ReverseSource{}, ConnConfig{}
	if c, ok = cfg.Connections[name]; !ok {
		return src, c
	}
	d := dialect.GetDialectByName(c.DriverName)
	if d != nil {
		src.Database = d.Name()
		src.ConnStr = d.GetDSN(c.Params)
	}
	return src, c
}
