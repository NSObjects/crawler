package model

import (
	"time"
)

type MerchantViewingMonitor struct {
	Id         int       `json:"id" xorm:"not null pk autoincr INT(11)"`
	UserId     int       `json:"user_id" xorm:"not null default 0 INT(11)"`
	MerchantId int       `json:"merchant_id" xorm:"not null default 0 INT(10)"`
	Created    time.Time `json:"created" xorm:"not null DATETIME"`
}
