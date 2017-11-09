package controller

import (
	"crawler/src/global"
	"crawler/src/ini"
	"crawler/src/model"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"sync"

	"fmt"

	"crawler/src/util"
	"io/ioutil"
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
)

const size int = 250

var (
	u                 [][]model.User
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

	if requestCount >= 10 {
		weekSalesPage := <-weekSalesPageChan
		weekSalesPageChan <- weekSalesPage
		_, err := ini.AppWish.Exec("update load_page set week_sales_page=?", weekSalesPage)
		if err != nil {
			fmt.Println(err)
		}
		page := <-pageChan
		_, err = ini.AppWish.Exec("update load_page set all_id_page=?", page)
		if err != nil {
			fmt.Println(err)
		}
		pageChan <- page
		requestCount = 0
	}

	var datas []string
	mutex.Lock()
	if requestCount >= 5 && global.WeekSalesCacheLenght > 0 {
		fmt.Println(global.WeekSalesCacheLenght, "一周销量")
		datas = wishIdByWeekSalesGtZero()
	} else if requestCount >= 8 && global.SalesGreaterThanZeroCacheLenght > 0 {
		fmt.Println(global.SalesGreaterThanZeroCacheLenght, "销量大于0")
		datas = wishIdBySalesGtZero()
	} else if global.AllWishIdCacheLenght > 0 {
		fmt.Println(global.AllWishIdCacheLenght, "销量等于0")
		datas = allWishId()
	} else {
		datas = nocacheWishId()
	}
	mutex.Unlock()
	requestCount++
	if len(datas) > 0 {
		JSONData.Data = datas
		JSONData.Users = getUsers()
	}

	return ctx.JSON(http.StatusOK, JSONData)
}

func (this *ProductCrawlerController) Post(ctx echo.Context) error {

	b, _ := ioutil.ReadAll(ctx.Request().Body)

	if len(b) > 0 {

		SaveProductToDBFrom(b)
	}
	//
	//fmt.Println("外面")
	//utility.Statistics(0, this.Ctx.Input.IP())
	//this.Ctx.WriteString("Ok")
	return ctx.String(http.StatusOK, "ok")
}

func Setup() {
	loadPage := &model.LoadPage{Id: 1}
	_, err := ini.AppWish.Get(loadPage)
	if err != nil {
		panic(err)
	}

	weekSalesPageChan = make(chan int, 50)
	weekSalesPageChan <- loadPage.WeekSalesPage
	getUsers()
	pageChan = make(chan int, 50)
	pageChan <- loadPage.AllIdPage
}

func SaveProductToDBFrom(jsonStr []byte) {
	var w WishProductJSON

	err := json.Unmarshal(jsonStr, &w)
	if err != nil {
		fmt.Println(err)
	}

	for _, j := range w.Data {

		if j.Code != 0 || len(j.WishId) <= 0 {
			continue
		}

		var product model.Product
		if _, err := ini.AppWish.Id(util.FNV(j.WishId)).Get(&product); err == nil {
			saveWishDataIncremental(j, product)
		} else {
			if err == xorm.ErrNotExist {
				var p model.Product
				p.Id = util.FNV(j.WishId)
				configProduct(j, &p)
				if _, err := ini.AppWish.Insert(&p); err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println(err)
			}
		}

	}
}

type WishIdJson struct {
	Message string       `json:"message"`
	Code    int          `json:"code"`
	Data    []string     `json:"data"`
	Users   []model.User `json:"users"`
	Page    int64        `json:"page"`
}

