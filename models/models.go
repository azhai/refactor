package models

type A struct {
	Id int `json:"Id" xorm:"integer"`
}

type B struct {
	Id int `json:"Id" xorm:"INTEGER"`
}
