package model

type ProductCategory struct {
	Id         int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	ProductId  int    `json:"product_id" xorm:"not null unique(product_id) INT(30)"`
	CategoryId string `json:"category_id" xorm:"not null default '' unique(product_id) VARCHAR(30)"`
}
