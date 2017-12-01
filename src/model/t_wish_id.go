package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type TWishId struct {
	Id      uint32    `json:"id" xorm:"pk autoincr BIGINT(20)"`
	WishId  string    `json:"wish_id" xorm:"not null unique VARCHAR(30)"`
	Created time.Time `json:"created" xorm:"not null index DATETIME"`
}

func (t *TWishId) TableName() string {
	return "t_wish_id"
}

func init() {
	orm.RegisterModel(new(TWishId))
}
