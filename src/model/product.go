package model

import (
	"time"
)

type Product struct {
	Id              uint32    `json:"id" xorm:"not null pk INT(30)"`
	GenerationTime  time.Time `json:"generation_time" xorm:"DATETIME"`
	RatingCount     int       `json:"rating_count" xorm:"not null default 0 INT(20)"`
	Name            string    `json:"name" xorm:"not null default '' VARCHAR(600)"`
	WishId          string    `json:"wish_id" xorm:"not null default '' index VARCHAR(30)"`
	Color           string    `json:"color" xorm:"TEXT"`
	Size            string    `json:"size" xorm:"TEXT"`
	Price           float64   `json:"price" xorm:"not null default 0.00 FLOAT(10,2)"`
	RetailPrice     float64   `json:"retail_price" xorm:"not null default 0.00 FLOAT(10,2)"`
	MerchantTags    string    `json:"merchant_tags" xorm:"TEXT"`
	Tags            string    `json:"tags" xorm:"TEXT"`
	Shipping        float64   `json:"shipping" xorm:"not null default 0.00 FLOAT(10,2)"`
	Description     string    `json:"description" xorm:"not null TEXT"`
	NumBought       int       `json:"num_bought" xorm:"not null default 0 index INT(20)"`
	MaxShippingTime int       `json:"max_shipping_time" xorm:"not null default 0 INT(11)"`
	MinShippingTime int       `json:"min_shipping_time" xorm:"not null default 0 INT(11)"`
	NumEntered      int       `json:"num_entered" xorm:"not null default 0 index INT(20)"`
	Created         time.Time `json:"created" xorm:"not null index DATETIME"`
	Updated         time.Time `json:"updated" xorm:"not null DATETIME"`
	Gender          int       `json:"gender" xorm:"index INT(11)"`
	MerchantId      uint32    `json:"merchant_id" xorm:"not null index INT(30)"`
	Merchant        string    `json:"merchant" xorm:"not null default '' VARCHAR(255)"`
	TrueTagIds      string    `json:"true_tag_ids" xorm:"VARCHAR(4000)"`
}
