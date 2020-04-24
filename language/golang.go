// Copyright 2019 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package language

import (
	"errors"
	"fmt"
	"html/template"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"gitea.com/azhai/refactor/rewrite"
	"xorm.io/xorm/schemas"
)

// Golang represents a golang language
var Golang = Language{
	Name:     "golang",
	Template: golangModelTemplate,
	Types:    map[string]string{},
	Funcs: template.FuncMap{
		"Type": typestring,
		"Tag":  tag,
	},
	Formatter: rewrite.WriteCodeFile,
	Importter: genGoImports,
	Packager:  genNameSpace,
	ExtName:   ".go",
}

func init() {
	RegisterLanguage(&Golang)
}

var (
	errBadComparisonType = errors.New("invalid type for comparison")
	errBadComparison     = errors.New("incompatible types for comparison")
	errNoComparison      = errors.New("missing argument for comparison")
)

type kind int

const (
	invalidKind kind = iota
	boolKind
	complexKind
	intKind
	floatKind
	integerKind
	stringKind
	uintKind
)

func basicKind(v reflect.Value) (kind, error) {
	switch v.Kind() {
	case reflect.Bool:
		return boolKind, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intKind, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintKind, nil
	case reflect.Float32, reflect.Float64:
		return floatKind, nil
	case reflect.Complex64, reflect.Complex128:
		return complexKind, nil
	case reflect.String:
		return stringKind, nil
	}
	return invalidKind, errBadComparisonType
}

func getCol(cols map[string]*schemas.Column, name string) *schemas.Column {
	return cols[strings.ToLower(name)]
}

func genNameSpace(targetDir string) string {
	nameSpace := filepath.Base(targetDir)
	if nameSpace == "default" { // it is golang keyword
		nameSpace = "db"
	}
	return nameSpace
}

func genGoImports(tables map[string]*schemas.Table) []string {
	imports := make(map[string]string)
	results := make([]string, 0)
	for _, table := range tables {
		for _, col := range table.Columns() {
			if typestring(col) == "time.Time" {
				if _, ok := imports["time"]; !ok {
					imports["time"] = "time"
					results = append(results, "time")
				}
			}
		}
	}
	return results
}

func typestring(col *schemas.Column) string {
	st := col.SQLType
	t := schemas.SQLType2Type(st)
	s := t.String()
	if s == "[]uint8" {
		return "[]byte"
	}
	return s
}

func tag(table *schemas.Table, col *schemas.Column, genJson bool) template.HTML {
	tj, tx := "", tagXorm(table, col)
	if genJson {
		tj = tagJson(col)
	} else {
		return template.HTML(tx)
	}
	if tx == "" {
		if tj == "" {
			return ""
		} else {
			return template.HTML(tj)
		}
	}
	return template.HTML(tj + " " + tx)
}

func tagJson(col *schemas.Column) string {
	if col.Name == "" {
		return ""
	}
	return fmt.Sprintf(`json:"%s"`, col.Name)
}

func tagXorm(table *schemas.Table, col *schemas.Column) string {
	isNameId := col.FieldName == "Id"
	isIdPk := isNameId && typestring(col) == "int64"

	var res []string
	if !col.Nullable {
		if !isIdPk {
			res = append(res, "not null")
		}
	}
	if col.IsPrimaryKey {
		res = append(res, "pk")
	}
	if col.Default != "" {
		res = append(res, "default "+col.Default)
	}
	if col.IsAutoIncrement {
		res = append(res, "autoincr")
	}

	/*if col.SQLType.IsTime() && include(created, col.Name) {
		res = append(res, "created")
	}

	if col.SQLType.IsTime() && include(updated, col.Name) {
		res = append(res, "updated")
	}

	if col.SQLType.IsTime() && include(deleted, col.Name) {
		res = append(res, "deleted")
	}*/

	if /*supportComment &&*/ col.Comment != "" {
		res = append(res, fmt.Sprintf("comment('%s')", col.Comment))
	}

	names := make([]string, 0, len(col.Indexes))
	for name := range col.Indexes {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		index := table.Indexes[name]
		var uistr string
		if index.Type == schemas.UniqueType {
			uistr = "unique"
		} else if index.Type == schemas.IndexType {
			uistr = "index"
		}
		if len(index.Cols) > 1 {
			uistr += "(" + index.Name + ")"
		}
		res = append(res, uistr)
	}

	nstr := col.SQLType.Name
	if col.Length != 0 {
		if col.Length2 != 0 {
			nstr += fmt.Sprintf("(%v,%v)", col.Length, col.Length2)
		} else {
			nstr += fmt.Sprintf("(%v)", col.Length)
		}
	} else if len(col.EnumOptions) > 0 { //enum
		nstr += "("
		opts := ""

		enumOptions := make([]string, 0, len(col.EnumOptions))
		for enumOption := range col.EnumOptions {
			enumOptions = append(enumOptions, enumOption)
		}
		sort.Strings(enumOptions)

		for _, v := range enumOptions {
			opts += fmt.Sprintf(",'%v'", v)
		}
		nstr += strings.TrimLeft(opts, ",")
		nstr += ")"
	} else if len(col.SetOptions) > 0 { //enum
		nstr += "("
		opts := ""

		setOptions := make([]string, 0, len(col.SetOptions))
		for setOption := range col.SetOptions {
			setOptions = append(setOptions, setOption)
		}
		sort.Strings(setOptions)

		for _, v := range setOptions {
			opts += fmt.Sprintf(",'%v'", v)
		}
		nstr += strings.TrimLeft(opts, ",")
		nstr += ")"
	}
	res = append(res, nstr)
	if len(res) > 0 {
		return fmt.Sprintf(`xorm:"%s"`, strings.Join(res, " "))
	}
	return ""
}

func include(source []string, target string) bool {
	for _, s := range source {
		if s == target {
			return true
		}
	}
	return false
}
