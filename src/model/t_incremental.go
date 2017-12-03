package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type TIncremental struct {
	Id                       int       `orm:"column(id);auto"`
	NumBought                int       `orm:"column(num_bought);null"`
	RatingCount              int       `orm:"column(rating_count);null"`
	NumCollection            int       `orm:"column(num_collection);null"`
	NumBoughtIncremental     int       `orm:"column(num_bought_incremental);null"`
	RatingCountIncremental   int       `orm:"column(rating_count_incremental);null"`
	NumCollectionIncremental int       `orm:"column(num_collection_incremental);null"`
	Price                    float64   `orm:"column(price);null;digits(10);decimals(2)"`
	Created                  time.Time `orm:"column(created);type(datetime)"`
	Updated                  time.Time `orm:"column(updated);type(datetime);null"`
	ProductId                uint32    `orm:"column(product_id);size(30);null"`
	PriceIncremental         float64   `orm:"column(price_incremental);null;digits(11);decimals(2)"`
}

func init() {
	orm.RegisterModel(new(TIncremental))
}