type WishProductJSON struct {
	Data []WishProduct `json:"data"`
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

func saveWishDataIncremental(jsonData WishProduct, product model.Product) {
	if len(jsonData.Name) <= 0 || len(jsonData.WishId) <= 0 || jsonData.Code != 0 {
		return
	}

	if jsonData.CurrentlyViewing > 0 {
		v := model.Viewings{
			Count:     jsonData.CurrentlyViewing,
			ProductId: util.FNV(jsonData.WishId),
		}

		if jsonData.Time > 0 {
			v.Created = time.Unix(jsonData.Time, 0)
		} else {
			v.Created = time.Now()
		}

		if _, err := ini.AppWish.Insert(&v); err != nil {
			if !strings.Contains(err.Error(), "Duplicate entry") {
				fmt.Println(err)
			}
		}
	}

	wishdataIncremental := model.Incremental{}
	if time.Now().YearDay()-product.Updated.YearDay() <= 1 {
		if jsonData.NumBought > product.NumBought {
			wishdataIncremental.NumBoughtIncremental = jsonData.NumBought - product.NumBought
		}
		if jsonData.NumEntered != product.NumEntered {
			wishdataIncremental.NumCollectionIncremental = jsonData.NumEntered - product.NumEntered
		}
		if jsonData.RatingCount != product.RatingCount {
			wishdataIncremental.RatingCountIncremental = jsonData.RatingCount - product.RatingCount
		}
		if jsonData.Price != product.Price {
			wishdataIncremental.PriceIncremental = jsonData.Price - product.Price
		}
		if wishdataIncremental.NumBoughtIncremental > 0 ||
			wishdataIncremental.NumCollectionIncremental > 0 {
			wishdataIncremental.Price = jsonData.Price
			wishdataIncremental.Created = time.Now()
			wishdataIncremental.Updated = time.Now()
			wishdataIncremental.NumBought = jsonData.NumBought
			wishdataIncremental.NumCollection = jsonData.NumEntered
			wishdataIncremental.RatingCount = jsonData.RatingCount
			wishdataIncremental.ProductId = util.FNV(jsonData.WishId)

			_, err := ini.AppWish.Insert(&wishdataIncremental)

			if err != nil {
				if strings.Contains(err.Error(), "Error 1062: Duplicate entry") == false {
					fmt.Println(err)
				}
			}
		}
	}

	product.Updated = time.Now()
	if product.Price != jsonData.Price ||
		product.RetailPrice != jsonData.RetailPrice ||
		product.NumBought != jsonData.NumBought ||
		product.NumEntered != jsonData.NumEntered ||
		product.RatingCount != jsonData.RatingCount {

		configProduct(jsonData, &product)

		if _, err := ini.AppWish.Id(product.Id).Cols("rating_count",
			"retail_price",
			"shipping",
			"price",
			"num_bought",
			"num_entered",
			"true_tag_ids",
			"updated",
			"merchant",
			"merchant_id").Update(&product); err != nil {
			fmt.Println(err)
		}
	}

}

func configProduct(jsonData WishProduct, product *model.Product) {

	if len(jsonData.TrueTagIds) > 0 {
		product.TrueTagIds = strings.Join(jsonData.TrueTagIds, ",")
	}

	product.Gender = jsonData.Gender
	product.RatingCount = jsonData.RatingCount
	product.Price = jsonData.Price
	product.Size = jsonData.Size
	product.Color = jsonData.Color
	product.WishId = jsonData.WishId
	product.MaxShippingTime = jsonData.MaxShippingTime
	product.MinShippingTime = jsonData.MinShippingTime
	product.RetailPrice = jsonData.RetailPrice
	product.Merchant = jsonData.MerchantName
	product.MerchantId = util.FNV(product.Merchant)
	product.Shipping = jsonData.Shipping
	product.GenerationTime = jsonData.GenerationTime
	product.RetailPrice = jsonData.RetailPrice
	product.Tags = jsonData.Tags
	product.MerchantTags = jsonData.MerchantTags
	product.NumEntered = jsonData.NumEntered
	product.NumBought = jsonData.NumBought
	product.Name = jsonData.Name
	product.Description = jsonData.Description
	product.Updated = time.Now()

}

func getUsers() []model.User {

	if len(u) <= 0 {

		contrys := []string{"Britain", "Canada", "Australia", "France", "Germany", "America"}
		for _, contry := range contrys {
			var user []model.User
			err := ini.AppWish.Where("has_address=1").And("country=?", contry).Find(&user)
			if err != nil {
				fmt.Println(err)
			}
			u = append(u, user)

		}
	}
	var users []model.User

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

func nocacheWishId() (datas []string) {

	page := <-pageChan
	var result []map[string][]byte
	var err error
	result, err = ini.AppWish.Query("select wish_id from wish_id limit ? offset ?", size, size*page)
	if err != nil {
		fmt.Println(err)
	}
	if len(result) <= 0 {
		pageChan <- 0
		if _, err = ini.RedisClient.HSet("load_page", "page", 1).Result(); err != nil {
			fmt.Println(err)
		}
		result, err = ini.AppWish.Query("select wish_id from wish_id limit ? offset ?", size, 0)

		if err != nil {
			fmt.Println(err)
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
		fmt.Println(err)
	}

	if len(datas) <= 0 {
		if ids, err := ini.RedisClient.
			LRange(global.WEEK_SALES_GREATER_THAN_ZERO, 0, int64(size)).
			Result(); err == nil {
			datas = ids
		} else {
			fmt.Println(err)
		}

		weekSalesPageChan <- 1
	} else {
		weekSalesPageChan <- cachePage + 1
	}
	return datas
}
