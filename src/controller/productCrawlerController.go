package controller

import (
	"crawler/src/global"
	"crawler/src/ini"
	"crawler/src/model"
	"crawler/src/util"

	"go.uber.org/zap"

	"github.com/jinzhu/now"
	"github.com/labstack/echo"

	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const size int = 250

var (
	u                 [][]model.TUser
	requestCount      int
	mutex             sync.Mutex
	weekSalesPageChan chan int
	pageChan          chan int
)

type ProductCrawlerController struct{}

func (this ProductCrawlerController) RegisterRoute(g *echo.Group) {
	g.GET("/wishdata", this.GetWishId)
	g.POST("/wishdata", this.Post)
}

func (this ProductCrawlerController) GetWishId(ctx echo.Context) error {
	var JSONData WishIdJson
	JSONData.Code = 200
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if requestCount >= 10 {
		weekSalesPage := <-weekSalesPageChan
		weekSalesPageChan <- weekSalesPage
		_, err := ini.AppWish.Exec("update t_load_page set week_sales_page=?", weekSalesPage)
		if err != nil {
			logger.Error(err.Error())
		}
		page := <-pageChan
		_, err = ini.AppWish.Exec("update t_load_page set all_id_page=?", page)
		if err != nil {
			logger.Error(err.Error())
		}
		pageChan <- page
		requestCount = 0
	}

	var datas []string
	mutex.Lock()
	if requestCount >= 5 && global.WeekSalesCacheLenght > 0 {
		datas = wishIdByWeekSalesGtZero()
	} else if requestCount >= 8 && global.SalesGreaterThanZeroCacheLenght > 0 {
		datas = wishIdBySalesGtZero()
	} else if global.AllWishIdCacheLenght > 0 {
		datas = allWishId()
	} else {
		datas = nocacheWishId()
	}
	mutex.Unlock()
	requestCount++
	if len(datas) > 0 {
		JSONData.Data = datas
		JSONData.Users = model.GetUsers()
	}

	return ctx.JSON(http.StatusOK, JSONData)
}

func (this *ProductCrawlerController) Post(ctx echo.Context) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	b, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		logger.Error(err.Error())
	}

	if len(b) > 0 {
		ip := strings.Split(ctx.Request().RemoteAddr, ":")
		if len(ip) > 0 {
			if ip[0] != "[" {
				fmt.Println(ip[0])
			}
		}

		SaveProductToDBFrom(b)
	}

	return ctx.String(http.StatusOK, "ok")
}

func Setup() {
	loadPage := &model.TLoadPage{Id: 1}
	_, err := ini.AppWish.Get(loadPage)
	if err != nil {
		panic(err)
	}

	weekSalesPageChan = make(chan int, 50)
	weekSalesPageChan <- loadPage.WeekSalesPage

	pageChan = make(chan int, 50)
	pageChan <- loadPage.AllIdPage
}

func SaveProductToDBFrom(jsonStr []byte) {
	logger, _ := zap.NewProduction()

	defer logger.Sync()
	var w WishProductJSON

	err := json.Unmarshal(jsonStr, &w)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	fmt.Printf("接收到数据%d条", len(w.Data))
	//logger.Debug(fmt.Sprintf("接收到数据%条", len(w.Data)))

	for _, j := range w.Data {

		if j.Code != 0 || len(j.Data.Contest.ID) <= 0 {
			continue
		}

		id, _ := ini.RedisClient.HGet(global.SNAPSHOT_IDS, j.Data.Contest.ID).Result()

		if len(id) <= 0 {

			ini.RedisClient.HSet(global.SNAPSHOT_IDS, j.Data.Contest.ID, "1")
			value, err := json.Marshal(&j)
			if err != nil {
				logger.Error(err.Error())
				continue
			}

			ps := model.TProductSnapshot{
				Data:    string(ZipBytes(value)),
				Created: now.BeginningOfDay(),
				WishId:  j.Data.Contest.ID,
			}
			_, err = ini.AppWish.Insert(&ps)
			if err != nil {
				logger.Error(err.Error())
			}

			var product model.TProduct
			product.Created = time.Now()
			product.Id = util.FNV(j.Data.Contest.ID)
			configProduct(j, &product)
			_, err = ini.AppWish.Insert(&product)
			if err != nil {
				logger.Error(err.Error())
			}

		} else {
			var product model.TProduct

			if _, err := ini.AppWish.Id(util.FNV(j.Data.Contest.ID)).Get(&product); err == nil {
				saveWishDataIncremental(j, product)
			} else {
				logger.Error(err.Error())
			}
		}
	}
}

