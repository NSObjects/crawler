package model

import (
	"time"
)

type ViewingStatistics struct {
	Id                 int       `json:"id" xorm:"not null pk autoincr INT(11)"`
	Median             int       `json:"median" xorm:"not null default 0 index index INT(11)"`
	MeanValue          int       `json:"mean_value" xorm:"not null default 0 INT(11)"`
	StandardDeviation  float64   `json:"standard_deviation" xorm:"not null default 0 index DOUBLE"`
	Range              int       `json:"range" xorm:"not null default 0 INT(11)"`
	MaximumValue       int       `json:"maximum_value" xorm:"not null default 0 INT(11)"`
	MinmunValue        int       `json:"minmun_value" xorm:"not null default 0 INT(11)"`
	StartValue         int       `json:"start_value" xorm:"not null default 0 INT(11)"`
	EndValue           int       `json:"end_value" xorm:"not null default 0 INT(11)"`
	StartEndDifference int       `json:"start_end_difference" xorm:"not null default 0 INT(11)"`
	Created            time.Time `json:"created" xorm:"not null DATETIME"`
	Updated            time.Time `json:"updated" xorm:"not null DATETIME"`
	Date               time.Time `json:"date" xorm:"not null index DATETIME"`
	ProductId          int       `json:"product_id" xorm:"not null default 0 index INT(10)"`
	Counts             string    `json:"counts" xorm:"not null TEXT"`
}
