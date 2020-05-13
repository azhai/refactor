package contrib

import (
	"encoding/hex"
	"time"

	db "gitea.com/azhai/refactor/tests/models/default"
	"github.com/azhai/gozzo-utils/cryptogy"
	"github.com/muyo/sno"
)

var saltPasswd ICipher

type ICipher interface {
	CreatePassword(plainText string) string
	VerifyPassword(plainText, cipherText string) bool
}

func Cipher() ICipher {
	if saltPasswd == nil { // 8位salt值，用$符号分隔开
		saltPasswd = cryptogy.NewSaltPassword(8, "$")
	}
	return saltPasswd
}

func NewSerialNo(n byte) string {
	return sno.New(n).String()
}

func NewTimeSerialNo(n byte, t time.Time) string {
	return sno.NewWithTime(n, t).String()
}

func NewToken(n byte) string {
	tail := cryptogy.RandSalt(2)
	token := sno.New(n).Bytes()
	token = append(token[:8], []byte(tail)...)
	return hex.EncodeToString(token)
}

type UserWithGroup struct {
	db.User   `xorm:"extends"`
	PrinGroup *GroupSummary `xorm:"extends"`
	ViceGroup *GroupSummary `xorm:"extends"`
}

func (UserWithGroup) TableName() string {
	return "t_user"
}

type GroupSummary struct {
	Title  string `json:"title" xorm:"notnull default '' comment('名称') VARCHAR(50)"`
	Remark string `json:"remark" xorm:"comment('说明备注') TEXT"`
}

func (GroupSummary) TableName() string {
	return "t_group"
}
