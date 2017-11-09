package model

import (
	"time"
)

type Tag struct {
	Id            int       `json:"id" xorm:"not null pk autoincr INT(11)"`
	Name          string    `json:"name" xorm:"not null default '' VARCHAR(255)"`
	NumBought     int       `json:"num_bought" xorm:"not null default 0 INT(20)"`
	NumCollection int       `json:"num_collection" xorm:"not null default 0 INT(20)"`
	Price         float32   `json:"price" xorm:"not null default 0.00 FLOAT(10,2)"`
	Shipping      float32   `json:"shipping" xorm:"not null default 0.00 FLOAT(10,2)"`
	Rating        int       `json:"rating" xorm:"not null default 0 INT(20)"`
	Created       time.Time `json:"created" xorm:"DATETIME"`
	Update        time.Time `json:"update" xorm:"DATETIME"`
}
