/*
 * Created  t_wish_id.go on 17-12-4 下午3:42
 * Copyright (c) 2017  dyt.Co.Ltd All right reserved
 * Author lintao
 * Last modified 17-12-3 下午6:16
 */

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
