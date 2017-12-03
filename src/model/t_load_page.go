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
