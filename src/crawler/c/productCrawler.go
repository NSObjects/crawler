/*
 * Created  productCrawler.go on 17-12-4 下午3:42
 * Copyright (c) 2017  dyt.Co.Ltd All right reserved
 * Author lintao
 * Last modified 17-12-3 下午2:55
 */

package c

import (
	"bytes"
	"compress/gzip"
	"crawler/src/model"
	"crawler/src/util"
	"os"

	"github.com/Sirupsen/logrus"

	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const Host = "45.76.220.102:2597"

//const Host string = "localhost:2596"

var log = logrus.New()

func init() {
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
}

func CrawlerProduct() {
	if taskData, err := requestTaskData(); err == nil {
		proudcts := crawlerWishData(taskData)

		if len(proudcts) > 0 {
			go sendRequest(proudcts)
		}
	} else {
		log.WithFields(logrus.Fields{
			"productCrawler.go": "49",
		}).Error(err)
	}

	CrawlerProduct()
}

func FeedCrawler() {
	u := model.RegistIdWith()
	crawlerProduct(u)
}

func requestTaskData() (taskData TaskData, err error) {

	client := &http.Client{}
	urlStr := fmt.Sprintf("http://%s/api/wishdata", Host)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.WithFields(logrus.Fields{
			"productCrawler.go": "69",
		}).Error(err)
		return
	}

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return taskData, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return taskData, err
	}

	if err := json.Unmarshal(respBody, &taskData); err != nil {
		return taskData, err
	}

	return taskData, nil
}

func crawlerWishData(taskData TaskData) (proudcts []model.WishOrginalData) {

	var wg sync.WaitGroup
	p := util.New(12)
	for _, wishId := range taskData.WishIds {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		u := taskData.Users[r.Intn(len(taskData.Users))]

		wg.Add(1)
		go func(id string, user model.TUser) {
			p.Run(func() {
				if product, err := requestProductData(id, user); err == nil {
					proudcts = append(proudcts, product)
				} else {
					log.WithFields(logrus.Fields{
						"productCrawler.go": "109",
					}).Error(err)
				}
				wg.Done()
			})

		}(wishId, u)
	}

	wg.Wait()
	p.Shutdown()

	return
}

func sendRequest(p []model.WishOrginalData) (err error) {

	data := make(map[string]interface{})

	data["data"] = p
	body, err := json.Marshal(&data)
	if err != nil {
		log.WithFields(logrus.Fields{
			"productCrawler.go": "133",
		}).Error(err)
		return err
	}

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()

	_, err = w.Write(body)
	if err != nil {
		log.WithFields(logrus.Fields{
			"productCrawler.go": "145",
		}).Error(err)
		return err
	}
	err = w.Flush()
	if err != nil {
		log.WithFields(logrus.Fields{
			"productCrawler.go": "152",
		}).Error(err)
		return err
	}

	client := &http.Client{}
	urlStr := fmt.Sprintf("http://%s/api/wishdata", Host)

	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(b.Bytes()))
	if err != nil {
		log.WithFields(logrus.Fields{
			"productCrawler.go": "163",
		}).Error(err)
		return err
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.WithFields(logrus.Fields{
			"productCrawler.go": "175",
		}).Error(err)
		return err
	}
	fmt.Printf("发送数据%d条", len(p))
	//logger.Debug(fmt.Sprintf("发送数据%条", len(p)))
	return nil
}

func requestProductData(wishID string, user model.TUser) (wishPorduct model.WishOrginalData, err error) {
	if len(wishID) == 0 {
		return wishPorduct, errors.New("wish id error")
	}

	if product, err := loadProductWith(wishID, user); err == nil {

		if len(product.Data.Contest.Name) > 0 &&
			len(product.Data.Contest.ID) > 0 {
			if n, err := addToCart(product); err == nil {
				product.Data.Contest.NumBought = n
			} else {
				fmt.Println(err)
			}
			fmt.Println("num_bought:", product.Data.Contest.NumBought)
			wishPorduct = product
		}
	}

	return wishPorduct, nil
}

var user model.TUser

