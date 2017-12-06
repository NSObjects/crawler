/*
 * Created  productCrawlerController.go on 17-12-4 下午3:42
 * Copyright (c) 2017  dyt.Co.Ltd All right reserved
 * Author lintao
 * Last modified 17-12-4 下午3:42
 */

package controller

import (
	"CrawlerMainController/utility"
	"crawler/src/global"
	"crawler/src/ini"
	"crawler/src/model"
	"crawler/src/util"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/now"
	"github.com/labstack/echo"

	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"strconv"

	"github.com/astaxie/beego/orm"
)

const size = 250

var (
	u            [][]model.TUser
	requestCount chan int
	mutex        sync.Mutex
	pageChan     chan int
)

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

type ProductCrawlerController struct{}

func (this ProductCrawlerController) RegisterRoute(g *echo.Group) {
	g.GET("/wishdata", this.GetWishId)
	g.POST("/wishdata", this.Post)
}

func (this ProductCrawlerController) GetWishId(ctx echo.Context) error {
	var JSONData WishIdJson
	JSONData.Code = 200
	o := orm.NewOrm()
	rc := <-requestCount
	if rc >= 10 {
		page := <-pageChan
		_, err := o.Raw("update t_load_page set all_id_page=?", page).Exec()
		if err != nil {
			log.WithFields(logrus.Fields{
				"productCrawlerController.go": "82",
			}).Error(err)
		}
		pageChan <- page
		requestCount <- 0
		rc = <-requestCount
	}

	var datas []string

	if rc >= 3 && rc < 8 && global.WeekSalesCacheLenght > 0 {
		datas = wishIdByWeekSalesGtZero()
	} else if rc >= 8 && global.SalesGreaterThanZeroCacheLenght > 0 {
		datas = wishIdBySalesGtZero()
	} else if global.AllWishIdCacheLenght > 0 {
		datas = allWishId()
	} else {
		datas = nocacheWishId()
	}

	rc += 1
	requestCount <- rc

	if len(datas) > 0 {
		JSONData.Data = datas
		JSONData.Users = model.GetUsers()
	}

	return ctx.JSON(http.StatusOK, JSONData)
}

func (this *ProductCrawlerController) Post(ctx echo.Context) error {

	var b []byte
	reader, err := gzip.NewReader(ctx.Request().Body)
	buf := bytes.NewBuffer(b)
	buf.ReadFrom(reader)

	if err != nil {
		log.WithFields(logrus.Fields{
			"productCrawlerController.go": "119",
		}).Error(err)
		return err
	}

	if len(buf.Bytes()) > 0 {
		ip := strings.Split(ctx.Request().RemoteAddr, ":")
		if len(ip) > 0 {
			if ip[0] != "[" {
				util.Statistics(0, ip[0])
			}
		}
		SaveProductToDBFrom(buf.Bytes())
	}

	return ctx.String(http.StatusOK, "ok")
}

func Setup() {

	requestCount = make(chan int, 1)
	requestCount <- 0
	o := orm.NewOrm()
	loadPage := model.TLoadPage{
		Id: 1,
	}
	if err := o.Read(&loadPage); err != nil {
		fmt.Println(err)
	}

	pageChan = make(chan int, 1000)
	pageChan <- loadPage.AllIdPage
}

func SaveProductToDBFrom(jsonStr []byte) {

	o := orm.NewOrm()
	var w WishProductJSON

	err := json.Unmarshal(jsonStr, &w)
	if err != nil {
		log.WithFields(logrus.Fields{
			"productCrawlerController.go": "164",
		}).Error(err)
		return
	}

	for _, j := range w.Data {

		if j.Code != 0 || len(j.Data.Contest.ID) <= 0 {
			continue
		}
		inserViewing(j)
		//先查redis中是否缓存了这个产品
		//如果没有就存一个快照

		id, _ := ini.RedisClient.HGet(global.SNAPSHOT_IDS, j.Data.Contest.ID).Result()

		if len(id) <= 0 && j.Data.Contest.NumBought > 0 {

			if value, err := json.Marshal(&j); err == nil {
				ps := model.TProductSnapshot{
					Data:    string(ZipBytes(value)),
					Created: now.BeginningOfDay(),
					WishID:  j.Data.Contest.ID,
				}

				_, err = o.Insert(&ps)
				//保存成功将Key设置为1
				if err != nil {
					if strings.Contains(err.Error(), "Duplicate entry") == true {
						ini.RedisClient.HSet(global.SNAPSHOT_IDS, j.Data.Contest.ID, "1")
					} else {
						log.WithFields(logrus.Fields{
							"productCrawlerController.go": "196",
						}).Error(err)
					}
				} else {
					ini.RedisClient.HSet(global.SNAPSHOT_IDS, j.Data.Contest.ID, "1")
				}

			} else {
				log.WithFields(logrus.Fields{
					"productCrawlerController.go": "205",
				}).Error(err)
			}
		}
		//
		//var product model.TProduct
		////查数据库中是否有这个产品
		////如果有新增一条增量，更新产品数据
		////没有就新增一条产品数据
		//
		//if err := o.QueryTable("t_product").Filter("wish_id", j.Data.Contest.ID).One(&product); err == nil {
		//	saveWishDataIncremental(j, &product)
		//	updateProduct(j, product)
		//} else {
		//	product.Created = time.Now()
		//	product.Id = util.FNV(j.Data.Contest.ID)
		//	configProduct(j, &product)
		//
		//	_, err = o.Insert(&product)
		//	if err != nil {
		//		if strings.Contains(err.Error(), " Duplicate entry") == false {
		//			log.WithFields(logrus.Fields{
		//				"productCrawlerController.go": "228",
		//			}).Error(err)
		//		}
		//	}
		//}

	}
}

