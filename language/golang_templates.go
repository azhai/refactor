package language

import (
	"fmt"
	"html/template"
	"strings"
)

var (
	golangModelTemplate = fmt.Sprintf(`package {{.Target.NameSpace}}

{{$ilen := len .Imports}}{{if gt $ilen 0 -}}
import (
	{{range $imp, $al := .Imports}}{{$al}} "{{$imp}}"{{end}}
)
{{end -}}
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
{{end}}
`, "`", "`")

	golangCacheTemplate = `package {{.Target.NameSpace}}

import (
	"gitea.com/azhai/refactor/builtin/base"
	"gitea.com/azhai/refactor/config"
	"github.com/azhai/gozzo-utils/redisw"
	"github.com/gomodule/redigo/redis"
)

var (
	sessreg *base.SessionRegistry
)

// 初始化、连接数据库和缓存
func Initialize(r *config.ReverseSource, verbose bool) {
	d := config.ReverseSource2RedisDialect(r)
	dial := func() (redis.Conn, error) {
		return d.Connect(verbose)
	}
	sessreg = base.NewRegistry(redisw.NewRedisPool(dial, -1))
}

// 获得当前会话管理器
func Registry() *base.SessionRegistry {
	return sessreg
}

// 获得用户会话
func Session(token string) *base.Session {
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

	golangConnTemplate = `package {{.Target.NameSpace}}

import (
	"gitea.com/azhai/refactor/builtin/base"
	"gitea.com/azhai/refactor/config"
	_ "{{.ImporterPath}}"
	"xorm.io/xorm"
)

var (
	engine  *xorm.Engine
)

// 初始化、连接数据库和缓存
func Initialize(r *config.ReverseSource, verbose bool) {
	var err error
	engine, err = r.Connect(verbose)
	if err != nil {
		panic(err)
	}
}

// 查询某张数据表
func Engine() *xorm.Engine {
	return engine
}

// 转义表名或字段名
func Quote(value string) string {
	if engine == nil {
		return value
	}
	return engine.Quote(value)
}

// 查询某张数据表
func Table(args ...interface{}) *xorm.Session {
	if engine == nil {
		return nil
	}
	if args == nil {
		return engine.NewSession()
	}
	return engine.Table(args[0])
}

// 执行事务
func ExecTx(modify base.ModifyFunc) error {
	tx := engine.NewSession() // 必须是新的session
	defer tx.Close()
	_ = tx.Begin()
	if _, err := modify(tx); err != nil {
		_ = tx.Rollback() // 失败回滚
		return err
	}
	return tx.Commit()
}

// 查询多行数据
func QueryAll(filter base.FilterFunc, pages ...int) *xorm.Session {
	query := engine.NewSession()
	if filter != nil {
		query = filter(query)
	}
	pageno, pagesize := 0, -1
	if len(pages) >= 1 {
		pageno = pages[0]
		if len(pages) >= 2 {
			pagesize = pages[1]
		}
	}
	return base.Paginate(query, pageno, pagesize)
}
`

	golangInitTemplate = `package models

{{$initns := .Target.InitNameSpace -}}
import (
	"gitea.com/azhai/refactor/config"
	"gitea.com/azhai/refactor/config/dialect"
	{{- range $dir, $al := .Imports}}
	{{if ne $al $dir}}{{$al}} {{end -}}
	"{{$initns}}/{{$dir}}"{{end}}
)

var (
	configFile = "./settings.yml"
	settings   *config.Settings
	verbose    bool
)

func init() {
	settings, err := config.ReadSettings(configFile)
	if err != nil {
		panic(err)
	}
	verbose = settings.Application.Debug
	ConnectDatabases(settings.GetConnConfigMap())
}

func Verbose() bool {
	return verbose
}

func ConnectDatabases(confs map[string]config.ConnConfig) {
	for key, c := range confs {
		r, _ := config.NewReverseSource(c)
		switch key {
		{{- range $dir, $al := .Imports}}
		case "{{$dir}}":
			{{$al}}.Initialize(r, verbose){{end}}
		}
	}
}
`

	golangQueryTemplate = `{{if not .Target.MultipleFiles}}package {{.Target.NameSpace}}

import (
	"time"

	{{range $imp, $al := .Imports}}{{$al}} "{{$imp}}"{{end}}
)
{{end -}}

{{range .Tables}}
{{$class := TableMapper .Name -}}
{{$pkey := GetSinglePKey . -}}
{{$created := GetCreatedColumn . -}}
func (m *{{$class}}) Load(where interface{}, args ...interface{}) (bool, error) {
	return Table().Where(where, args...).Get(m)
}

{{if ne $pkey "" -}}
func (m *{{$class}}) Save(changes map[string]interface{}) error {
	return ExecTx(func(tx *xorm.Session) (int64, error) {
		if changes == nil || m.{{$pkey}} == 0 {
			{{if ne $created "" -}}changes["{{$created}}"] = time.Now()
			{{else}}{{end -}}
			return tx.Table(m).Insert(changes)
		} else {
			return tx.Table(m).ID(m.{{$pkey}}).Update(changes)
		}
	})
}
{{end -}}
{{end -}}
`
)

func GetGolangTemplate(name string, funcs template.FuncMap) *template.Template {
	var content string
	switch strings.ToLower(name) {
	case "cache":
		name, content = "cache", golangCacheTemplate
	case "conn":
		name, content = "conn", golangConnTemplate
	case "init":
		name, content = "init", golangInitTemplate
	case "query":
		name, content = "query", golangQueryTemplate
	default:
		name, content = "model", golangModelTemplate
	}
	if tmpl := GetPresetTemplate(name); tmpl != nil {
		return tmpl
	}
	return NewTemplate(name, content, funcs)
}
