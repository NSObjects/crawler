package model

import (
	"time"
)

type ProductVariations struct {
	Id            int       `json:"id" xorm:"not null pk autoincr INT(11)"`
	ProductId     int       `json:"product_id" xorm:"not null unique(product_id_2) index INT(30)"`
	Created       time.Time `json:"created" xorm:"not null unique(product_id_2) index DATETIME"`
	WishShipping  int       `json:"wish_shipping" xorm:"not null INT(11)"`
	OwnerShipping int       `json:"owner_shipping" xorm:"not null INT(11)"`
	Title         int       `json:"title" xorm:"not null INT(11)"`
	WishPrice     int       `json:"wish_price" xorm:"not null INT(11)"`
	OwnerPrice    int       `json:"owner_price" xorm:"not null INT(11)"`
	Tag           int       `json:"tag" xorm:"not null INT(11)"`
	Category      int       `json:"category" xorm:"not null INT(11)"`
	WishVerified  int       `json:"wish_verified" xorm:"not null INT(11)"`
	IsVariations  int       `json:"is_variations" xorm:"not null default 0 index INT(11)"`
	Sku           int       `json:"sku" xorm:"not null INT(11)"`
	NumBought     int       `json:"num_bought" xorm:"not null INT(11)"`
	MerchantTags  int       `json:"merchant_tags" xorm:"default 0 INT(11)"`
}
