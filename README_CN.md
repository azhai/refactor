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

一个典型的配置文件看起来如下：

```yml
application:
   debug: true
   plural_table: false  #表名是否使用复数

mysql: &mysql  #共用数据库配置
   driver_name: "mysql"
   params:
      host: "127.0.0.1"
      port: 3306
      username: "root"
      password: ""
      database: "test"
      options: { charset: "utf8mb4" }

connections:
   another:
      driver_name: "sqlite"
      params:
         database: "./testdata/test.db"
   cache:
      driver_name: "redis"
      params:
         host: "127.0.0.1"
         port: 6379
         password: ""
         database: "0"
   default:
      read_only: false
      table_prefix: "t_" # 表前缀
      include_tables: # 包含的表，以下可以用
         - "a*"
         - "b*"
      exclude_tables: # 排除的表，以下可以用
         - "c"
      <<: *mysql  #引用mysql配置

reverse_targets:
   -  type: codes
      table_mapper: snake # 表名到代码类或结构体的映射关系
      column_mapper: snake # 字段名到代码或结构体成员的映射关系
      output_dir: "./models" # 代码生成目录
      multiple_files: false # 是否生成多个文件
      gen_json_tag: true   # 生成JSON标签
      gen_table_name: true # 生成TableName()方法
      gen_query_methods: true # 生成查询方法
      apply_mixins: true   # 使用已知的Mixin替换部分字段
      template_path: "./template/goxorm.tmpl" # 生成的模板的路径，优先级比 language 中的默认模板高
      query_template_path: "" # 自定义查询方法模板
```
