package c

import (
	"crawler/src/model"
	"errors"

	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func CrawlerWishId() {

	if w, err := requestMerchnt(); err == nil {
		if m, err := getWishIdFromMerchant(w.MerchntName, w.Users); err == nil {
			if m["wish_ids"] != nil {
				sendWishid(m)
				time.Sleep(30 * time.Second)
			}
		}
	} else {
		log.Error(err)
	}

	CrawlerWishId()
}

func requestMerchnt() (w WishIdJSON, err error) {

	client := &http.Client{}
	urlStr := fmt.Sprintf("http://%s/api/merchantCrawler", Host)
	req, err := http.NewRequest("GET", urlStr, nil)
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return w, err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(respBody, &w); err != nil {
		return w, err
	}

	return

}

func getWishIdFromMerchant(merchant string, user model.TUser) (m map[string]interface{}, err error) {

	if v, err := login(user); err == nil {
		user = v
	}
	m = make(map[string]interface{})
	if len(merchant) <= 0 {
		return m, errors.New("merchant name error")
	}
	invalidmerChantCount := 0
	page := 0

	var ids []string
	for {
		if merchant, err := loadMerchantProduct(page, merchant, user); err == nil {
			if len(merchant.Data.Results) <= 0 {
				break
			}
			for _, v := range merchant.Data.Results {
				ids = append(ids, v.ID)
			}
			page++
			if len(merchant.Data.MerchantInfo.Name) > 0 {
				m["merchant_info"] = merchant.Data.MerchantInfo
			}

		} else {
			if err.Error() == "Invalid merchant name" {
				invalidmerChantCount++
				if invalidmerChantCount < 100 {
					continue
				} else {
					break
				}
			} else {
				log.Error(err)
				break
			}
		}

	}

	m["wish_ids"] = ids
	return
}

func loadMerchantProduct(page int, merchant string, user model.TUser) (feeds *MerchantJSON, err error) {

	body := bodyWith(merchant, page, user)

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://www.wish.com/api/merchant", body)
	if err != nil {
		return nil, err
	}
	req = headerWish(req, user)
	// Fetch Request
	resp, err := client.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	var reader io.ReadCloser

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	default:
		reader = resp.Body
	}

	var b []byte
	buf := bytes.NewBuffer(b)
	buf.ReadFrom(reader)

	if err = json.Unmarshal(buf.Bytes(), &feeds); err != nil {
		return nil, err
	}

	if feeds.Code != 0 {
		if feeds.Msg == "Invalid merchant name" {
			return nil, errors.New(feeds.Msg)
		} else {
			return nil, fmt.Errorf("error :%s", fmt.Sprintln(feeds))
		}
	}

	if feeds.Msg == "Invalid merchant name" {
		return nil, errors.New(feeds.Msg)
	}

	return

}

func bodyWith(merchant string, page int, user model.TUser) *bytes.Buffer {
	params := url.Values{}

	params.Set("_app_type", "wish")
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

	//params.Set("advertiser_id", "FF22EA82-D474-4F28-9279-264E5F81946C")
	//.Set("app_device_model", "iPhone9,2")

	//params.Set("_riskified_session_token", "1BC70A62-009A-450E-9759-5C2DF3E032A7")
	//params.Set("app_device_id", "d5c712a51c5ef40cc341a5fcda73d5fc5b64de7d")
	//params.Set("_threat_metrix_session_token", "7625-67523522-9ED8-4B41-B56F-E3D2E7AC5D88")
	//body := bytes.NewBufferString(params.Encode())

	params.Set("_version", "3.20.6")
	params.Set("_client", "iosapp")
	params.Set("_xsrf", "1")
	params.Set("app_device_model", "iPhone9,2")

	params.Set("query", merchant)

	params.Set("start", fmt.Sprintf("%d", page*30))
	params.Set("count", "30")

	//params.Set("advertiser_id", user.AdvertiserId)
	//params.Set("_riskified_session_token", user.RiskifiedSessionToken)
	//params.Set("app_device_id", user.AppDeviceID)
	//AdvertiserID := strings.ToUpper(uuid.NewV4().String())
	//params.Set("_threat_metrix_session_token", AdvertiserID)
	body := bytes.NewBufferString(params.Encode())

	return body
}

func headerWish(req *http.Request, user model.TUser) *http.Request {
	// Headers
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Accept-Language", "zh-Hans-CN;q=1")
	cookie := fmt.Sprintf("_xsrf=1; _timezone=8; _appLocale=zh-Hans-CN; sweeper_session=\"%s\"; bsid=%s", user.SweeperSession, user.Baid)
	req.Header.Add("Cookie", cookie)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Wish/3.20.6 (iPhone; iOS 10.3.2; Scale/3.00)")

	return req
}

func sendWishid(wishIds map[string]interface{}) {

	jsonString, err := json.Marshal(wishIds)
	if err != nil {
		log.Error(err)
	}

	body := bytes.NewBuffer(jsonString)
	client := &http.Client{}
	urlStr := fmt.Sprintf("http://%s/api/merchantCrawler", Host)
	req, err := http.NewRequest("POST", urlStr, body)
	if err != nil {
		log.Error(err)
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
}

type MerchantJSON struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data Data   `json:"data"`
}

type Data struct {
	MerchantInfo MerchantInfo `json:"merchant_info"`
	Results      []struct {
		ID string `json:"id"`
	} `json:"results"`
	FeedEnded bool `json:"feed_ended"`
}

type MerchantInfo struct {
	PercentPositiveFeedback float64 `json:"percent_positive_feedback"`
	DisplayName             string  `json:"display_name"`
	Name                    string  `json:"name"`
	DisplayPic              string  `json:"display_pic"`
	AvgRating               float64 `json:"avg_rating"`
	ApprovedDate            float64 `json:"approved_date"`
	RatingCount             int     `json:"rating_count"`
}

type WishIdJSON struct {
	Users       model.TUser `json:"users"`
	MerchntName string      `json:"merchnt_name"`
}
