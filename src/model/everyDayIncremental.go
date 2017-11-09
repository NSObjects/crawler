package model

import (
	"time"
)

type EveryDayIncremental struct {
	Id                     int       `json:"id" xorm:"not null pk autoincr INT(30)"`
	BeginningOfDay         time.Time `json:"beginning_of_day" xorm:"DATETIME"`
	EndOfDay               time.Time `json:"end_of_day" xorm:"DATETIME"`
	Created                time.Time `json:"created" xorm:"not null pk index DATETIME"`
	Updated                time.Time `json:"updated" xorm:"not null DATETIME"`
	ProductId              int       `json:"product_id" xorm:"not null index INT(30)"`
	PriceIncremental       float32   `json:"price_incremental" xorm:"not null default 0.00 FLOAT(10,2)"`
	Price                  float32   `json:"price" xorm:"not null FLOAT(10,2)"`
	NumEnteredIncremental  int       `json:"num_entered_incremental" xorm:"not null INT(11)"`
	NumEntered             int       `json:"num_entered" xorm:"not null INT(11)"`
	NumBoughtIncremental   int       `json:"num_bought_incremental" xorm:"INT(11)"`
	NumBought              int       `json:"num_bought" xorm:"not null index INT(11)"`
	RatingCountIncremental int       `json:"rating_count_incremental" xorm:"not null INT(11)"`
	RatingCount            int       `json:"rating_count" xorm:"not null INT(11)"`
}
