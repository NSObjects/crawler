package model

import (
	"time"
)

type PbKeywordRecords struct {
	Id                  int       `json:"id" xorm:"not null pk autoincr INT(11)"`
	Keyword             string    `json:"keyword" xorm:"default '' index(keyword) VARCHAR(200)"`
	KeywordHash         string    `json:"keyword_hash" xorm:"TEXT"`
	Reach               int       `json:"reach" xorm:"index(competition) INT(11)"`
	ReachIncrease       float32   `json:"reach_increase" xorm:"FLOAT(7,2)"`
	ReachText           string    `json:"reach_text" xorm:"index(keyword) VARCHAR(20)"`
	Competition         int       `json:"competition" xorm:"index(competition) INT(11)"`
	CompetitionIncrease float32   `json:"competition_increase" xorm:"FLOAT(7,2)"`
	CompetitionText     string    `json:"competition_text" xorm:"index(keyword) VARCHAR(20)"`
	HighBid             float32   `json:"high_bid" xorm:"index(competition) FLOAT(5,2)"`
	BidIncrease         float32   `json:"bid_increase" xorm:"FLOAT(7,2)"`
	Date                string    `json:"date" xorm:"index(competition) VARCHAR(12)"`
	Created             time.Time `json:"created" xorm:"default 'CURRENT_TIMESTAMP' TIMESTAMP"`
}
