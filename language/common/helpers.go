package common

import (
	"xorm.io/xorm"
)

/**
 * 过滤查询
 */
type FilterFunc = func(query *xorm.Session) *xorm.Session

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
			offset = NegativeOffset(pageno * pagesize, pagesize, int(total))
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
