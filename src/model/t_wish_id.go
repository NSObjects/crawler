package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type TWishId struct {
	Id      uint32    `orm:"column(id);pk"`
	WishId  string    `orm:"column(wish_id);size(30)"`
	Created time.Time `orm:"column(created);type(datetime)"`
}

func (t *TWishId) TableName() string {
	return "t_wish_id"
}

func init() {
	orm.RegisterModel(new(TWishId))
}
