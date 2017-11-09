package model

import (
	"time"
)

type ProductTag struct {
	Id        int       `json:"id" xorm:"not null pk autoincr INT(11)"`
	TagId     int64     `json:"tag_id" xorm:"not null default 0 unique(tag_id) BIGINT(20)"`
	ProductId int       `json:"product_id" xorm:"not null unique(tag_id) INT(30)"`
	Created   time.Time `json:"created" xorm:"not null DATETIME"`
	Updated   time.Time `json:"updated" xorm:"not null DATETIME"`
	WishId    string    `json:"wish_id" xorm:"not null default '' VARCHAR(255)"`
}
