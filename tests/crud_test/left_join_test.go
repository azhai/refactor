package crud_test

import (
	"testing"

	"github.com/azhai/refactor/builtin/base"
	"github.com/azhai/refactor/builtin/join"
	"github.com/azhai/refactor/inspect"
	"github.com/azhai/refactor/tests/contrib"
	_ "github.com/azhai/refactor/tests/models"
	db "github.com/azhai/refactor/tests/models/default"
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
	"xorm.io/xorm"
)

func getFilter(engine *xorm.Engine, table string) base.FilterFunc {
	return func(query *xorm.Session) *xorm.Session {
		cond := base.Qprintf(engine, "%s.%s IS NOT NULL", table, "prin_gid")
		sort := base.Qprintf(engine, "%s.%s ASC", table, "id")
		return query.Where(cond).OrderBy(sort)
	}
}

func TestJoin01FindUserGroups(t *testing.T) {
	m := &contrib.UserWithGroup{
		PrinGroup: new(contrib.GroupSummary),
	}
	engine, table := db.Engine(), m.TableName()
	filter := getFilter(engine, table)

	query := db.Table(table)
	total, err := filter(query).Count()
	assert.NoError(t, err)
	var objs []*contrib.UserWithGroup
	if err == nil && total > 0 {
		var cols []string
		cols = inspect.GetColumns(m.User, m.User.TableName(), cols)
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
	filter := getFilter(engine, table)

	group := contrib.GroupSummary{}
	query := base.NewLeftJoinQuery(engine, native)
	query.AddLeftJoin(group, "gid", "prin_gid", "P")
	query.AddLeftJoin(group, "gid", "vice_gid", "V")

	var objs []*contrib.UserWithGroup
	pageno, pagesize := 1, 20
	query = query.AddFilter(filter)
	_, err := query.FindPaginate(pageno, pagesize, &objs)
	assert.NoError(t, err)
	if testing.Verbose() {
		pp.Println(objs)
	}
}
