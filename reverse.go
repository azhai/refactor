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
	"github.com/azhai/gozzo-utils/filesystem"
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
		"Lower":         strings.ToLower,
		"Upper":         strings.ToUpper,
		"Title":         strings.Title,
		"Camelize":      inflect.Camelize,
		"Underscore":    inflect.Underscore,
		"Singularize":   inflect.Singularize,
		"Pluralize":     inflect.Pluralize,
		"DiffPluralize": DiffPluralize,
		"GetSinglePKey": GetSinglePKey,
	}
)

// 如果复数形式和单数相同，人为增加后缀
func DiffPluralize(word, suffix string) string {
	words := inflect.Pluralize(word)
	if words == word {
		words += suffix
	}
	return words
}

func GetSinglePKey(table *schemas.Table) (pkey string) {
	if cols := table.PKColumns(); len(cols) == 1 {
		pkey = cols[0].FieldName
	}
	return
}

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

	var tmpl *template.Template
	if isRedis {
		tmpl = language.GetGolangTemplate("cache", nil)
	} else {
		tmpl = language.GetGolangTemplate("conn", nil)
	}
	buf := new(bytes.Buffer)
	data := map[string]interface{}{
		"Target":    target,
		"NameSpace": target.NameSpace,
		"ConnKey":   source.ConnKey,
	}
	if err := tmpl.Execute(buf, data); err != nil {
		return err
	}
	fileName := target.GetFileName(config.INIT_FILE_NAME)
	_, err := formatter(fileName, buf.Bytes())

	if target.ApplyMixins {
		files, _ := filesystem.FindFiles("./language/common/", ".go")
		baseImport := "gitea.com/azhai/refactor/language/common"
		for fileName := range files {
			_ = rewrite.AddFormerMixins(fileName, baseImport, "base")
		}
		if target.MixinDirPath != "" {
			files, _ := filesystem.FindFiles(target.MixinDirPath, ".go")
			for fileName := range files {
				_ = rewrite.AddFormerMixins(fileName, target.MixinNameSpace, "")
			}
		}
		files, _ = filesystem.FindFiles(target.OutputDir, ".go")
		for fileName := range files {
			_err := rewrite.ParseAndMixinFile(fileName, true)
			if _err != nil {
				err = _err
			}
		}
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
	if lang != nil {
		bs = []byte(lang.Template)
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

	// 配置模板优先于语言模板
	tmplQuery := language.GetGolangTemplate("query", funcs)
	if target.QueryTemplatePath != "" {
		qt, err := ioutil.ReadFile(target.QueryTemplatePath)
		if err == nil || len(qt) > 0 {
			tmplQuery = language.NewTemplate("custom-query", string(qt), funcs)
		} else {
			target.GenQueryMethods, tmplQuery = false, nil
		}
	}
	if target.TemplatePath != "" {
		bs, err = ioutil.ReadFile(target.TemplatePath)
		if err != nil {
			return err
		}
	}

	if bs == nil {
		return errors.New("you have to indicate template / template path or a language")
	}
	tmpl := language.NewTemplate("custom-model", string(bs), funcs)
	queryImports := map[string]string{"xorm.io/xorm":""}

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
		data := map[string]interface{}{
			"Target":  target,
			"Tables":  tables,
			"Imports": packages,
		}
		if err = tmpl.Execute(buf, data); err != nil {
			return err
		}
		fileName := target.GetFileName(config.SINGLE_FILE_NAME)
		if _, err = formatter(fileName, buf.Bytes()); err != nil {
			return err
		}
		if target.GenQueryMethods && tmplQuery != nil {
			buf.Reset()
			data["Imports"] = queryImports
			if err = tmplQuery.Execute(buf, data); err != nil {
				return err
			}
			fileName := target.GetFileName(config.QUERY_FILE_NAME)
			if _, err = formatter(fileName, buf.Bytes()); err != nil {
				return err
			}
		}
	} else {
		for tableName, table := range tables {
			tbs := map[string]*schemas.Table{tableName: table}
			packages := importter(tbs)
			data := map[string]interface{}{
				"Target":  target,
				"Tables":  tbs,
				"Imports": packages,
			}
			buf.Reset()
			if err = tmpl.Execute(buf, data); err != nil {
				return err
			}
			if target.GenQueryMethods && tmplQuery != nil {
				data["Imports"] = queryImports
				if err = tmplQuery.Execute(buf, data); err != nil {
					return err
				}
			}
			fileName := target.GetFileName(table.Name)
			if _, err = formatter(fileName, buf.Bytes()); err != nil {
				return err
			}
		}
	}
	return nil
}

func ExecReverseSettings(cfg config.IReverseSettings, names ...string) error {
	conns := cfg.GetConnections(names...)
	for key, conf := range conns {
		d := config.NewDataSource(key, conf)
		if d.ReverseSource == nil {
			continue
		}
		for _, target := range cfg.GetReverseTargets() {
			target = target.MergeOptions(d.ConnKey, d.PartConfig)
			if err := Reverse(d, &target); err != nil {
				return err
			}
		}
	}
	return nil
}
