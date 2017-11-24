package model

import (
	"crawler/src/ini"
	"math/rand"
	"time"

	"github.com/labstack/gommon/log"
)

type TUser struct {
	Id                    int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	Baid                  string `json:"baid" xorm:"not null default '' VARCHAR(255)"`
	SweeperSession        string `json:"sweeper_session" xorm:"not null default '' VARCHAR(255)"`
	Email                 string `json:"email" xorm:"not null default '' VARCHAR(255)"`
	Password              string `json:"password" xorm:"not null default '' VARCHAR(255)"`
	RiskifiedSessionToken string `json:"riskified_session_token" xorm:"not null default '' VARCHAR(255)"`
	AdvertiserId          string `json:"advertiser_id" xorm:"not null default '' VARCHAR(255)"`
	AppDeviceID           string `json:"app_device_i_d" xorm:"not null default '' VARCHAR(255)"`
	Country               string `json:"country" xorm:"not null default '' VARCHAR(255)"`
	FullName              string `json:"full_name" xorm:"not null default '' VARCHAR(255)"`
	HasAddress            int    `json:"has_address" xorm:"not null default 0 INT(11)"`
	Invalid               int    `json:"invalid" xorm:"not null default 0 INT(11)"`
	UserId                string `json:"user_id" xorm:"not null default '' VARCHAR(255)"`
	Gender                string `json:"gender" xorm:"VARCHAR(11)"`
}

var u [][]TUser

func GetUsers() []TUser {
	if len(u) <= 0 {
		contrys := []string{"Britain", "Canada", "Australia", "France", "Germany", "America"}
		for _, contry := range contrys {
			var user []TUser
			err := ini.AppWish.Where("has_address=1").And("country=?", contry).Find(&user)
			if err != nil {
				log.Print(err)
			}
			u = append(u, user)
		}
	}
	var users []TUser

	for _, userList := range u {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		if len(userList) <= 0 {
			continue
		}
		user := userList[r.Intn(len(userList))]
		if len(user.SweeperSession) > 0 && len(user.Baid) > 0 {
			users = append(users, user)
		}
	}
	return users
}
