package common

import (
	"xorm.io/xorm"
)

/**
 * 过滤查询
 */
type FilterFunc = func(query *xorm.Session) *xorm.Session

func Where(query *xorm.Session, conds []interface{}) *xorm.Session {
	if len(conds) >= 1 {
		query = query.And(conds[0], conds[1:]...)
	}
	return query
}

func OrWhere(query *xorm.Session, conds []interface{}) *xorm.Session {
	if len(conds) >= 1 {
		query = query.Or(conds[0], conds[1:]...)
	}
	return query
}

func Paginate(query *xorm.Session, pageno, pagesize int) *xorm.Session {
	if pagesize < 0 {
		return query
	} else if pagesize == 0 {
		return query.ID(0)
	}
	var offset int
	if pageno > 0 {
		offset = (pageno - 1) * pagesize
	} else if pageno < 0 {
		total, err := query.Count()
		if err == nil && total > 0 {
			offset = pageno * pagesize + int(total)
		}
	}
	return query.Limit(pagesize, offset)
}
