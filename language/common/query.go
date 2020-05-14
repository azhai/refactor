package common

import (
	"fmt"
	"reflect"
	"strings"

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

func GetIndirectType(v interface{}) (rt reflect.Type) {
	var ok bool
	if rt, ok = v.(reflect.Type); !ok {
		rt = reflect.TypeOf(v)
	}
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	return
}

func GetFinalType(v interface{}) (rt reflect.Type) {
	rt = GetIndirectType(v)
	for {
		switch rt.Kind() {
		default:
			break
		case reflect.Ptr, reflect.Chan:
			rt = rt.Elem()
		case reflect.Array, reflect.Slice:
			rt = rt.Elem()
		case reflect.Map:
			kk := rt.Key().Kind()
			if kk == reflect.String || kk <= reflect.Float64 {
				rt = rt.Elem()
			} else {
				break
			}
		}
	}
	return
}

func GetColumns(v interface{}, alias string, cols []string) []string {
	rt := GetIndirectType(v)
	if rt.Kind() != reflect.Struct {
		return cols
	}
	for i := 0; i < rt.NumField(); i++ {
		t := rt.Field(i).Tag.Get("json")
		if t == "" || t == "-" {
			continue
		} else if strings.HasSuffix(t, "inline") {
			cols = GetColumns(rt.Field(i).Type, alias, cols)
		} else {
			if alias != "" {
				t = fmt.Sprintf("%s.%s", alias, t)
			}
			cols = append(cols, t)
		}
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