func saveWishDataIncremental(jsonData model.WishOrginalData, product model.TProduct) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if len(jsonData.Data.Contest.Name) <= 0 || len(jsonData.Data.Contest.ID) <= 0 || jsonData.Code != 0 {
		return
	}

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

		if _, err := ini.AppWish.Insert(&v); err != nil {
			logger.Error(err.Error())
		}
	}

	wishdataIncremental := model.TIncremental{}
	if time.Now().YearDay()-product.Updated.YearDay() <= 1 {
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

		if wishdataIncremental.NumBoughtIncremental > 0 ||
			wishdataIncremental.NumCollectionIncremental > 0 {
			wishdataIncremental.Created = time.Now()
			wishdataIncremental.Updated = time.Now()
			wishdataIncremental.NumBought = jsonData.Data.Contest.NumBought
			wishdataIncremental.NumCollection = jsonData.Data.Contest.NumEntered
			wishdataIncremental.RatingCount = int(jsonData.Data.Contest.ProductRating.RatingCount)
			wishdataIncremental.ProductId = util.FNV(jsonData.Data.Contest.ID)

			_, err := ini.AppWish.Insert(&wishdataIncremental)

			if err != nil {
				logger.Error(err.Error())
			}
		}
	}

	product.Updated = time.Now()
	if product.NumBought != jsonData.Data.Contest.NumBought ||
		product.NumEntered != jsonData.Data.Contest.NumEntered ||
		product.RatingCount != int(jsonData.Data.Contest.ProductRating.RatingCount) {

		configProduct(jsonData, &product)

		if _, err := ini.AppWish.Id(product.Id).Cols(
			"retail_price",
			"price",
			"shipping",
			"num_bought",
			"num_entered",
			"updated",
			"rating_count").Update(&product); err != nil {
			logger.Error(err.Error())
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

func nocacheWishId() (datas []string) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	page := <-pageChan
	var result []map[string][]byte
	var err error
	result, err = ini.AppWish.Query("select wish_id from wish_id limit ? offset ?", size, size*page)
	if err != nil {
		logger.Error(err.Error())
	}
	if len(result) <= 0 {
		pageChan <- 0
		if _, err = ini.RedisClient.HSet("load_page", "page", 1).Result(); err != nil {
			logger.Error(err.Error())
		}
		result, err = ini.AppWish.Query("select wish_id from wish_id limit ? offset ?", size, 0)

		if err != nil {
			logger.Error(err.Error())
		}
	} else {
		pageChan <- page + 1
	}
	for _, id := range result {
		datas = append(datas, string(id["wish_id"]))
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
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	cachePage := <-weekSalesPageChan
	var start = 0
	var end = 0
	if cachePage*size+size > global.WeekSalesCacheLenght {
		start = cachePage * size
		end = global.WeekSalesCacheLenght - cachePage*size
		global.WeekSalesCacheLenght = 0
	} else {
		start = cachePage * size
		end = cachePage*size + size
	}

	if ids, err := ini.RedisClient.
		LRange(global.WEEK_SALES_GREATER_THAN_ZERO, int64(start), int64(end)).
		Result(); err == nil {
		datas = ids
	} else {
		logger.Error(err.Error())
	}

	if len(datas) <= 0 {
		if ids, err := ini.RedisClient.
			LRange(global.WEEK_SALES_GREATER_THAN_ZERO, 0, int64(size)).
			Result(); err == nil {
			datas = ids
		} else {
			logger.Error(err.Error())
		}

		weekSalesPageChan <- 1
	} else {
		weekSalesPageChan <- cachePage + 1
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
