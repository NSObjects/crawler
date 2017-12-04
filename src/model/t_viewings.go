/*
 * Created  t_viewings.go on 17-12-4 下午3:42
 * Copyright (c) 2017  dyt.Co.Ltd All right reserved
 * Author lintao
 * Last modified 17-12-3 下午6:16
 */

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