func addToCart(product model.WishOrginalData) (numBought int, err error) {
	// 加入购物车 (POST https://www.wish.com/api/cart/update)
	if len(user.Baid) <= 0 {
		user = model.RegistIdWith()
	}

	params := url.Values{}
	//params.Set("_riskified_session_token", "E1EEB61A-FFC3-40FA-B1F0-9165E217CB21")
	params.Set("_app_type", "wish")
	//params.Set("advertiser_id", "FF22EA82-D474-4F28-9279-264E5F81946C")
	params.Set("_xsrf", "1")
	//params.Set("_threat_metrix_session_token", "7625-7334AB75-1D44-4F5E-ACCB-0637B34E736C")
	params.Set("app_device_id", "d5c712a51c5ef40cc341a5fcda73d5fc5b64de7d")
	params.Set("app_device_model", "iPhone9,2")
	params.Set("_client", "iosapp")

	params.Set("_capabilities[]", "2")
	params.Set("_capabilities[]", "4")
	params.Set("_capabilities[]", "3")
	params.Set("_capabilities[]", "6")
	params.Set("_capabilities[]", "7")
	params.Set("_capabilities[]", "9")
	params.Set("_capabilities[]", "11")
	params.Set("_capabilities[]", "12")
	params.Set("_capabilities[]", "13")
	params.Set("_capabilities[]", "15")
	params.Set("_capabilities[]", "21")
	params.Set("_capabilities[]", "24")
	params.Set("_capabilities[]", "25")
	params.Set("_capabilities[]", "40")
	params.Set("_capabilities[]", "28")
	params.Set("_capabilities[]", "32")
	params.Set("_capabilities[]", "33")
	params.Set("_capabilities[]", "35")
	params.Set("_capabilities[]", "39")
	params.Set("_capabilities[]", "40")
	params.Set("_capabilities[]", "43")
	params.Set("_capabilities[]", "46")
	params.Set("_capabilities[]", "47")
	params.Set("_capabilities[]", "49")
	params.Set("_capabilities[]", "51")
	params.Set("_capabilities[]", "52")
	params.Set("_capabilities[]", "55")
	params.Set("_capabilities[]", "64")
	params.Set("_capabilities[]", "65")
	params.Set("_capabilities[]", "70")
	params.Set("_capabilities[]", "72")
	params.Set("_capabilities[]", "76")
	params.Set("_capabilities[]", "77")

	params.Set("add_to_cart", "true")
	params.Set("should_clear_cart", "true")

	for _, v := range product.Data.Contest.CommerceProductInfo.Variations {
		if len(v.VariationID) > 0 {
			params.Set("variation_id", v.VariationID)
			break
		}
	}

	params.Set("quantity", "1")

	params.Set("product_id", product.Data.Contest.ID)
	params.Set("shipping_option_id", "standard")
	params.Set("_version", "4.1.5")
	body := bytes.NewBufferString(params.Encode())

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "https://www.wish.com/api/cart/update", body)
	cookie := fmt.Sprintf("_xsrf=1; _timezone=8; _appLocale=zh-Hans-CN; sweeper_session=\"%s\"; bsid=%s", user.SweeperSession, user.Baid)
	req.Header.Add("Cookie", cookie)
	// Headers
	//req.Header.Add("Cookie", "sweeper_session=\"2|1:0|10:1512354599|15:sweeper_session|84:MzIyZDNiMjctNGEyMi00ODhkLWI0ODYtN2ExZjQxMWE4MGU4MjAxNy0xMi0wNCAwMjoyOTo1OS4yMTYzODI=|13219a254cb930574dd363c86bbf468e07a77fe270427575d9d42877130359ec\"; bsid=f250b6bcef394904b3a3db14896db65d; _xsrf=1; _appLocale=zh")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	// Display Results

	if resp.Status == "403 Forbidde" {
		fmt.Println("response Status : ", resp.Status)
		user = model.RegistIdWith()
	}

	var v Cart
	err = json.Unmarshal(respBody, &v)
	if err != nil {
		return numBought, err
	}

	if v.Code != 0 {
		return numBought, fmt.Errorf(v.Msg)
	}
	if len(v.Data.CartInfo.Items) > 0 {
		return v.Data.CartInfo.Items[0].NumBought, nil
	}
	return
}

type Cart struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data struct {
		CartInfo struct {
			Items []struct {
				VariationID string `json:"variation_id"`
				NumBought   int    `json:"num_bought"`
			} `json:"items"`
		} `json:"cart_info"`
	} `json:"data"`
}

func loadProductWith(wishID string, user model.TUser) (p model.WishOrginalData, e error) {

	var wishProduct model.WishOrginalData

	body := wbodyWish(wishID, user)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://www.wish.com/api/product/get", body)
	if err != nil {
		log.WithFields(logrus.Fields{
			"productCrawler.go": "210",
		}).Error(err)
		return wishProduct, err
	}
	req = wheaderWish(req, user)
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.WithFields(logrus.Fields{
			"productCrawler.go": "219",
		}).Error(err)
		return wishProduct, err
	}
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.WithFields(logrus.Fields{
				"productCrawler.go": "227",
			}).Error(err)
			return wishProduct, err
		}
	default:
		reader = resp.Body
	}
	var b []byte
	buf := bytes.NewBuffer(b)
	buf.ReadFrom(reader)
	if err = json.Unmarshal(buf.Bytes(), &wishProduct); err != nil {
		log.WithFields(logrus.Fields{
			"productCrawler.go": "241",
		}).Error(err)
		return wishProduct, err
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 405 {
			log.WithFields(logrus.Fields{
				"productCrawler.go": "244",
			}).Error(wishProduct.Msg)
		} else if wishProduct.Code != 12 && wishProduct.Code != 13 && wishProduct.Code != 11 {
			log.WithFields(logrus.Fields{
				"productCrawler.go": "248",
			}).Error(wishProduct.Msg)
		}
		return wishProduct, nil
	}

	return wishProduct, nil
}

