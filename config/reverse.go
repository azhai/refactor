// Copyright 2019 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"os"
	"path/filepath"

	"gitea.com/azhai/refactor/rewrite"
)

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

	Funcs     map[string]string `yaml:"funcs"`
	Formatter string            `yaml:"formatter"`
	Importter string            `yaml:"importter"`
	ExtName   string            `yaml:"ext_name"`

	NameSpace    string `yaml:"name_space"`
	GenJsonTag   bool   `yaml:"gen_json_tag"`
	GenTableName bool   `yaml:"gen_table_name"`
	ApplyMixins  bool   `yaml:"apply_mixins"`
}

func (t ReverseTarget) GetFileName(name string) string {
	_ = os.MkdirAll(t.OutputDir, rewrite.DEFAULT_DIR_MODE)
	return filepath.Join(t.OutputDir, name+t.ExtName)
}

func (t ReverseTarget) MergeOptions(name, tablePrefix string) ReverseTarget {
	if t.Type == "codes" && t.Language == "" {
		t.Language = "golang"
	}
	if name != "" {
		t.OutputDir = filepath.Join(t.OutputDir, name)
	}
	if t.TablePrefix == "" {
		t.TablePrefix = tablePrefix
	}
	return t
}
