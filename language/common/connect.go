package common

import (
	"time"

	"gitea.com/azhai/refactor/config"
	"gitea.com/azhai/refactor/config/dialect"
	"github.com/k0kubun/pp"
	"xorm.io/xorm"
)

func InitConn(cfg *config.Settings, name string, verbose bool) (*xorm.Engine, error) {
	var drv, dsn string
	if c, ok := cfg.Connections[name]; ok {
		d := dialect.GetDialectByName(c.DriverName)
		if d != nil {
			drv, dsn = d.Name(), d.GetDSN(c.Params)
		}
	}
	if drv == "" || dsn == "" {
		return nil, nil
	} else if verbose {
		pp.Println(drv, dsn)
	}
	engine, err := xorm.NewEngine(drv, dsn)
	if err == nil {
		engine.ShowSQL(verbose)
	}
	return engine, err
}

/**
 * 数据表名
 */
type ITableName interface {
	TableName() string
}

/**
 * 数据表注释
 */
type ITableComment interface {
	TableComment() string
}

/**
 * 带自增主键的基础Model
 */
type BaseModel struct {
	Id uint `json:"id" xorm:"not null pk autoincr INT(10)"`
}

func (BaseModel) TableComment() string {
	return ""
}

/**
 * 时间相关的三个典型字段
 */
type TimeModel struct {
	CreatedAt time.Time `json:"created_at" xorm:"created comment('创建时间') TIMESTAMP"`       // 创建时间
	UpdatedAt time.Time `json:"updated_at" xorm:"updated comment('更新时间') TIMESTAMP"`       // 更新时间
	DeletedAt time.Time `json:"deleted_at" xorm:"deleted comment('删除时间') index TIMESTAMP"` // 删除时间
}
