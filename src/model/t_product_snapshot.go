/*
 * Created  t_product_snapshot.go on 17-12-4 下午3:42
 * Copyright (c) 2017  dyt.Co.Ltd All right reserved
 * Author lintao
 * Last modified 17-12-3 下午6:15
 */

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
