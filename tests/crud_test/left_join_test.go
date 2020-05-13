package crud_test

import (
	"testing"

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
	pageno, pagesize := 1, 20
	query := db.Table(table)
	total, err := filter(query).Count()
	assert.NoError(t, err)
	if err == nil && total > 0 {
		query = base.LeftJoinQuery(engine, m, *m.PrinGroup, "prin_gid")
		base.Paginate(filter(query), pageno, pagesize).Find(&objs)
	}
	pp.Println(objs)
}
