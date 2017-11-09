package model

import (
	"time"
)

type ProductSnapshot struct {
	Id      int       `json:"id" xorm:"not null pk autoincr INT(11)"`
	Created time.Time `json:"created" xorm:"not null unique(created_2) index DATETIME"`
	WishId  string    `json:"wish_id" xorm:"not null default '' index unique(created_2) VARCHAR(30)"`
	Data    string    `json:"data" xorm:"not null LONGTEXT"`
}
