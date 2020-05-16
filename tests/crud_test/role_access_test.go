package crud_test

import (
	"testing"

	"gitea.com/azhai/refactor/contrib"
	"gitea.com/azhai/refactor/contrib/access"
	"gitea.com/azhai/refactor/contrib/usertype"
	_ "gitea.com/azhai/refactor/tests/models"
	db "gitea.com/azhai/refactor/tests/models/default"
	"github.com/stretchr/testify/assert"
)

func TestAccess01Insert(t *testing.T) {
	m := new(db.Access)
	err := contrib.TruncTable(m.TableName())
	assert.NoError(t, err)
	// 超管可以访问所有菜单
	m, err = contrib.AddAccess("superuser", "menu", access.ALL, "*")
	assert.NoError(t, err)
	assert.Equal(t, m.Id, 1)
	// 普通用户
	ops := access.VIEW | access.GET | access.POST
	m, err = contrib.AddAccess("member", "menu", ops, "/dashboard")
	assert.NoError(t, err)
	assert.Equal(t, m.Id, 2)
	m, err = contrib.AddAccess("member", "menu", access.VIEW, "/error/404")
	assert.NoError(t, err)
	assert.Equal(t, m.Id, 3)
	// 未登录用户
	m, err = contrib.AddAccess("", "menu", access.VIEW, "/error/404")
	assert.NoError(t, err)
	assert.Equal(t, m.Id, 4)
}

func TestAccess02Anonymous(t *testing.T) {
	var err error
	anonymous := &contrib.UserAuth{}
	err = usertype.Authorize(anonymous, access.VIEW, "style.css")
	assert.NoError(t, err)
	err = usertype.Authorize(anonymous, access.POST, "/images/abc.jpg")
	assert.NoError(t, err)
	err = usertype.Authorize(anonymous, access.POST, "/error/404")
	assert.NoError(t, err)
	err = usertype.Authorize(anonymous, access.POST, "/dashboard")
	assert.Error(t, err) // 无权限！
}

func TestAccess03Demo(t *testing.T) {
	var err error
	demo := &contrib.UserAuth{User: new(db.User)}
	demo.User.Load("username = ?", "demo")
	err = usertype.Authorize(demo, access.DISABLE, "/images/abc.jpg")
	assert.NoError(t, err)
	err = usertype.Authorize(demo, access.POST, "/dashboard")
	assert.NoError(t, err)
	err = usertype.Authorize(demo, access.VIEW, "/notExists")
	assert.Error(t, err) // 无权限！
}

func TestAccess04Admin(t *testing.T) {
	var err error
	admin := &contrib.UserAuth{User: new(db.User)}
	admin.User.Load("username = ?", "admin")
	err = usertype.Authorize(admin, access.POST, "/images/abc.jpg")
	assert.NoError(t, err)
	err = usertype.Authorize(admin, access.REMOVE, "/dashboard")
	assert.NoError(t, err)
	err = usertype.Authorize(admin, access.GET, "/notExists")
	assert.NoError(t, err)
	err = usertype.Authorize(admin, access.NONE, "")
	assert.NoError(t, err)
}
