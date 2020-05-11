package crud_test

import (
	"testing"

	"gitea.com/azhai/refactor/tests/models"
	"gitea.com/azhai/refactor/tests/models/default"
	"github.com/stretchr/testify/assert"
)

var verbose = true

func init() {
	if c, ok := models.GetConnConfig("default"); ok {
		db.Initialize(c, verbose)
	}
}

func TestInsertRole(t *testing.T) {
	roles := []map[string]interface{}{
		{"name": "superuser", "remark": "超级用户，无上权限的超级管理员。"},
		{"name": "member", "remark": "普通用户，除权限外的其他页面。"},
		{"name": "visitor", "remark": "基本用户，只能看到面板页。"},
	}
	m := new(db.Role)
	for _, row := range roles {
		err := m.Save(row)
		assert.NoError(t, err)
	}
}
