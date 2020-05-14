package common

import (
	"fmt"

	"gitea.com/azhai/refactor/defines/join"
	"xorm.io/xorm"
)

// 联表查询
func JoinQuery(engine *xorm.Engine, query *xorm.Session, table, fkey string, foreign ForeignTable) *xorm.Session {
	frgTable, frgAlias := foreign.TableName(), foreign.AliasName()
	cond := Qprintf(engine, "%s.%s = %s.%s", table, fkey, frgAlias, foreign.Index)
	if query == nil {
		query = engine.Table(table)
	}
	query = query.Join(string(foreign.Join), frgTable, cond)
	return query.Cols(GetColumns(foreign.Table, frgAlias)...)
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
	filter      FilterFunc
	nativeTable string
	Native      ITableName
	Foreigns    map[string]ForeignTable
	*xorm.Session
}

func NewLeftJoinQuery(engine *xorm.Engine, native ITableName) *LeftJoinQuery {
	nativeTable := native.TableName()
	return &LeftJoinQuery{
		engine:      engine,
		filter:      nil,
		nativeTable: nativeTable,
		Native:      native,
		Foreigns:    make(map[string]ForeignTable),
		Session:     engine.Table(nativeTable),
	}
}

func (q LeftJoinQuery) Quote(value string) string {
	return q.engine.Quote(value)
}

func (q *LeftJoinQuery) SetFilter(filter FilterFunc) *LeftJoinQuery {
	q.filter = filter
	return q
}

func (q *LeftJoinQuery) LeftJoin(foreign ITableName, fkey string) *LeftJoinQuery {
	q.AddLeftJoin(foreign, "", fkey, "")
	return q
}

func (q *LeftJoinQuery) AddLeftJoin(foreign ITableName, pkey, fkey, alias string) {
	if pkey == "" {
		col := GetPrimarykey(q.engine, foreign)
		if col != nil {
			pkey = col.Name
		}
	}
	q.Foreigns[fkey] = ForeignTable{
		Join:  join.LEFT_JOIN,
		Table: foreign,
		Alias: alias,
		Index: pkey,
	}
}

func (q *LeftJoinQuery) GetQuery() *xorm.Session {
	cols := GetColumns(q.Native, "")
	query := q.Session.Clone().Cols(cols...)
	if q.filter != nil {
		query = q.filter(query)
	}
	for fkey, foreign := range q.Foreigns {
		query = JoinQuery(q.engine, query, q.nativeTable, fkey, foreign)
	}
	return query
}

func (q *LeftJoinQuery) Paginate(pageno, pagesize int) *xorm.Session {
	return Paginate(q.GetQuery(), pageno, pagesize)
}

func (q *LeftJoinQuery) Count(bean ...interface{}) (int64, error) {
	if q.filter != nil {
		q.Session = q.filter(q.Session)
	}
	return q.Session.Count(bean...)
}

func (q *LeftJoinQuery) FindAndCount(rowsSlicePtr interface{}, condiBean ...interface{}) (int64, error) {
	total, err := q.Count()
	if err != nil || total == 0 {
		return total, err
	}
	err = q.GetQuery().Find(rowsSlicePtr, condiBean...)
	return total, err
}
