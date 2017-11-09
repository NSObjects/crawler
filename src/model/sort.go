package model

type Sort struct {
	Id                int     `json:"id" xorm:"not null pk autoincr INT(11)"`
	SevenDaySalesRate float32 `json:"seven_day_sales_rate" xorm:"not null unique(seven_day_sales_rate) FLOAT(11,2)"`
	ProductId         int     `json:"product_id" xorm:"not null unique(seven_day_sales_rate) INT(30)"`
}
