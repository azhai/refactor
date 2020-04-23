// Copyright 2019 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import "path/filepath"

// ReverseConfig represents a reverse configuration
type ReverseConfig struct {
	Kind    string          `yaml:"kind"`
	Name    string          `yaml:"name"`
	Source  ReverseSource   `yaml:"source"`
	Targets []ReverseTarget `yaml:"targets"`
}

// ReverseSource represents a reverse source which should be a database connection
type ReverseSource struct {
	Database string `yaml:"database"`
	ConnStr  string `yaml:"conn_str"`
}

// ReverseTarget represents a reverse target
type ReverseTarget struct {
	Type          string   `yaml:"type"`
	IncludeTables []string `yaml:"include_tables"`
	ExcludeTables []string `yaml:"exclude_tables"`
	TableMapper   string   `yaml:"table_mapper"`
	ColumnMapper  string   `yaml:"column_mapper"`
	TemplatePath  string   `yaml:"template_path"`
	Template      string   `yaml:"template"`
	MultipleFiles bool     `yaml:"multiple_files"`
	OutputDir     string   `yaml:"output_dir"`
	TablePrefix   string   `yaml:"table_prefix"`
	Language      string   `yaml:"language"`

	GenJsonTag   bool              `yaml:"gen_json_tag"`
	GenTableName bool              `yaml:"gen_table_name"`
	Funcs        map[string]string `yaml:"funcs"`
	Formatter    string            `yaml:"formatter"`
	Importter    string            `yaml:"importter"`
	ExtName      string            `yaml:"ext_name"`
}

func (t ReverseTarget) FixTarget(name, tablePrefix string) ReverseTarget {
	if t.Type == "codes" && t.Language == "" {
		t.Language = "golang"
	}
	if t.TablePrefix == "" {
		t.TablePrefix = tablePrefix
	}
	if name != "" {
		t.OutputDir = filepath.Join(t.OutputDir, name)
	}
	return t
}
