package crud_test

import (
	"testing"

	"gitea.com/azhai/refactor/defines/access"
	"gitea.com/azhai/refactor/tests/contrib"
	_ "gitea.com/azhai/refactor/tests/models"
	db "gitea.com/azhai/refactor/tests/models/default"
	"github.com/stretchr/testify/assert"
)

func Test01InsertAccess(t *testing.T) {
	m := new(db.Access)
	err := contrib.TruncTable(m.TableName())
	assert.NoError(t, err)

	// 超管可以访问所有菜单
	m, err = contrib.AddAccess("superuser", "menu", access.ALL, "*")
	assert.NoError(t, err)
	assert.Equal(t, m.Id, 1)
	// 基本用户
	m, err = contrib.AddAccess("visitor", "menu", access.VIEW, "/dashboard")
	assert.NoError(t, err)
	assert.Equal(t, m.Id, 2)
	m, err = contrib.AddAccess("visitor", "menu", access.VIEW, "/error/404")
	assert.NoError(t, err)
	assert.Equal(t, m.Id, 3)
}
