package model

import (
	"time"
)

type ThreeSpeedUp struct {
	Id          int64     `json:"id" xorm:"pk autoincr BIGINT(20)"`
	ProductId   int       `json:"product_id" xorm:"not null index unique(p_c_s_unique) INT(30)"`
	SpeedUp     int       `json:"speed_up" xorm:"not null default 0 INT(20)"`
	Created     time.Time `json:"created" xorm:"not null pk unique(p_c_s_unique) DATETIME"`
	SpeedUpType int       `json:"speed_up_type" xorm:"not null unique(p_c_s_unique) INT(11)"`
	Start       time.Time `json:"start" xorm:"DATETIME"`
	End         time.Time `json:"end" xorm:"DATETIME"`
}
