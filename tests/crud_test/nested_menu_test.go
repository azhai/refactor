package crud_test

import (
	"testing"

	"gitea.com/azhai/refactor/tests/contrib"
	_ "gitea.com/azhai/refactor/tests/models"
	db "gitea.com/azhai/refactor/tests/models/default"
	"github.com/stretchr/testify/assert"
)

var allMenuData = []map[string]string{
	{"path": "/dashboard", "title": "面板", "icon": "dashboard"},
	{"path": "/permission", "title": "权限", "icon": "lock"},
	{"path": "role", "title": "角色权限"},
	{"path": "/table", "title": "Table", "icon": "table"},
	{"path": "complex-table", "title": "复杂Table"},
	{"path": "inline-edit-table", "title": "内联编辑"},
	{"path": "/excel", "title": "Excel", "icon": "excel"},
	{"path": "export-selected-excel", "title": "选择导出"},
	{"path": "upload-excel", "title": "上传Excel"},
	{"path": "/theme/index", "title": "主题", "icon": "theme"},
	{"path": "/error/404", "title": "404错误", "icon": "404"},
	{"path": "https://cn.vuejs.org/", "title": "外部链接", "icon": "link"},
}

// 添加菜单
func TestNested01InsertMenus(t *testing.T) {
	parent := new(db.Menu)
	err := contrib.TruncTable(parent.TableName())
	assert.NoError(t, err)
	icon, ok := "", false
	for _, row := range allMenuData {
		if icon, ok = row["icon"]; ok && icon != "" {
			parent = contrib.NewMenu(row["path"], row["title"], icon)
			err := contrib.AddMenuToParent(parent, nil)
			assert.NoError(t, err)
		} else {
			menu := contrib.NewMenu(row["path"], row["title"], icon)
			err := contrib.AddMenuToParent(menu, parent)
			assert.NoError(t, err)
		}
	}
}

// 查找祖先菜单
func TestNested02FindAncestors(t *testing.T) {
	menu := new(db.Menu)
	_, err := menu.Load("path = ?", "inline-edit-table")
	assert.NoError(t, err)
	filter := menu.AncestorsFilter(true)
	var menus []*db.Menu
	err = filter(db.Table(menu)).Find(&menus)
	assert.NoError(t, err)
	if assert.Len(t, menus, 1) {
		assert.Equal(t, menus[0].Path, "/table")
	}
}

// 查找子菜单
func TestNested03FindChildren(t *testing.T) {
	menu := new(db.Menu)
	_, err := menu.Load("path = ?", "/excel")
	assert.NoError(t, err)
	filter := menu.ChildrenFilter(-1)
	var menus []*db.Menu
	err = filter(db.Table(menu)).Find(&menus)
	assert.NoError(t, err)
	if assert.Len(t, menus, 2) {
		assert.Equal(t, menus[0].Path, "export-selected-excel")
		assert.Equal(t, menus[1].Path, "upload-excel")
	}
}
