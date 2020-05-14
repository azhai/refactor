package common

import (
	"fmt"
	"reflect"

	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
)

// 对参数先进行转义Quote
func Qprintf(engine *xorm.Engine, format string, args ...interface{}) string {
	if engine != nil {
		for i, arg := range args {
			args[i] = engine.Quote(arg.(string))
		}
	}
	return fmt.Sprintf(format, args...)
}

// 获取Model的字段列表
func GetPrimarykey(engine *xorm.Engine, m interface{}) *schemas.Column {
	table, err := engine.TableInfo(m)
	if err != nil {
		return nil
	}
	if cols := table.PKColumns(); len(cols) > 0 {
		return cols[0]
	}
	return nil
}

// 获取Model的字段列表
func GetColumns(m ITableName, alias string) (cols []string) {
	var st reflect.Type
	v := reflect.ValueOf(m)
	if v.Kind() == reflect.Ptr {
		st = reflect.Indirect(v).Type()
	} else {
		st = reflect.TypeOf(m)
	}

	if alias == "" {
		alias = m.TableName()
	}
	for i := 0; i < st.NumField(); i++ {
		t := st.Field(i).Tag.Get("json")
		if t == "" || t == "-" {
			continue
		}
		cols = append(cols, fmt.Sprintf("%s.%s", alias, t))
	}
	return cols
}

// 调整从后往前翻页
func NegativeOffset(offset, pagesize, total int) int {
	if remain := total % pagesize; remain > 0 {
		offset += pagesize - remain
	}
	return offset + total
}

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
