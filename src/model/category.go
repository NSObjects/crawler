package model

type Category struct {
	Id       int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	Name     string `json:"name" xorm:"not null default '' unique(name_filter_unique) VARCHAR(255)"`
	FilterId string `json:"filter_id" xorm:"not null default '' unique(name_filter_unique) VARCHAR(255)"`
}
