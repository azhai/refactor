[![](https://goreportcard.com/badge/gitea.com/azhai/refactor)](https://goreportcard.com/report/gitea.com/azhai/refactor)

# Xorm-Refactor

一个灵活高效的数据库反转工具。将数据库中的表生成对应的Model和
相应的查询方法，特点是自动套用包名和将已知的 Mixin 嵌入Model 中。

## 常见用法

以MySQL中的一个菜单表为例，建表SQL语句:

```sql
CREATE TABLE `t_menu`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `lft` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '左边界',
  `rgt` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '右边界',
  `depth` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '高度',
  `path` varchar(100) NOT NULL DEFAULT '' COMMENT '路径',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',
  `icon` varchar(30) NULL DEFAULT NULL COMMENT '图标',
  `remark` text NULL COMMENT '说明备注',
  `created_at` timestamp(0) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp(0) NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` timestamp(0) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_t_menu_rgt`(`rgt`) USING BTREE,
  INDEX `idx_t_menu_depth`(`depth`) USING BTREE,
  INDEX `idx_t_menu_path`(`path`) USING BTREE,
  INDEX `idx_t_menu_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COMMENT = '菜单' ROW_FORMAT = Compact;
```

默认在 models/default/ 下生成 3 个代码文件 init.go 、 models.go 和 queries.go

```go
package db
// Filename is models.go

import (
	base "gitea.com/azhai/refactor/language/common"
)

type Menu struct {
	Id                int `json:"id" xorm:"notnull pk autoincr INT(10)"`
	*base.NestedMixin `xorm:"extends"`
	Path              string `json:"path" xorm:"notnull default '' comment('路径') index VARCHAR(100)"`
	Title             string `json:"title" xorm:"notnull default '' comment('名称') VARCHAR(50)"`
	Icon              string `json:"icon" xorm:"comment('图标') VARCHAR(30)"`
	Remark            string `json:"remark" xorm:"comment('说明备注') TEXT"`
	*base.TimeMixin   `xorm:"extends"`
}
```

```go
package db
// Filename is queries.go

import (
	 "xorm.io/xorm"
)

func (m *Menu) Load(where interface{}, args ...interface{}) (bool, error) {
	return Table().Where(where, args...).Get(m)
}

func (m *Menu) Save(changes map[string]interface{}) error {
	return ExecTx(func(tx *xorm.Session) (int64, error) {
		if changes == nil || m.Id == 0 {
			return tx.Insert(m)
		} else {
			return tx.Table(m).ID(m.Id).Update(changes)
		}
	})
}
```

init.go 中含有 Initialize() 方法，通过下面的方法，在程序入口 main() 或 init() 中
初始化数据库连接。

```go
package main

import (
	"gitea.com/azhai/refactor/config"
	"my-project/models/db"
)

var verbose bool // 详细输出

func init() {
	cfg, err := config.ReadSettings("settings.yml")
	if err != nil {
		panic(err)
	}
	db.Initialize(cfg, verbose)
}
```

init.go 中含有 QueryAll() 方法用于查询表中多行数据，而 queries.go 中的  Load() 方法
只查询指定的一行数据。配合 NestedMixin 我们可以查询子菜单：

```go
package main

import (
	"my-project/models/db"
)

// 根据 id 找出菜单及其子菜单（最多三层）
func GetMenus(id int) ([]*db.Menu, error) {
	m := new(db.Menu)
	has, err := m.Load("id = ?", id) // 找出指定 id 的菜单
	if err != nil || has == false {
		return nil, err
	}
	filter := m.ChildrenFilter(3) // 最多三层子菜单，如不限制传递参数 0
	pageno, pagesize := 0, -1     // 符合条件所有数据，即不分页
	var menus []*db.Menu
	err = db.QueryAll(filter, pageno, pagesize).Find(&menus)
	return menus, err
}
```



## 安装

```
go get gitea.com/azhai/refactor
```

## 编译使用

```
make all
./reverse -f tests/settings.yml
```

## 配置文件

一个典型的配置文件看起来如下：

```yml
application:
   debug: true
   plural_table: false  #表名是否使用复数

mysql: &mysql           #共用数据库配置
   driver_name: "mysql"
   params:
      host: "127.0.0.1"
      port: 3306
      username: "root"
      password: ""
      database: "test"
      options: { charset: "utf8" }

connections:
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
      include_tables:    # 包含的表，以下可以用
         - "a*"
         - "b*"
      exclude_tables:    # 排除的表，以下可以用
         - "c"
      <<: *mysql         #引用mysql配置

reverse_targets:
   -  type: codes
      table_mapper: snake     # 表名到代码类或结构体的映射关系
      column_mapper: snake    # 字段名到代码或结构体成员的映射关系
      output_dir: "./models"  # 代码生成目录
      multiple_files: false   # 是否生成多个文件
      template_path: ""       # 生成的模板的路径，优先级比 language 中的默认模板高
      query_template_path: "./tests/golang_query.tmpl" # 自定义查询方法模板
      gen_json_tag: true      # 生成JSON标签
      gen_table_name: true    # 生成TableName()方法
      gen_query_methods: true # 生成查询方法
      apply_mixins: true      # 使用已知的Mixin替换部分字段
      mixin_dir_path: ""      # 额外的mixin目录
      mixin_name_space: ""    # 额外的mixin包名
```
