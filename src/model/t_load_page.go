/*
 * Created  t_load_page.go on 17-12-4 下午3:42
 * Copyright (c) 2017  dyt.Co.Ltd All right reserved
 * Author lintao
 * Last modified 17-12-3 下午6:13
 */

package model

import "github.com/astaxie/beego/orm"

type TLoadPage struct {
	Id                 int `orm:"column(id);auto"`
	WeeSalesPage       int `orm:"column(week_sales_page)"`
	AllIdPage          int `orm:"column(all_id_page)"`
	SalesGtZeroPage    int `orm:"column(sales_gt_zero_page)"`
	AllWishIdCachePage int `orm:"column(all_wishid_cache_page)"`
}

func init() {
	orm.RegisterModel(new(TLoadPage))
}
