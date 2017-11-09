package model

import (
	"time"
)

type CategoriesTag struct {
	Id           int       `json:"id" xorm:"not null pk autoincr INT(11)"`
	TagId        int       `json:"tag_id" xorm:"not null default 0 unique(tag_categores) INT(20)"`
	CategoriesId string    `json:"categories_id" xorm:"not null default '' unique(tag_categores) VARCHAR(255)"`
	Created      time.Time `json:"created" xorm:"not null DATETIME"`
	Updated      time.Time `json:"updated" xorm:"not null DATETIME"`
}
