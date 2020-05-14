package crud_test

import (
	"testing"

	"gitea.com/azhai/refactor/defines/join"

	base "gitea.com/azhai/refactor/language/common"
	"gitea.com/azhai/refactor/tests/contrib"
	_ "gitea.com/azhai/refactor/tests/models"
	db "gitea.com/azhai/refactor/tests/models/default"
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"

	"xorm.io/xorm"
)

func Test01FindUserGroups(t *testing.T) {
	m := &contrib.UserWithGroup{
		PrinGroup: new(contrib.GroupSummary),
	}
	engine, table := db.Engine(), m.TableName()
	filter := func(query *xorm.Session) *xorm.Session {
		cond := base.Qprintf(engine, "%s.%s IS NOT NULL", table, "prin_gid")
		sort := base.Qprintf(engine, "%s.%s ASC", table, "id")
		return query.Where(cond).OrderBy(sort)
	}

	var objs []*contrib.UserWithGroup
	query := db.Table(table)
	total, err := filter(query).Count()
	assert.NoError(t, err)
	if err == nil && total > 0 {
		cols := base.GetColumns(m.User, "")
		query = filter(query).Cols(cols...)
		foreign := base.ForeignTable{
			Join:  join.LEFT_JOIN,
			Table: *m.PrinGroup,
			Alias: "", Index: "gid",
		}
		foreign.Alias = "P"
		query = base.JoinQuery(engine, query, table, "prin_gid", foreign)
		foreign.Alias = "V"
		query = base.JoinQuery(engine, query, table, "vice_gid", foreign)
		pageno, pagesize := 1, 20
		base.Paginate(query, pageno, pagesize).Find(&objs)
	}
	pp.Println(objs)
}

func Test02LeftJoinQuery(t *testing.T) {
	engine, native := db.Engine(), db.User{}
	table := native.TableName()
	filter := func(query *xorm.Session) *xorm.Session {
		cond := base.Qprintf(engine, "%s.%s IS NOT NULL", table, "prin_gid")
		sort := base.Qprintf(engine, "%s.%s ASC", table, "id")
		return query.Where(cond).OrderBy(sort)
	}

	var objs []*contrib.UserWithGroup
	group := contrib.GroupSummary{}
	query := base.NewLeftJoinQuery(engine, native)
	query = query.SetFilter(filter)
	query.AddLeftJoin(group, "gid", "prin_gid", "P")
	query.AddLeftJoin(group, "gid", "vice_gid", "V")

	/*
		_, err := query.FindAndCount(&objs)
		assert.NoError(t, err)
	*/
	total, err := query.Count()
	assert.NoError(t, err)
	if err == nil && total > 0 {
		pageno, pagesize := 1, 20
		query.Paginate(pageno, pagesize).Find(&objs)
	}
	pp.Println(objs)
}
