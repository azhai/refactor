// Copyright 2019 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gitea.com/azhai/refactor/config"
	"gitea.com/azhai/refactor/language"
	"github.com/gobwas/glob"
	"github.com/grsmv/inflect"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
	"xorm.io/xorm/schemas"
)

var (
	formatters   = map[string]language.Formatter{}
	importters   = map[string]language.Importter{}
	defaultFuncs = template.FuncMap{
		"Lower":      strings.ToLower,
		"Upper":      strings.ToUpper,
		"Title":      strings.Title,
		"Camelize":   inflect.Camelize,
		"Underscore": inflect.Underscore,
	}
)

func filterTables(tables []*schemas.Table, target *config.ReverseTarget) []*schemas.Table {
	var res = make([]*schemas.Table, 0, len(tables))
	for _, tb := range tables {
		var remove bool
		for _, exclude := range target.ExcludeTables {
			s, _ := glob.Compile(exclude)
			remove = s.Match(tb.Name)
			if remove {
				break
			}
		}
		if remove {
			continue
		}
		if len(target.IncludeTables) == 0 {
			res = append(res, tb)
			continue
		}

		var keep bool
		for _, include := range target.IncludeTables {
			s, _ := glob.Compile(include)
			keep = s.Match(tb.Name)
			if keep {
				break
			}
		}
		if keep {
			res = append(res, tb)
		}
	}
	return res
}

func newFuncs() template.FuncMap {
	var m = make(template.FuncMap)
	for k, v := range defaultFuncs {
		m[k] = v
	}
	return m
}

func convertMapper(mapname string) names.Mapper {
	switch mapname {
	case "gonic":
		return names.LintGonicMapper
	case "same":
		return names.SameMapper{}
	default:
		return names.SnakeMapper{}
	}
}

func RunReverse(source *config.ReverseSource, target *config.ReverseTarget) error {
	orm, err := xorm.NewEngine(source.Database, source.ConnStr)
	if err != nil {
		return err
	}

	tables, err := orm.DBMetas()
	if err != nil {
		return err
	}

	// filter tables according includes and excludes
	tables = filterTables(tables, target)

	// load configuration from language
	lang := language.GetLanguage(target.Language)
	funcs := newFuncs()
	formatter := formatters[target.Formatter]
	importter := importters[target.Importter]
	var packager language.Packager

	// load template
	var bs []byte
	if target.Template != "" {
		bs = []byte(target.Template)
	} else if target.TemplatePath != "" {
		bs, err = ioutil.ReadFile(target.TemplatePath)
		if err != nil {
			return err
		}
	}

	if lang != nil {
		if bs == nil {
			bs = []byte(lang.Template)
		}
		for k, v := range lang.Funcs {
			funcs[k] = v
		}
		if formatter == nil {
			formatter = lang.Formatter
		}
		if importter == nil {
			importter = lang.Importter
		}
		packager = lang.Packager
		target.ExtName = lang.ExtName
	}
	if !strings.HasPrefix(target.ExtName, ".") {
		target.ExtName = "." + target.ExtName
	}

	var tableMapper = convertMapper(target.TableMapper)
	var colMapper = convertMapper(target.ColumnMapper)

	funcs["TableMapper"] = tableMapper.Table2Obj
	funcs["ColumnMapper"] = colMapper.Table2Obj

	if bs == nil {
		return errors.New("You have to indicate template / template path or a language")
	}

	t := template.New("reverse")
	t.Funcs(funcs)

	tmpl, err := t.Parse(string(bs))
	if err != nil {
		return err
	}

	for _, table := range tables {
		if target.TablePrefix != "" {
			table.Name = strings.TrimPrefix(table.Name, target.TablePrefix)
		}
		for _, col := range table.Columns() {
			col.FieldName = colMapper.Table2Obj(col.Name)
		}
	}

	err = os.MkdirAll(target.OutputDir, os.ModePerm)
	if err != nil {
		return err
	}

	nameSpace := "models"
	if packager != nil {
		nameSpace = packager(target.OutputDir)
	}
	if !target.MultipleFiles {
		packages := importter(tables)

		newbytes := bytes.NewBufferString("")
		err = tmpl.Execute(newbytes, map[string]interface{}{
			"NameSpace": nameSpace,
			"Target":    target,
			"Tables":    tables,
			"Imports":   packages,
		})
		if err != nil {
			return err
		}

		sourceCode, err := ioutil.ReadAll(newbytes)
		if err != nil {
			return err
		}

		fileName := filepath.Join(target.OutputDir, "models"+target.ExtName)
		if _, err = formatter(fileName, sourceCode); err != nil {
			return err
		}
	} else {
		for _, table := range tables {
			// imports
			tbs := []*schemas.Table{table}
			packages := importter(tbs)

			newbytes := bytes.NewBufferString("")
			err = tmpl.Execute(newbytes, map[string]interface{}{
				"NameSpace": nameSpace,
				"Target":    target,
				"Tables":    tbs,
				"Imports":   packages,
			})
			if err != nil {
				return err
			}

			sourceCode, err := ioutil.ReadAll(newbytes)
			if err != nil {
				return err
			}

			fileName := filepath.Join(target.OutputDir, table.Name+target.ExtName)
			if _, err = formatter(fileName, sourceCode); err != nil {
				return err
			}
		}
	}

	return nil
}
