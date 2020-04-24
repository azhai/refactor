// Copyright 2019 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package refactor

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"gitea.com/azhai/refactor/config"
	"gitea.com/azhai/refactor/language"
	"gitea.com/azhai/refactor/rewrite"
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
		"Lower":       strings.ToLower,
		"Upper":       strings.ToUpper,
		"Title":       strings.Title,
		"Camelize":    inflect.Camelize,
		"Underscore":  inflect.Underscore,
		"Pluralize":   inflect.Pluralize,
		"Singularize": inflect.Singularize,
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

func Reverse(source *config.DataSource, target *config.ReverseTarget) error {
	formatter := formatters[target.Formatter]
	lang := language.GetLanguage(target.Language)
	if lang != nil {
		lang.FixTarget(target)
		formatter = lang.Formatter
	}
	if formatter == nil {
		formatter = rewrite.WriteCodeFile
	}

	isRedis := true
	if source.ReverseSource.Database != "redis" {
		isRedis = false
		err := RunReverse(source.ReverseSource, target)
		if err != nil {
			return err
		}
	}
	if target.Language != "golang" {
		return nil
	}
	var _err error
	if target.ApplyMixins {
		_err = rewrite.ScanModelDir(target.OutputDir, true)
	}

	var tmpl *template.Template
	if isRedis {
		tmpl = language.GetGolangTemplate("cache")
	} else {
		tmpl = language.GetGolangTemplate("conn")
	}
	buf := new(bytes.Buffer)
	_ = tmpl.Execute(buf, map[string]interface{}{
		"Target":    target,
		"NameSpace": target.NameSpace,
		"ConnKey":   source.ConnKey,
	})
	fileName := target.GetFileName("init")
	_, err := formatter(fileName, buf.Bytes())
	if err == nil {
		err = _err
	}
	return err
}

func RunReverse(source *config.ReverseSource, target *config.ReverseTarget) error {
	orm, err := xorm.NewEngine(source.Database, source.ConnStr)
	if err != nil {
		return err
	}

	tableSchemas, err := orm.DBMetas()
	if err != nil {
		return err
	}

	// filter tables according includes and excludes
	tableSchemas = filterTables(tableSchemas, target)

	// load configuration from language
	lang := language.GetLanguage(target.Language)
	funcs := newFuncs()
	formatter := formatters[target.Formatter]
	importter := importters[target.Importter]

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

	tables := make(map[string]*schemas.Table)
	for _, table := range tableSchemas {
		tableName := table.Name
		if target.TablePrefix != "" {
			table.Name = strings.TrimPrefix(table.Name, target.TablePrefix)
		}
		for _, col := range table.Columns() {
			col.FieldName = colMapper.Table2Obj(col.Name)
		}
		tables[tableName] = table
	}

	err = os.MkdirAll(target.OutputDir, os.ModePerm)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if !target.MultipleFiles {
		packages := importter(tables)
		if err = tmpl.Execute(buf, map[string]interface{}{
			"Target":  target,
			"Tables":  tables,
			"Imports": packages,
		}); err != nil {
			return err
		}
		fileName := target.GetFileName("models")
		if _, err = formatter(fileName, buf.Bytes()); err != nil {
			return err
		}
	} else {
		for tableName, table := range tables {
			tbs := map[string]*schemas.Table{tableName: table}
			packages := importter(tbs)
			buf.Reset()
			if err = tmpl.Execute(buf, map[string]interface{}{
				"Target":  target,
				"Tables":  tbs,
				"Imports": packages,
			}); err != nil {
				return err
			}
			fileName := target.GetFileName(table.Name)
			if _, err = formatter(fileName, buf.Bytes()); err != nil {
				return err
			}
		}
	}
	return nil
}
