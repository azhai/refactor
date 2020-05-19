package contrib

import (
	"strings"

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

type CronTimerCluster struct {
	cron.CronTimer     `json:",inline" xorm:"extends"`
	*base.ClusterMixin `json:"-" xorm:"-"`
}

func NewCronTimerCluster() *CronTimerCluster {
	ct := cron.CronTimer{}
	prefix := ct.TableName() + "_"
	cm := base.GetClusterMixinFor("Monthly", prefix, cron.Engine())
	curr := prefix + cm.GetSuffix()
	m := &CronTimerCluster{ct, cm}
	m.ResetTable(curr, false)
	return m
}

func (m CronTimerCluster) TableName() string {
	table := m.CronTimer.TableName()
	if m.ClusterMixin == nil {
		return table
	}
	suffix := m.ClusterMixin.GetSuffix()
	return table + "_" + suffix
}

func (m CronTimerCluster) ResetTable(curr string, trunc bool) error {
	table := m.CronTimer.TableName()
	create, err := base.CreateTableLike(cron.Engine(), curr, table)
	if err == nil && !create && trunc {
		err = TruncTable(curr)
	}
	return err
}
