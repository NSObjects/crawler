package model

import (
	"crawler/src/ini"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"

	"crawler/src/util"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/labstack/echo"
)

var (
	merchantLenght int64
)

const MERCHANT_NAME_CACHE = "mernchant_name"

var log = logrus.New()

func init() {
	file, err := os.OpenFile("err.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	local, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
	}
	time.Local = local
}

type TMerchant struct {
	Id                      uint32    `orm:"column(id)"`
	MerchantName            string    `orm:"column(merchant_name);size(255)"`
	ProductCount            int       `orm:"column(product_count)"`
	PercentPositiveFeedback float32   `orm:"column(percent_positive_feedback)"`
	DisplayName             string    `orm:"column(display_name);size(255)"`
	AvgRating               float32   `orm:"column(avg_rating)"`
	RatingCount             int       `orm:"column(rating_count)"`
	DisplayPic              string    `orm:"column(display_pic);size(255)"`
	ApprovedDate            time.Time `orm:"column(approved_date);type(datetime)"`
}

func (t *TMerchant) TableName() string {
	return "t_merchant"
}

func init() {
	orm.RegisterModel(new(TMerchant))
}

func (this *TMerchant) GetMerchantName() (wishIdJSON MerchantNameJSON, err error) {
	o := orm.NewOrm()

	if merchantLenght == 0 {
		merchantLenght, _ = ini.RedisClient.LLen(MERCHANT_NAME_CACHE).Result()
		if merchantLenght <= 0 {
			var list []orm.ParamsList
			num, err := o.Raw("select distinct(t_merchant_name) from merchant").ValuesList(&list)
			if err == nil && num > 0 {
				for _, name := range list {
					ini.RedisClient.RPush(MERCHANT_NAME_CACHE, fmt.Sprint(name))
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

	o := orm.NewOrm()

	for _, id := range dat.WishIds {
		if len(id) == 24 {
			wishID := TWishId{WishId: id}
			wishID.Created = time.Now()

			if _, err := o.Insert(&wishID); err != nil {
				if strings.Contains(err.Error(), "Error 1062: Duplicate entry") == false {
					log.WithFields(logrus.Fields{
						"t_merchant.go": "111",
					}).Error(err)
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

		if err := o.QueryTable("t_merchant").Filter("id", merchant.Id).One(&merchant); err != nil {

			if _, err := o.Insert(&merchant); err != nil {
				return err
			}
		}

		merchant.MerchantName = dat.MerchantInfo.Name
		merchant.DisplayPic = dat.MerchantInfo.DisplayPic
		merchant.DisplayName = dat.MerchantInfo.DisplayName
		merchant.AvgRating = dat.MerchantInfo.AvgRating
		merchant.RatingCount = dat.MerchantInfo.RatingCount
		merchant.PercentPositiveFeedback = dat.MerchantInfo.PercentPositiveFeedback

		if _, err := o.Update(&merchant); err != nil {
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
