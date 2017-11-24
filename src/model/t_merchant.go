package model

import (
	"crawler/src/ini"
	"log"

	"crawler/src/util"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/labstack/echo"
)

var (
	merchantLenght int64
)

const MERCHANT_NAME_CACHE = "mernchant_name"

type TMerchant struct {
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

func (this *TMerchant) GetMerchantName() (wishIdJSON MerchantNameJSON, err error) {
	if merchantLenght == 0 {
		merchantLenght, _ = ini.RedisClient.LLen(MERCHANT_NAME_CACHE).Result()
		if merchantLenght <= 0 {
			results, err := ini.AppWish.Query("select distinct(merchant_name) from merchant")

			if err != nil {
				return wishIdJSON, err
			}
			for _, r := range results {
				if string(r["merchant_name"]) != "" {
					err = ini.RedisClient.RPush(MERCHANT_NAME_CACHE, string(r["merchant_name"])).Err()
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	users := GetUsers()
	wishIdJSON.Users = users[r.Intn(len(users))]
	value, err := ini.RedisClient.LPop(MERCHANT_NAME_CACHE).Result()

	if err != nil {
		return wishIdJSON, err
	}
	if len(value) == 0 {
		merchantLenght = 0
	}

	if len(value) > 0 {
		value = strings.Replace(value, "[", "", 1)
		value = strings.Replace(value, "]", "", 1)
		wishIdJSON.MerchntName = value
	}

	return
}

func (this *TMerchant) MerchantInfoHandler(ctx echo.Context) error {
	var dat MerchantJSON
	b, err := ioutil.ReadAll(ctx.Request().Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &dat)
	if err != nil {
		return err
	}

	for _, id := range dat.WishIds {
		if len(id) == 24 {
			wishID := TWishId{WishId: id}
			wishID.Created = time.Now()
			if _, err := ini.AppWish.Insert(&wishID); err != nil {
				if strings.Contains(err.Error(), "Error 1062: Duplicate entry") == false {
					log.Println(err)
					log.Println(err)
				}

			}
		}
	}

	if len(dat.MerchantInfo.Name) > 0 {
		merchant := TMerchant{Id: util.FNV(dat.MerchantInfo.Name)}
		if dat.MerchantInfo.ApprovedDate > 0 {
			merchant.ApprovedDate = time.Unix(int64(dat.MerchantInfo.ApprovedDate), 0)
		} else {
			merchant.ApprovedDate = time.Now()
		}

		if exit, _ := ini.AppWish.Id(merchant.Id).Get(&merchant); exit == false {
			if _, err := ini.AppWish.Insert(&merchant); err != nil {
				return err
			}
		}

		merchant.MerchantName = dat.MerchantInfo.Name
		merchant.DisplayPic = dat.MerchantInfo.DisplayPic
		merchant.DisplayName = dat.MerchantInfo.DisplayName
		merchant.AvgRating = dat.MerchantInfo.AvgRating
		merchant.RatingCount = dat.MerchantInfo.RatingCount
		merchant.PercentPositiveFeedback = dat.MerchantInfo.PercentPositiveFeedback

		if _, err := ini.AppWish.Update(&merchant); err != nil {
			return err
		}
	}

	return nil
}

type MerchanApiInterface interface {
	GetMerchantName() (MerchantNameJSON, error)
	MerchantInfoHandler(ctx echo.Context) error
}

type MerchantJSON struct {
	MerchantInfo MerchantInfo `json:"merchant_info"`
	WishIds      []string     `json:"wish_ids"`
}

type MerchantNameJSON struct {
	Users       TUser  `json:"users"`
	MerchntName string `json:"merchnt_name"`
}

type MerchantInfo struct {
	PercentPositiveFeedback float32 `json:"percent_positive_feedback"`
	DisplayName             string  `json:"display_name"`
	Name                    string  `json:"name"`
	DisplayPic              string  `json:"display_pic"`
	AvgRating               float32 `json:"avg_rating"`
	ApprovedDate            float32 `json:"approved_date"`
	RatingCount             int     `json:"rating_count"`
}
