package model

import (
	"time"
)

type EveryDayRank struct {
	Id        int       `json:"id" xorm:"not null pk autoincr INT(30)"`
	Position  int       `json:"position" xorm:"not null default 0 INT(20)"`
	ProductId int       `json:"product_id" xorm:"not null unique(p_r_c_unique) index INT(30)"`
	Created   time.Time `json:"created" xorm:"not null pk unique(p_r_c_unique) DATETIME"`
	RankType  int       `json:"rank_type" xorm:"not null unique(p_r_c_unique) INT(11)"`
}
