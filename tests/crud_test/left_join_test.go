package crud_test

import (
	"testing"

	"github.com/k0kubun/pp"

	"gitea.com/azhai/refactor/defines/join"
	base "gitea.com/azhai/refactor/language/common"
	"gitea.com/azhai/refactor/tests/contrib"
	_ "gitea.com/azhai/refactor/tests/models"
	db "gitea.com/azhai/refactor/tests/models/default"
	"github.com/stretchr/testify/assert"
	"xorm.io/xorm"
)

func TestJoin01FindUserGroups(t *testing.T) {
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
		var cols []string
		cols = base.GetColumns(m.User, m.User.TableName(), cols)
		query = filter(query).Cols(cols...)
		if testing.Verbose() {
			pp.Println(cols)
		}
		foreign := base.ForeignTable{
			Join:  join.LEFT_JOIN,
			Table: *m.PrinGroup,
			Alias: "", Index: "gid",
		}
		foreign.Alias = "P"
		query, cols = base.JoinQuery(engine, query, table, "prin_gid", foreign)
		query = query.Cols(cols...)
		if testing.Verbose() {
			pp.Println(cols)
		}
		foreign.Alias = "V"
		query, cols = base.JoinQuery(engine, query, table, "vice_gid", foreign)
		query = query.Cols(cols...)
		if testing.Verbose() {
			pp.Println(cols)
		}
		pageno, pagesize := 1, 20
		base.Paginate(query, pageno, pagesize).Find(&objs)
		if testing.Verbose() {
			pp.Println(objs)
		}
	}
}

func TestJoin02LeftJoinQuery(t *testing.T) {
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
	query.AddLeftJoin(group, "gid", "prin_gid", "P")
	query.AddLeftJoin(group, "gid", "vice_gid", "V")

	query = query.AddFilter(filter).Limit(20)
	_, err := query.FindAndCount(&objs)
	if testing.Verbose() {
		pp.Println(objs)
	}
	assert.NoError(t, err)
}