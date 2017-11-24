package model

import (
	"time"
)

type TViewings struct {
	Id        int       `json:"id" xorm:"not null pk autoincr INT(11)"`
	Created   time.Time `json:"created" xorm:"not null index unique(created) DATETIME"`
	Count     int       `json:"count" xorm:"not null default 0 index INT(11)"`
	ProductId uint32    `json:"product_id" xorm:"not null index unique(created) INT(30)"`
}
