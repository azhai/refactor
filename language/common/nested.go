package common

import (
	"xorm.io/builder"
	"xorm.io/xorm"
)

// 嵌套集合树
type NestedModel struct {
	Lft   int `json:"lft" xorm:"notnull default 0 comment('左边界') INT(10)"`           // 左边界
	Rgt   int `json:"rgt" xorm:"notnull default 0 comment('右边界') index INT(10)"`     // 右边界
	Depth int `json:"depth" xorm:"notnull default 1 comment('高度') index TINYINT(3)"` // 高度
}

// 是否叶子节点
func (n NestedModel) IsLeaf() bool {
	return n.Rgt-n.Lft == 1
}

// 有多少个子孙节点
func (n NestedModel) CountChildren() int {
	return int(n.Rgt-n.Lft-1) / 2
}

// 找出所有直系祖先节点
func (n NestedModel) AncestorsFilter(Backward bool) FilterFunc {
	return func(query *xorm.Session) *xorm.Session {
		query = query.Where("rgt > ? AND lft < ?", n.Rgt, n.Lft)
		if Backward { // 从子孙往祖先方向排序，即时间倒序
			return query.OrderBy("rgt ASC")
		} else {
			return query.OrderBy("rgt DESC")
		}
	}
}

// 找出所有子孙节点
func (n NestedModel) ChildrenFilter(rank uint8) FilterFunc {
	return func(query *xorm.Session) *xorm.Session {
		if n.Rgt > 0 && n.Lft > 0 { // 当前不是第0层，即具体某分支以下的节点
			query = query.Where("rgt < ? AND lft > ?", n.Rgt, n.Lft)
		}
		if rank > 0 { // 限制层级
			query = query.Where("depth < ?", uint8(n.Depth)+rank)
		}
		if rank != 1 { // 多层先按高度排序
			query = query.OrderBy("depth ASC")
		}
		return query.OrderBy("rgt ASC")

	}
}

// 添加到父节点最末，tbQuery一定要使用db.Table(...)
func (n *NestedModel) AddToParent(parent *NestedModel, tbQuery *xorm.Session) error {
	var query = tbQuery.OrderBy("rgt DESC")
	if parent == nil {
		n.Depth = 1
	} else {
		n.Depth = parent.Depth + 1
		query = query.Where("rgt < ? AND lft > ?", parent.Rgt, parent.Lft)
	}
	query = query.Where("depth = ?", n.Depth)
	sibling := new(NestedModel)
	has, err := query.Get(&sibling)
	if has == false || err != nil {
		return err
	}
	// 重建受影响的左右边界
	if sibling.Depth > 0 {
		n.Lft = sibling.Rgt + 1
	} else if parent != nil {
		n.Lft = parent.Lft + 1
		parent.Rgt += 2 // 上面的数据更新使 parent.Rgt 变成脏数据
	} else {
		n.Lft = 1
	}
	n.Rgt = n.Lft + 1
	if n.Depth > 1 {
		err = MoveEdge(tbQuery, uint(n.Lft), "+ 2")
	}
	return err
}

// 左右边界整体移动
func MoveEdge(query *xorm.Session, base uint, offset string) error {
	// 更新右边界
	query = query.Where("rgt >= ?", base) // 下面的更新lft也要用rgt作为索引
	affected, err := query.Update("rgt", builder.Expr("rgt "+offset))
	if affected == 0 || err != nil {
		return err
	}
	// 更新左边界，范围一定在上面更新右边界的所有行之内
	// 要么和上面一起为空，要么比上面少>=n行，n为直系祖先数量
	if affected > 1 {
		query = query.Where("lft >= ?", base)
		_, err = query.Update("lft", builder.Expr("lft "+offset))
	}
	return err
}
