package model

type Price struct {
	Id          int     `json:"id" xorm:"not null pk autoincr INT(11)"`
	RetailPrice float32 `json:"retail_price" xorm:"not null FLOAT"`
	Price       float32 `json:"price" xorm:"FLOAT"`
	ProductId   int     `json:"product_id" xorm:"not null INT(30)"`
	Date        int     `json:"date" xorm:"INT(11)"`
}
