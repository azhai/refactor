// Copyright 2019 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package language

import (
	"html/template"
	"xorm.io/xorm/schemas"
)

type Formatter func(fileName string, sourceCode []byte) ([]byte, error)
type Importter func(tables []*schemas.Table) []string
type Packager func(targetDir string) string

// Language represents a languages supported when reverse codes
type Language struct {
	Name      string
	ExtName   string
	Template  string
	Types     map[string]string
	Funcs     template.FuncMap
	Formatter Formatter
	Importter Importter
	Packager  Packager
}

var (
	languages = make(map[string]*Language)
)

// RegisterLanguage registers a language
func RegisterLanguage(l *Language) {
	languages[l.Name] = l
}

// GetLanguage returns a language if exists
func GetLanguage(name string) *Language {
	return languages[name]
}
