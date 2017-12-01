package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type TProduct struct {
	Id          uint32    `json:"id" xorm:"not null pk INT(30)"`
	RatingCount int       `json:"rating_count" xorm:"not null default 0 INT(20)"`
	WishId      string    `json:"wish_id" xorm:"not null default '' index VARCHAR(30)"`
	Price       float64   `json:"price" xorm:"not null default 0.00 FLOAT(10,2)"`
	RetailPrice float64   `json:"retail_price" xorm:"not null default 0.00 FLOAT(10,2)"`
	Shipping    float64   `json:"shipping" xorm:"not null default 0.00 FLOAT(10,2)"`
	NumBought   int       `json:"num_bought" xorm:"not null default 0 index INT(20)"`
	NumEntered  int       `json:"num_entered" xorm:"not null default 0 index INT(20)"`
	Created     time.Time `json:"created" xorm:"not null index DATETIME"`
	Updated     time.Time `json:"updated" xorm:"not null DATETIME"`
}

func (t *TProduct) TableName() string {
	return "t_product"
}

func init() {
	orm.RegisterModel(new(TProduct))
}
