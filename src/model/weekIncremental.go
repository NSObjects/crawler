package model

type WeekIncremental struct {
	Id                    int `json:"id" xorm:"not null pk autoincr INT(11)"`
	SalesIncremental      int `json:"sales_incremental" xorm:"not null index INT(11)"`
	CollectionIncremental int `json:"collection_incremental" xorm:"not null index INT(11)"`
	ProductId             int `json:"product_id" xorm:"index INT(30)"`
}
