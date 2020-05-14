package common

import (
	"bytes"
	"fmt"

	"gitea.com/azhai/refactor/defines/join"
	"xorm.io/xorm"
)

// 联表查询
func JoinQuery(engine *xorm.Engine, query *xorm.Session,
	table, fkey string, foreign ForeignTable) (*xorm.Session, []string) {
	frgTable, frgAlias := foreign.TableName(), foreign.AliasName()
	cond := Qprintf(engine, "%s.%s = %s.%s", table, fkey, frgAlias, foreign.Index)
	if query == nil {
		query = engine.Table(table)
	}
	var cols []string
	cols = GetColumns(foreign.Table, frgAlias, cols)
	query = query.Join(string(foreign.Join), frgTable, cond)
	return query, cols
}

// 关联表
type ForeignTable struct {
	Join  join.JoinOp
	Table ITableName
	Alias string
	Index string
}

func (f ForeignTable) AliasName() string {
	if f.Alias != "" {
		return f.Alias
	}
	return f.Table.TableName()
}

func (f ForeignTable) TableName() string {
	table := f.Table.TableName()
	if f.Alias != "" {
		return fmt.Sprintf("%s as %s", table, f.Alias)
	}
	return table
}

// Left Join 联表查询
type LeftJoinQuery struct {
	engine      *xorm.Engine
	filters     []FilterFunc
	nativeTable string
	Native      ITableName
	ForeignKeys []string
	Foreigns    map[string]ForeignTable
	*xorm.Session
}

func NewLeftJoinQuery(engine *xorm.Engine, native ITableName) *LeftJoinQuery {
	nativeTable := native.TableName()
	return &LeftJoinQuery{
		engine:      engine,
		filters:     nil,
		nativeTable: nativeTable,
		Native:      native,
		Foreigns:    make(map[string]ForeignTable),
		Session:     engine.Table(nativeTable),
	}
}

func (q LeftJoinQuery) Quote(value string) string {
	return q.engine.Quote(value)
}

func (q *LeftJoinQuery) ClearFilters() *LeftJoinQuery {
	q.filters = make([]FilterFunc, 0)
	return q
}

func (q *LeftJoinQuery) AddFilter(filter FilterFunc) *LeftJoinQuery {
	q.filters = append(q.filters, filter)
	return q
}

func (q *LeftJoinQuery) LeftJoin(foreign ITableName, fkey string) *LeftJoinQuery {
	q.AddLeftJoin(foreign, "", fkey, "")
	return q
}

// 添加次序要和 struct 定义一致
func (q *LeftJoinQuery) AddLeftJoin(foreign ITableName, pkey, fkey, alias string) *LeftJoinQuery {
	if pkey == "" {
		col := GetPrimarykey(q.engine, foreign)
		if col != nil {
			pkey = col.Name
		}
	}
	if _, ok := q.Foreigns[fkey]; !ok {
		q.ForeignKeys = append(q.ForeignKeys, fkey)
	}
	q.Foreigns[fkey] = ForeignTable{
		Join:  join.LEFT_JOIN,
		Table: foreign,
		Alias: alias,
		Index: pkey,
	}
	return q
}

func (q *LeftJoinQuery) Limit(limit int, start ...int) *LeftJoinQuery {
	q.AddFilter(func(query *xorm.Session) *xorm.Session {
		return query.Limit(limit, start...)
	})
	return q
}

func (q *LeftJoinQuery) OrderBy(order string) *LeftJoinQuery {
	q.AddFilter(func(query *xorm.Session) *xorm.Session {
		return query.OrderBy(order)
	})
	return q
}

func (q *LeftJoinQuery) GetQuery() *xorm.Session {
	buf := new(bytes.Buffer)
	buf.WriteString(Qprintf(q.engine, "%s.*", q.Native.TableName()))
	query := q.Session.Clone()
	for _, filter := range q.filters {
		query = filter(query)
	}
	var cols []string
	for _, fkey := range q.ForeignKeys {
		foreign := q.Foreigns[fkey]
		query, cols = JoinQuery(q.engine, query, q.nativeTable, fkey, foreign)
		buf.WriteString(", ")
		buf.WriteString(BlindlyQuote(q.engine, ", ", cols...))
	}
	return query.Select(buf.String())
}

func (q *LeftJoinQuery) Paginate(pageno, pagesize int) *xorm.Session {
	return Paginate(q.GetQuery(), pageno, pagesize)
}

func (q *LeftJoinQuery) Count(bean ...interface{}) (int64, error) {
	query := q.Session.Clone()
	for _, filter := range q.filters {
		query = filter(query)
	}
	return query.Count(bean...)
}

func (q *LeftJoinQuery) FindAndCount(rowsSlicePtr interface{}, condiBean ...interface{}) (int64, error) {
	total, err := q.Count()
	if err != nil || total == 0 {
		return total, err
	}
	err = q.GetQuery().Find(rowsSlicePtr, condiBean...)
	return total, err
}
