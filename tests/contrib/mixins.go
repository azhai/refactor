package contrib

import (
	"strings"

	"gitea.com/azhai/refactor/builtin/access"
	db "gitea.com/azhai/refactor/tests/models/default"
)

type UserWithGroup struct {
	db.User   `json:",inline" xorm:"extends"`
	PrinGroup *GroupSummary `json:",inline" xorm:"extends"`
	ViceGroup *GroupSummary `json:",inline" xorm:"extends"`
}

type GroupSummary struct {
	Title  string `json:"title" xorm:"notnull default '' comment('名称') VARCHAR(50)"`
	Remark string `json:"remark" xorm:"comment('说明备注') TEXT"`
}

func (GroupSummary) TableName() string {
	return "t_group"
}

type Permission struct {
	db.Access `json:",inline" xorm:"extends"`
}

func (p Permission) CheckPerm(act uint16, url string) bool {
	if !p.RevokedAt.IsZero() || p.RoleName == "" || p.PermCode == 0 {
		return false
	}
	if p.ResourceArgs != "*" && !strings.HasPrefix(url, p.ResourceArgs) {
		return false
	}
	return access.ContainAction(uint16(p.PermCode), act, false)
}