func wheaderWish(req *http.Request, user model.TUser) *http.Request {
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Accept-Language", "zh-Hans-CN;q=1, en-CN;q=0.9")
	cookie := fmt.Sprintf("_xsrf=1; _timezone=8; _appLocale=zh-Hans-CN; sweeper_session=\"%s\"; bsid=%s", user.SweeperSession, user.Baid)
	req.Header.Add("Cookie", cookie)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Wish/3.18.5 (iPhone; iOS 10.2.1; Scale/3.00)")
	return req
}

func wbodyWish(wishId string, user model.TUser) *bytes.Buffer {
	params := url.Values{}
	params.Set("_capabilities[]", "11")
	params.Set("_capabilities[]", "12")
	params.Set("_capabilities[]", "13")
	params.Set("_capabilities[]", "15")
	params.Set("_capabilities[]", "2")
	params.Set("_capabilities[]", "21")
	params.Set("_capabilities[]", "24")
	params.Set("_capabilities[]", "25")
	params.Set("_capabilities[]", "28")
	params.Set("_capabilities[]", "32")
	params.Set("_capabilities[]", "35")
	params.Set("_capabilities[]", "39")
	params.Set("_capabilities[]", "4")
	params.Set("_capabilities[]", "40")
	params.Set("_capabilities[]", "43")
	params.Set("_capabilities[]", "6")
	params.Set("_capabilities[]", "7")
	params.Set("_capabilities[]", "8")
	params.Set("_capabilities[]", "9")
	params.Set("_app_type", "wish")
	params.Set("_version", "3.18.5")
	params.Set("_xsrf", "1")
	params.Set("app_device_model", "iPhone7,1")
	params.Set("_client", "iosapp")
	params.Set("advertiser_id", user.AdvertiserId)
	params.Set("app_device_id", user.AppDeviceID)
	params.Set("cid", wishId)
	body := bytes.NewBufferString(params.Encode())
	return body
}

func removeDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}
	for v := range elements {
		encountered[elements[v]] = true
	}
	result := []string{}
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

type TaskData struct {
	Message string        `json:"message"`
	Code    int           `json:"code"`
	WishIds []string      `json:"data"`
	Users   []model.TUser `json:"users"`
	Page    int           `json:"page"`
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
	Price                  float32   `json:"price"`
	RetailPrice            float32   `json:"retail_price"`
	MerchantTags           string    `json:"merchant_tags"`
	Tags                   string    `json:"tags"`
	Shipping               float32   `json:"shipping"`
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

type UserHeap []model.TUser

func (h UserHeap) Len() int           { return len(h) }
func (h UserHeap) Less(i, j int) bool { return h[i].Id < h[j].Id }
func (h UserHeap) Swap(i, j int)      { h[i].Id, h[j].Id = h[j].Id, h[i].Id }
func (h *UserHeap) Push(x interface{}) {
	*h = append(*h, x.(model.TUser))
}

func (h *UserHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func login(user model.TUser) (model.TUser, error) {
	// Wish登录 (POST http://www.wish.com/api/email-login)
	params := url.Values{}
	params.Set("email", user.Email)
	params.Set("password", user.Password)
	params.Set("_experiments", "")
	params.Set("_buckets", "")
	body := bytes.NewBufferString(params.Encode())
	// Create client
	client := &http.Client{}
	// Create request
	req, err := http.NewRequest("POST", "http://www.wish.com/api/email-login", body)
	// Headers
	req.Header.Add("Cookie", "_xsrf=C8B10FD5747D3B6B413A0F3F11422F55; IR_PI=1502072051896-f8zqdqa3znn; IR_EV=1502072051896|4953|0|1502072051896; __utmt=1; __utma=96128154.140752188.1502072052.1502072052.1502072052.1; __utmb=96128154.1.10.1502072052; __utmc=96128154; __utmz=96128154.1502072052.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); bsid=87e3bd6849794aaabb003290ae30cc6f; sweeper_uuid=77190f1ea92c4741aa11fb5dc4e07c79")
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("X-XSRFToken", "C8B10FD5747D3B6B413A0F3F11422F55")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8,en;q=0.6,zh-TW;q=0.4")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req)
	if err != nil {
		return user, err
	}
	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)
	for _, cookie := range resp.Cookies() {
		switch cookie.Name {
		case "bsid":
			user.Baid = cookie.Value
		case "sweeper_session":
			user.SweeperSession = cookie.Value
		}
	}
	if resp.StatusCode != 200 {
		return user, errors.New(string(respBody))
	}
	return user, nil
}
