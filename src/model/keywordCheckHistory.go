package model

import (
	"time"
)

type KeywordCheckHistory struct {
	Id           int       `json:"id" xorm:"not null pk autoincr INT(11)"`
	WishId       string    `json:"wish_id" xorm:"not null default '' VARCHAR(255)"`
	Data         string    `json:"data" xorm:"not null TEXT"`
	Gender       string    `json:"gender" xorm:"not null default '' VARCHAR(255)"`
	IsNewAccount int       `json:"is_new_account" xorm:"not null default 0 INT(11)"`
	Created      time.Time `json:"created" xorm:"not null DATETIME"`
}
