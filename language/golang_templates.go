package language

import (
	"fmt"
	"html/template"
	"strings"
)

var (
	initTemplates map[string]*template.Template

	golangModelTemplate = fmt.Sprintf(`package {{.Target.NameSpace}}

{{$ilen := len .Imports}}{{if gt $ilen 0}}import (
	{{range .Imports}}"{{.}}"{{end}}
){{end}}
{{$gen_json := .Target.GenJsonTag -}}
{{$gen_table := .Target.GenTableName -}}

{{range $table_name, $table := .Tables}}
{{$class := TableMapper $table.Name}}
type {{$class}} struct { {{- range $table.ColumnsSeq}}{{$col := $table.GetColumn .}}
	{{ColumnMapper $col.Name}} {{Type $col}} %s{{Tag $table $col $gen_json}}%s{{end}}
}

{{if $gen_table -}}
func ({{$class}}) TableName() string {
	return "{{$table_name}}"
}
{{end -}}
{{end -}}
`, "`", "`")

	golangConnTemplate = `package {{.Target.NameSpace}}

import (
	"gitea.com/azhai/refactor/config"
	"gitea.com/azhai/refactor/language/common"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var (
	engine  *xorm.Engine
)

// 初始化、连接数据库和缓存
func Initialize(cfg *config.Settings, verbose bool) {
	var err error
	engine, err = common.InitConn(cfg, "{{.ConnKey}}", verbose)
	if err != nil || engine == nil {
		panic(err)
	}
}

// 查询某张数据表
func Engine() *xorm.Engine {
	return engine
}

// 查询某张数据表
func Table(name interface{}) *xorm.Session {
	if engine == nil {
		return nil
	}
	return engine.Table(name)
}
`

	golangCacheTemplate = `package {{.Target.NameSpace}}

import (
	"gitea.com/azhai/refactor/config"
	"gitea.com/azhai/refactor/language/common"
	"xorm.io/xorm"
)

var (
	sessreg *common.SessionRegistry
)

// 初始化、连接数据库和缓存
func Initialize(cfg *config.Settings, verbose bool) {
	var err error
	sessreg, err = common.InitCache(cfg, "{{.ConnKey}}", verbose)
	if err != nil {
		panic(err)
	}
}

// 获得当前会话管理器
func Registry() *common.SessionRegistry {
	return sessreg
}

// 获得用户会话
func Session(token string) *common.Session {
	if sessreg == nil {
		return nil
	}
	return sessreg.GetSession(token)
}

// 删除会话
func DelSession(token string) bool {
	if sessreg == nil {
		return false
	}
	return sessreg.DelSession(token)
}
`

	golangQueryTemplate = `package {{.Target.NameSpace}}

import (
	base "gitea.com/azhai/refactor/language/common"
)

{{$class := TableMapper $table.Name -}}
{{$classes := DiffPluralize $class "Rows" -}}
type {{$class}} []{{$class}}

func (m *{{$class}}) One(where interface{}, args ...interface{}) (has bool, err error) {
	query := engine.NewSession().Where(where, args)
	return query.Get(&m)
}

func (m *{{$class}}) FindOne(filters ...base.FilterFunc) (has bool, err error) {
	query := engine.NewSession()
	for _, ft := range filters {
		query = ft(query)
	}
	return query.Get(&m)
}

func (s *{{$classes}}) All(pageno, pagesize int, where interface{}, args ...interface{}) (total int64, err error) {
	query := engine.NewSession().Where(where, args)
	return base.Paginate(query, pageno, pagesize).FindAndCount(s)
}

func (s *{{$classes}}) FindAll(pageno, pagesize int, filters ...base.FilterFunc) (total int64, err error) {
	query := engine.NewSession()
	for _, ft := range filters {
		query = ft(query)
	}
	return base.Paginate(query, pageno, pagesize).FindAndCount(s)
}
`
)

func GetGolangTemplate(name string) *template.Template {
	var content string
	switch strings.ToLower(name) {
	default:
		name, content = "model", golangModelTemplate
	case "cache":
		name, content = "cache", golangCacheTemplate
	case "conn":
		name, content = "conn", golangConnTemplate
	case "query":
		name, content = "query", golangQueryTemplate
	}
	if tmpl, ok := initTemplates[name]; ok {
		return tmpl
	}
	t := template.New(name)
	if tmpl, err := t.Parse(content); err == nil {
		return tmpl
	}
	return nil
}