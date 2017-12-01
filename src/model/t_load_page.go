package model

type TLoadPage struct {
	Id                 int `json:"id" xorm:"not null pk autoincr INT(11)"`
	WeekSalesPage      int `json:"week_sales_page" xorm:"INT(11)"`
	AllIdPage          int `json:"all_id_page" xorm:"INT(11)"`
	SalesGtZeroPage    int `json:"sales_gt_zero_page" xorm:"INT(11)"`
	AllWishidCachePage int `json:"all_wishid_cache_page" xorm:"INT(11)"`
}
