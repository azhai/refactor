// Copyright 2019 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package language

import (
	"html/template"
	"strings"

	"gitea.com/azhai/refactor/config"
	"xorm.io/xorm/schemas"
)

type Formatter func(fileName string, sourceCode []byte) ([]byte, error)
type Importter func(tables map[string]*schemas.Table) map[string]string
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

func (l *Language) FixTarget(target *config.ReverseTarget) {
	if target.ExtName == "" && l.ExtName != "" {
		if !strings.HasPrefix(l.ExtName, ".") {
			l.ExtName = "." + l.ExtName
		}
		target.ExtName = l.ExtName
	}
	if target.NameSpace == "" {
		if pck := l.Packager; pck != nil {
			target.NameSpace = pck(target.OutputDir)
		}
		if target.NameSpace == "" {
			target.NameSpace = "models"
		}
	}
}

func NewTemplate(name, content string, funcs template.FuncMap) *template.Template {
	t := template.New(name).Funcs(funcs)
	tmpl, err := t.Parse(content)
	if err != nil {
		panic(err)
	}
	initTemplates[name] = tmpl
	return tmpl
}
