package common

import (
	"time"

	"xorm.io/xorm"
)

// 过滤查询
type FilterFunc = func(query *xorm.Session) *xorm.Session

// 修改操作，用于事务
type ModifyFunc = func(tx *xorm.Session) (int64, error)

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
 * 时间相关的三个典型字段
 */
type TimeMixin struct {
	CreatedAt time.Time `json:"created_at" xorm:"created comment('创建时间') TIMESTAMP"`       // 创建时间
	UpdatedAt time.Time `json:"updated_at" xorm:"updated comment('更新时间') TIMESTAMP"`       // 更新时间
	DeletedAt time.Time `json:"deleted_at" xorm:"deleted comment('删除时间') index TIMESTAMP"` // 删除时间
}
