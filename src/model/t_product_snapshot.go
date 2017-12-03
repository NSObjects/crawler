package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type TProductSnapshot struct {
	Id      uint32    `orm:"column(id);pk"`
	Data    string    `orm:"column(data)"`
	WishID  string    `orm:"column(wish_id)"`
	Created time.Time `orm:"column(created);type(datetime)"`
}

func (u *TProductSnapshot) TableName() string {
	return "t_product_snapshot"
}
func init() {
	orm.RegisterModel(new(TProductSnapshot))
}
