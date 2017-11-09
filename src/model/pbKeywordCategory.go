package model

import (
	"time"
)

type PbKeywordCategory struct {
	Id         int       `json:"id" xorm:"not null pk autoincr INT(11)"`
	Keyword    string    `json:"keyword" xorm:"not null default '' unique(keyword) VARCHAR(200)"`
	WishId     string    `json:"wish_id" xorm:"not null default '' unique(keyword) index VARCHAR(30)"`
	CategoryId int       `json:"category_id" xorm:"not null unique(keyword) INT(11)"`
	Gender     int       `json:"gender" xorm:"not null unique(keyword) INT(11)"`
	Date       string    `json:"date" xorm:"not null default '' VARCHAR(30)"`
	Created    time.Time `json:"created" xorm:"DATETIME"`
}
