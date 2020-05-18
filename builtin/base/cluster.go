package base

import (
	"time"

	"xorm.io/xorm"
)

/**
 * 创建或清空表
 */
type IResetTable interface {
	ResetTable(curr string, trunc bool) error
	ITableName
}

/**
 * 写入多行到分布式表中
 */
func ClusterInsertMulti(xi xorm.Interface, objs []IResetTable, reset, trunc bool) (total int64, err error) {
	var count int64
	last, max := 0, len(objs)-1
	for i, m := range objs {
		table := m.TableName()
		if i < max && objs[i+1].TableName() == table {
			continue
		}
		if reset {
			if err = m.ResetTable(table, trunc); err != nil {
				return
			}
		}
		count, err = xi.Table(table).InsertMulti(objs[last : i+1])
		last, total = i+1, total+count
		if err != nil {
			return
		}
	}
	return
}

// 按月分表
type MonthlyMixin struct {
	Suffixes []string
	Date     time.Time
	Format   string
}

func NewMonthlyMixin(t time.Time) *MonthlyMixin {
	m := &MonthlyMixin{Format: "200601"}
	return m.SetTime(t)
}

func (m MonthlyMixin) GetSuffix() string {
	return m.Date.Format(m.Format)
}

func (m *MonthlyMixin) SetTime(t time.Time) *MonthlyMixin {
	if t.IsZero() {
		m.Date = t
		return m
	}
	year, month, _ := t.Date()
	loc := time.Local
	m.Date = time.Date(year, month, 1, 0, 0, 0, 0, loc)
	return m
}

func (m *MonthlyMixin) Move(n int) *MonthlyMixin {
	m.Date = m.Date.AddDate(0, n, 0)
	return m
}

func (m *MonthlyMixin) Prev() *MonthlyMixin {
	return m.Move(-1)
}

func (m *MonthlyMixin) Next() *MonthlyMixin {
	return m.Move(1)
}
