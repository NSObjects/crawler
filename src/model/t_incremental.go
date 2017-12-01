package model

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type TIncremental struct {
	Id                       int       `json:"id" xorm:"not null pk autoincr INT(30)"`
	NumBought                int       `json:"num_bought" xorm:"INT(15)"`
	RatingCount              int       `json:"rating_count" xorm:"INT(15)"`
	NumCollection            int       `json:"num_collection" xorm:"INT(15)"`
	NumBoughtIncremental     int       `json:"num_bought_incremental" xorm:"INT(15)"`
	RatingCountIncremental   int       `json:"rating_count_incremental" xorm:"INT(15)"`
	NumCollectionIncremental int       `json:"num_collection_incremental" xorm:"INT(15)"`
	Price                    float64   `json:"price" xorm:"DECIMAL(10,2)"`
	Created                  time.Time `json:"created" xorm:"not null pk unique(pid_created) DATETIME"`
	Updated                  time.Time `json:"updated" xorm:"DATETIME"`
	ProductId                uint32    `json:"product_id" xorm:"not null unique(pid_created) index INT(30)"`
	PriceIncremental         float64   `json:"price_incremental" xorm:"DECIMAL(11,2)"`
}

func (t *TIncremental) TableName() string {
	return "t_incremental"
}

func init() {
	orm.RegisterModel(new(TIncremental))
}