func saveWishDataIncremental(jsonData model.WishOrginalData, product *model.TProduct) {

	if len(jsonData.Data.Contest.Name) <= 0 ||
		len(jsonData.Data.Contest.ID) <= 0 ||
		jsonData.Code != 0 {
		return
	}

	//如果这个产品更新时间距离现在超过一天，则不更新增量
	if time.Now().YearDay()-product.Updated.YearDay() > 1 {
		//return
	}

	wishdataIncremental := model.TIncremental{}
	if jsonData.Data.Contest.NumBought > product.NumBought {
		wishdataIncremental.NumBoughtIncremental = jsonData.Data.Contest.NumBought - product.NumBought
	}
	if jsonData.Data.Contest.NumEntered != product.NumEntered {
		wishdataIncremental.NumCollectionIncremental = jsonData.Data.Contest.NumEntered - product.NumEntered
	}
	if int(jsonData.Data.Contest.ProductRating.RatingCount) != product.RatingCount {
		wishdataIncremental.RatingCountIncremental = int(jsonData.Data.Contest.ProductRating.RatingCount) - product.RatingCount
	}

	for _, v := range jsonData.Data.Contest.CommerceProductInfo.Variations {
		if v.Price > 0 {
			if v.Price != product.Price {
				wishdataIncremental.PriceIncremental = v.Price - product.Price
				wishdataIncremental.Price = v.Price
			}
			break
		}
	}

	if wishdataIncremental.NumBoughtIncremental <= 0 &&
		wishdataIncremental.NumCollectionIncremental <= 0 {
		return
	}

	wishdataIncremental.Created = time.Now()
	wishdataIncremental.Updated = time.Now()
	wishdataIncremental.NumBought = jsonData.Data.Contest.NumBought
	wishdataIncremental.NumCollection = jsonData.Data.Contest.NumEntered
	wishdataIncremental.RatingCount = int(jsonData.Data.Contest.ProductRating.RatingCount)

	wishdataIncremental.ProductId = util.FNV(jsonData.Data.Contest.ID)
	o := orm.NewOrm()

	_, err := o.Insert(&wishdataIncremental)

	if err != nil {
		log.WithFields(logrus.Fields{
			"productCrawlerController.go": "289",
		}).Error(err)
	}
}

func updateProduct(jsonData model.WishOrginalData, product model.TProduct) {
	product.Updated = time.Now()
	if product.NumBought != jsonData.Data.Contest.NumBought ||
		product.NumEntered != jsonData.Data.Contest.NumEntered ||
		product.RatingCount != int(jsonData.Data.Contest.ProductRating.RatingCount) {

		configProduct(jsonData, &product)
		o := orm.NewOrm()

		if _, err := o.Update(&product); err != nil {
			log.WithFields(logrus.Fields{
				"productCrawlerController.go": "305",
			}).Error(err)
		}
	}
}

func configProduct(jsonData model.WishOrginalData, product *model.TProduct) {

	product.RatingCount = int(jsonData.Data.Contest.ProductRating.RatingCount)

	var price float64
	var retailPrice float64
	var shipping float64

	variations := jsonData.Data.Contest.CommerceProductInfo.Variations

	if len(variations) > 0 {
		retailPrice = variations[0].RetailPrice
		price = variations[0].Price
		shipping = variations[0].Shipping
		for _, v := range jsonData.Data.Contest.CommerceProductInfo.Variations {
			if v.RetailPrice < retailPrice {
				retailPrice = v.RetailPrice
			}

			if v.Shipping < shipping {
				shipping = v.Shipping
			}

			if v.Price < price {
				price = v.Price
			}
		}
	}

	product.Price = price
	product.RetailPrice = retailPrice
	product.Shipping = shipping
	product.WishId = jsonData.Data.Contest.ID
	product.NumEntered = jsonData.Data.Contest.NumEntered
	product.NumBought = jsonData.Data.Contest.NumBought
	product.Updated = time.Now()
}

