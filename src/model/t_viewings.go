package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type TViewings struct {
	Id        int
	Created   time.Time
	Count     int    `orm:"column(count)"`
	ProductId uint32 `orm:"column(product_id)"`
}

func init() {
	orm.RegisterModel(new(TViewings))
}

func (t *TViewings) TableName() string {
	return "t_viewings"
}
