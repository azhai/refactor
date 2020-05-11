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

	golangConnTemplate = `package {{.Target.NameSpace}}

import (
	"gitea.com/azhai/refactor/config"
	base "gitea.com/azhai/refactor/language/common"
	_ "{{.ImporterPath}}"
	"xorm.io/xorm"
)

var (
	engine  *xorm.Engine
)

// 初始化、连接数据库和缓存
func Initialize(c config.ConnConfig, verbose bool) {
	var err error
	engine, err = c.Connect(verbose)
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

	golangCacheTemplate = `package {{.Target.NameSpace}}

import (
	"gitea.com/azhai/refactor/config"
	base "gitea.com/azhai/refactor/language/common"
)

var (
	sessreg *base.SessionRegistry
)

// 初始化、连接数据库和缓存
func Initialize(c config.ConnConfig, verbose bool) {
	var err error
	sessreg, err = base.InitCache(c, verbose)
	if err != nil {
		panic(err)
	}
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
)

func GetGolangTemplate(name string, funcs template.FuncMap) *template.Template {
	var content string
	switch strings.ToLower(name) {
	case "cache":
		name, content = "cache", golangCacheTemplate
	case "conn":
		name, content = "conn", golangConnTemplate
	default:
		name, content = "model", golangModelTemplate
	}
	if tmpl := GetPresetTemplate(name); tmpl != nil {
		return tmpl
	}
	return NewTemplate(name, content, funcs)
}
