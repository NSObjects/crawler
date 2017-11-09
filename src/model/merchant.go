package model

import (
	"time"
)

type Merchant struct {
	Id                      uint32    `json:"id" xorm:"not null pk autoincr INT(255)"`
	MerchantName            string    `json:"merchant_name" xorm:"not null default '' unique VARCHAR(255)"`
	ProductCount            int       `json:"product_count" xorm:"not null INT(255)"`
	PercentPositiveFeedback float32   `json:"percent_positive_feedback" xorm:"not null FLOAT(255,2)"`
	DisplayName             string    `json:"display_name" xorm:"not null default '' VARCHAR(255)"`
	AvgRating               float32   `json:"avg_rating" xorm:"not null FLOAT(255,2)"`
	RatingCount             int       `json:"rating_count" xorm:"not null INT(255)"`
	DisplayPic              string    `json:"display_pic" xorm:"not null default '' VARCHAR(255)"`
	ApprovedDate            time.Time `json:"approved_date" xorm:"index DATETIME"`
}
