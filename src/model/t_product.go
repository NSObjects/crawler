package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type TProduct struct {
	Id          uint32    `orm:"column(id);pk"`
	RatingCount int       `orm:"column(rating_count)"`
	WishId      string    `orm:"column(wish_id);size(30)"`
	Price       float64   `orm:"column(price);digits(10);decimals(2)"`
	RetailPrice float64   `orm:"column(retail_price);digits(10);decimals(2)"`
	Shipping    float64   `orm:"column(shipping);digits(10);decimals(2)"`
	NumBought   int       `orm:"column(num_bought)"`
	NumEntered  int       `orm:"column(num_entered)"`
	Created     time.Time `orm:"column(created);auto_now_add;type(datetime)"`
	Updated     time.Time `orm:"column(updated);auto_now;type(datetime)"`
}

func (t *TProduct) TableName() string {
	return "t_product"
}

func init() {
	orm.RegisterModel(new(TProduct))
}
