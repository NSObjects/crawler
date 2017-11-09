package controller

import (
	"crawler/src/ini"
	"crawler/src/model"
	"crawler/src/util"
	"encoding/json"
	"math/rand"
	"net/http"

	"fmt"

	"time"

	"strings"

	"io/ioutil"

	"github.com/labstack/echo"
)

var (
	merchantLenght int64
)

type MerchantController struct{}

func (this MerchantController) RegisterRoute(g *echo.Group) {
	g.GET("/merchantCrawler", this.GetMerchant)
	g.POST("/merchantCrawler", this.Post)
}

const MERCHANT_NAME_CACHE = "mernchant_name"

func (this *MerchantController) GetMerchant(ctx echo.Context) error {

	if merchantLenght == 0 {
		merchantLenght, _ = ini.RedisClient.LLen(MERCHANT_NAME_CACHE).Result()
		if merchantLenght <= 0 {
			results, err := ini.AppWish.Query("select distinct(merchant_name) from merchant")

			if err != nil {
				return ctx.String(http.StatusInternalServerError, err.Error())
			}

			for _, r := range results {
				if string(r["merchant_name"]) != "" {
					err = ini.RedisClient.RPush(MERCHANT_NAME_CACHE, string(r["merchant_name"])).Err()
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}

	if len(u) <= 0 {
		return ctx.String(http.StatusOK, "")
	}
	var wishIdJSON wishIdJSON
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	users := u[r.Intn(len(u))]
	wishIdJSON.Users = users[r.Intn(len(users))]
	value, err := ini.RedisClient.LPop(MERCHANT_NAME_CACHE).Result()

	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	if len(value) == 0 {
		merchantLenght = 0
	}

	if len(value) > 0 {
		value = strings.Replace(value, "[", "", 1)
		value = strings.Replace(value, "]", "", 1)
		wishIdJSON.MerchntName = value
	}

	return ctx.JSON(http.StatusOK, wishIdJSON)
}

func (this *MerchantController) Post(ctx echo.Context) error {
	var dat merchantJSON
	b, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		fmt.Println(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	json.Unmarshal(b, &dat)

	for _, id := range dat.WishIds {
		if len(id) == 24 {
			wishID := model.WishId{WishId: id}
			wishID.Created = time.Now()
			if _, err := ini.AppWish.Insert(&wishID); err != nil {
				if !strings.Contains(err.Error(), "Duplicate entry") {
					fmt.Println(err)
				}
			}
		}
	}

	if len(dat.MerchantInfo.Name) > 0 {
		merchant := model.Merchant{Id: util.FNV(dat.MerchantInfo.Name)}
		if dat.MerchantInfo.ApprovedDate > 0 {
			merchant.ApprovedDate = time.Unix(int64(dat.MerchantInfo.ApprovedDate), 0)
		} else {
			merchant.ApprovedDate = time.Now()
		}

		ini.AppWish.Cols("")

		if exit, _ := ini.AppWish.Id(merchant.Id).Get(&merchant); exit == false {
			if _, err := ini.AppWish.Insert(&merchant); err != nil {
				fmt.Println(err)
				return ctx.String(http.StatusBadRequest, err.Error())
			}
		}

		merchant.MerchantName = dat.MerchantInfo.Name
		merchant.DisplayPic = dat.MerchantInfo.DisplayPic
		merchant.DisplayName = dat.MerchantInfo.DisplayName
		merchant.AvgRating = dat.MerchantInfo.AvgRating
		merchant.RatingCount = dat.MerchantInfo.RatingCount
		merchant.PercentPositiveFeedback = dat.MerchantInfo.PercentPositiveFeedback

		if _, err := ini.AppWish.Update(&merchant); err != nil {
			fmt.Println(err)
			return ctx.String(http.StatusBadRequest, err.Error())
		}
	}

	return ctx.String(http.StatusOK, "ok")
}

type merchantJSON struct {
	MerchantInfo struct {
		PercentPositiveFeedback float32 `json:"percent_positive_feedback"`
		DisplayName             string  `json:"display_name"`
		Name                    string  `json:"name"`
		DisplayPic              string  `json:"display_pic"`
		AvgRating               float32 `json:"avg_rating"`
		ApprovedDate            float32 `json:"approved_date"`
		RatingCount             int     `json:"rating_count"`
	} `json:"merchant_info"`
	WishIds []string `json:"wish_ids"`
}

type wishIdJSON struct {
	Users       model.User `json:"users"`
	MerchntName string     `json:"merchnt_name"`
}
