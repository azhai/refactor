[![](https://goreportcard.com/badge/gitea.com/azhai/refactor)](https://goreportcard.com/report/gitea.com/azhai/refactor)

# Reverse

一个灵活高效的数据库反转工具。

## 安装

```
go get gitea.com/azhai/refactor
```

## 使用

```
reverse -f settings.yml
```

## 配置文件

一个最简单的配置文件看起来如下：

```yml
application:
   debug: true
   plural_table: false  #表名是否使用复数

connections:
   another:
      driver_name: "sqlite"
      params:
         database: ../testdata/test.db
   cache:
      driver_name: "redis"
      params:
         host: "127.0.0.1"
         port: 6379
         password: ""
         database: "0"
   default:
      driver_name: "mysql"
      table_prefix: "t_"
      read_only: false
      params:
         host: "127.0.0.1"
         port: 3306
         username: "root"
         password: ""
         database: "test"
         options: { charset: "utf8mb4" }

reverse_targets:
   -  type: codes
      include_tables: # 包含的表，以下可以用 **
         - a
         - b
      exclude_tables: # 排除的表，以下可以用 **
         - c
      table_mapper: snake # 表名到代码类或结构体的映射关系
      column_mapper: snake # 字段名到代码或结构体成员的映射关系
      table_prefix: "" # 表前缀
      output_dir: ./models # 代码生成目录
      multiple_files: false # 是否生成多个文件
      gen_json_tag: true
      gen_table_name: true
      template_path: ./template/goxorm.tmpl # 生成的模板的路径，优先级比 template 低，但比 language 中的默认模板高
      template: | # 生成模板，如果这里定义了，优先级比 template_path 高
        package {{.Target.NameSpace}}

        {{$ilen := len .Imports}}{{if gt $ilen 0}}import (
            {{range .Imports}}"{{.}}"{{end}}
        ){{end}}
        {{$gen_json := .Target.GenJsonTag -}}
        {{$gen_table := .Target.GenTableName -}}

        {{range $table_name, $table := .Tables}}
        {{$class := TableMapper $table.Name}}
        type {{$class}} struct { {{- range $table.ColumnsSeq}}{{$col := $table.GetColumn .}}
            {{ColumnMapper $col.Name}} {{Type $col}} `{{Tag $table $col $gen_json}}`{{end}}
        }

        {{if $gen_table -}}
        func ({{$class}}) TableName() string {
            return "{{$table_name}}"
        }
        {{end -}}
        {{end -}}
```
