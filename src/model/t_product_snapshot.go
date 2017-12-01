package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type TProductSnapshot struct {
	Id      int       `json:"id" xorm:"not null pk autoincr INT(11)"`
	Data    string    `json:"data" xorm:"not null BLOB"`
	Created time.Time `json:"created" xorm:"not null unique(created_2) index DATETIME"`
	WishId  string    `json:"wish_id" xorm:"not null default '' index unique(created_2) VARCHAR(30)"`
}

func (t *TProductSnapshot) TableName() string {
	return "t_product_snapshot"
}

func init() {
	orm.RegisterModel(new(TProductSnapshot))
}