func inserViewing(jsonData model.WishOrginalData) {
	if len(jsonData.Data.Contest.CurrentlyViewing.MessageList) > 0 {
		currentlyViewing := 0
		for _, v := range jsonData.Data.Contest.CurrentlyViewing.MessageList {
			for _, d := range strings.Split(v, " ") {
				if s, err := strconv.Atoi(d); err == nil {
					currentlyViewing += s
				}
			}
		}
		v := model.TViewings{
			Count:     currentlyViewing,
			ProductId: util.FNV(jsonData.Data.Contest.ID),
		}

		v.Created = time.Now()
		o := orm.NewOrm()

		if _, err := o.Insert(&v); err != nil {
			log.WithFields(logrus.Fields{
				"productCrawlerController.go": "332",
			}).Error(err)
		}
	}
}

func nocacheWishId() (datas []string) {

	o := orm.NewOrm()
	page := <-pageChan
	var list orm.ParamsList
	_, err := o.Raw("select wish_id from t_wish_id limit ? offset ?", size, size*page).ValuesFlat(&list)
	if err != nil {
		utility.Errorln(4, err)
	}
	if len(list) <= 0 {
		pageChan <- 0
		if _, err := ini.RedisClient.HSet("load_page", "page", 1).Result(); err != nil {
			log.WithFields(logrus.Fields{
				"productCrawlerController.go": "386",
			}).Error(err)
		}
		_, err = o.Raw("select wish_id from t_wish_id limit ? offset ?", size, 0).ValuesFlat(&list)
		if err != nil {
			log.WithFields(logrus.Fields{
				"productCrawlerController.go": "392",
			}).Error(err)
		}

	} else {
		pageChan <- page + 1
	}
	for _, id := range list {
		if wishId, ok := id.(string); ok {
			datas = append(datas, wishId)
		}
	}
	return datas
}

func allWishId() (datas []string) {
	if ids, err := ini.RedisClient.LRange(global.ALL_WISH_ID_CACHE, 0, 250).Result(); err == nil {
		ini.RedisClient.LTrim(global.ALL_WISH_ID_CACHE, 250, -1)
		return ids
	}

	return
}

func wishIdBySalesGtZero() (datas []string) {
	if ids, err := ini.RedisClient.LRange(global.SALES_GREATER_THAN_ZERO, 0, 250).Result(); err == nil {
		ini.RedisClient.LTrim(global.SALES_GREATER_THAN_ZERO, 250, -1)
		return ids
	}
	return
}

func wishIdByWeekSalesGtZero() (datas []string) {

	if ids, err := ini.RedisClient.LRange(global.WEEK_SALES_GREATER_THAN_ZERO, 0, 250).Result(); err == nil {
		ini.RedisClient.LTrim(global.WEEK_SALES_GREATER_THAN_ZERO, 250, -1)
		return ids
	}

	return datas
}

func ZipBytes(input []byte) []byte {
	var buf bytes.Buffer
	compressor, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	if err != nil {
		return input
	}
	compressor.Write(input)
	compressor.Close()
	return buf.Bytes()
}

type WishIdJson struct {
	Message string        `json:"message"`
	Code    int           `json:"code"`
	Data    []string      `json:"data"`
	Users   []model.TUser `json:"users"`
	Page    int64         `json:"page"`
}

type WishProductJSON struct {
	Data []model.WishOrginalData `json:"data"`
	Ip   string
}

type WishProduct struct {
	HasHwc                 int       `json:"has_hwc"`
	GenerationTime         time.Time `json:"generation_time"` //店铺开张时间
	RatingCount            int       `json:"rating_count"`
	Keyword                string    `json:"keyword"`
	Merchant               string    `json:"merchant"`
	MerchantName           string    `json:"merchant_name"`
	ContestSelectedPicture string    `json:"contest_selected_picture"`
	ExternalUrl            string    `json:"external_url"`
	Name                   string    `json:"name"`
	Countrys               string    `json:"countrys"`
	ExtraPhotoUrls         string    `json:"extra_photo_urls"`
	WishId                 string    `json:"wish_id"`
	Color                  string    `json:"color"`
	Size                   string    `json:"size"`
	Price                  float64   `json:"price"`
	RetailPrice            float64   `json:"retail_price"`
	MerchantTags           string    `json:"merchant_tags"`
	Tags                   string    `json:"tags"`
	Shipping               float64   `json:"shipping"`
	NumBought              int       `json:"num_bought"`
	MaxShippingTime        int       `json:"max_shipping_time"`
	MinShippingTime        int       `json:"min_shipping_time"`
	NumEntered             int       `json:"num_entered"`
	Code                   int       `json:"code"`
	Description            string    `json:"description"`
	Gender                 int       `json:"gender"`
	IsVerified             bool      `json:"is_verified"`
	CurrentlyViewing       int       `json:"currently_viewing"`
	Time                   int64     `json:"time"`
	TrueTagIds             []string  `json:"true_tag_ids"`
}
