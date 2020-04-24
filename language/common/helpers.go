package common

import (
	"xorm.io/xorm"
)

/**
 * 过滤查询
 */
type FilterFunc = func(query *xorm.Session) *xorm.Session


/**
 * 翻页查询，out参数需要传引用
 * 使用方法 total, err := Paginate(query, &rows, pageno, pagesize)
 */
func Paginate(query *xorm.Session, out interface{}, pageno, pagesize int) (int, error) {
	var (
		total int64
		err   error
	)
	offset, limit := 0, -1 // 初始值
	if pageno < 0 {
		total, err = query.Count()
		if err != nil || total <= 0 {
			return -1, err
		}
		offset = int(total) + pageno*pagesize
	}
	// 参数校正
	if pagesize >= 0 {
		limit = pagesize
		offset = (pageno - 1) * pagesize
	}
	if limit >= 0 && offset >= 0 {
		query = query.Limit(limit, offset)
	}
	if total > 0 {
		err = query.Find(out)
	} else {
		total, err = query.FindAndCount(out)
	}
	return int(total), err
}
