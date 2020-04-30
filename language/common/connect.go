package common

import (
	"time"

	"gitea.com/azhai/refactor/config"
	"xorm.io/xorm"
)

func InitConn(cfg config.IConnectSettings, name string, verbose bool) (*xorm.Engine, error) {
	conns := cfg.GetConnections(name)
	if len(conns) == 0 {
		return nil, nil
	}
	ds := config.NewDataSource(name, conns[name])
	return ds.Connect(verbose)
}

// 过滤查询
type FilterFunc = func(query *xorm.Session) *xorm.Session

// 修改操作，用于事务
type ModifyFunc = func(tx *xorm.Session) (int64, error)

// 计算翻页
func Paginate(query *xorm.Session, pageno, pagesize int) *xorm.Session {
	if pagesize < 0 {
		return query
	} else if pagesize == 0 {
		return query.Limit(0)
	}
	var offset int
	if pageno > 0 {
		offset = (pageno - 1) * pagesize
	} else if pageno < 0 {
		total, err := query.Count()
		if err == nil && total > 0 {
			offset = NegativeOffset(pageno*pagesize, pagesize, int(total))
		}
	}
	return query.Limit(pagesize, offset)
}

// 调整从后往前翻页
func NegativeOffset(offset, pagesize, total int) int {
	if remain := total % pagesize; remain > 0 {
		offset += pagesize - remain
	}
	return offset + total
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
 * 时间相关的三个典型字段
 */
type TimeMixin struct {
	CreatedAt time.Time `json:"created_at" xorm:"created comment('创建时间') TIMESTAMP"`       // 创建时间
	UpdatedAt time.Time `json:"updated_at" xorm:"updated comment('更新时间') TIMESTAMP"`       // 更新时间
	DeletedAt time.Time `json:"deleted_at" xorm:"deleted comment('删除时间') index TIMESTAMP"` // 删除时间
}
