package contrib

import (
	"strings"
	"time"

	"gitea.com/azhai/refactor/builtin/access"
	"gitea.com/azhai/refactor/builtin/base"
	"gitea.com/azhai/refactor/tests/models/cron"
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

type CronTimerMonthly struct {
	cron.CronTimer     `json:",inline" xorm:"extends"`
	*base.MonthlyMixin `json:"-" xorm:"-"`
}

func (m CronTimerMonthly) TableName() string {
	if m.MonthlyMixin == nil {
		m.MonthlyMixin = base.NewMonthlyMixin(time.Now())
	}
	table := m.CronTimer.TableName()
	suffix := m.MonthlyMixin.GetSuffix()
	return table + "_" + suffix
}

func (m CronTimerMonthly) ResetTable(curr string, trunc bool) error {
	table := m.CronTimer.TableName()
	create, err := base.CreateTableLike(cron.Engine(), curr, table)
	if err == nil && !create && trunc {
		err = TruncTable(curr)
	}
	return err
}
