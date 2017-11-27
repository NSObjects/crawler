package c

import (
	"bytes"
	"compress/gzip"
	"crawler/src/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

func crawlerProduct(user model.TUser) {
	for {
		getWishIdFromFeed("tabbed_feed_latest", user)
	}

}

func getWishIdFromFeed(categoryId string, user model.TUser) {

	page := 0
	for {
		if page > 1000 {
			return
		}
		if err := loadFeed(page, categoryId, user); err != nil {
			return
		}
		page++
	}

}

func loadFeed(page int, categoryId string, user model.TUser) error {

	body := feedBodyWith(page, user, categoryId)

	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://www.wish.com/api/feed/get-filtered-feed", body)
	if err != nil {
		return err
	}

	req = headerWish(req, user)

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Print(err)
		}

	default:
		reader = resp.Body
	}

	var b []byte
	buf := bytes.NewBuffer(b)
	buf.ReadFrom(reader)

	if resp.StatusCode != 200 {
		return fmt.Errorf("StatusCode: %d %s", resp.StatusCode, buf.Bytes())
	}

	var feeds Feed

	if err = json.Unmarshal(buf.Bytes(), &feeds); err != nil {
		return err
	}

	if feeds.Code == 10 {
		return fmt.Errorf("not more product")
	} else {

		if len(feeds.Data.Products) <= 0 {
			return fmt.Errorf("not more product")
		}

		for _, data := range feeds.Data.Products {

			if p, err := loadProductWith(data.ID, user); err == nil {

				for _, v := range p.Data.Contest.CommerceProductInfo.Variations {
					if len(v.MerchantName) > 0 {
						if m, err := getWishIdFromMerchant(v.MerchantName, user); err == nil {
							if m["wish_ids"] != nil {
								sendWishid(m)
							}
						}
						break
					}
				}
			}

		}

	}

	return nil
}

func feedBodyWith(page int, user model.TUser, category string) *bytes.Buffer {

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
	params.Set("request_id", category)
	params.Set("_app_type", "wish")
	params.Set("_version", "3.20.6")
	params.Set("_client", "iosapp")
	params.Set("_xsrf", "1")
	params.Set("app_device_model", "iPhone9,2")

	params.Set("advertiser_id", user.AdvertiserId)
	params.Set("_riskified_session_token", user.RiskifiedSessionToken)
	params.Set("app_device_id", user.AppDeviceID)
	//params.Set("_threat_metrix_session_token", user.)

	params.Set("count", "30")
	params.Set("offset", fmt.Sprintf("%d", page*30))

	body := bytes.NewBufferString(params.Encode())
	return body
}

func RandomUA() string {
	USER_AGENTS := []string{
		"Mozilla/5.0 (Windows; U; MSIE 9.0; Windows NT 9.0; en-US)",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Win64; x64; Trident/5.0; .NET CLR 3.5.30729; .NET CLR 3.0.30729; .NET CLR 2.0.50727; Media Center PC 6.0)",
		"Mozilla/5.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0; WOW64; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; .NET CLR 1.0.3705; .NET CLR 1.1.4322)",
		"Mozilla/4.0 (compatible; MSIE 7.0b; Windows NT 5.2; .NET CLR 1.1.4322; .NET CLR 2.0.50727; InfoPath.2; .NET CLR 3.0.04506.30)",
		"Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN) AppleWebKit/523.15 (KHTML, like Gecko, Safari/419.3) Arora/0.3 (Change: 287 c9dfb30)",
		"Mozilla/5.0 (X11; U; Linux; en-US) AppleWebKit/527+ (KHTML, like Gecko, Safari/419.3) Arora/0.6",
		"Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.8.1.2pre) Gecko/20070215 K-Ninja/2.1.1",
		"Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN; rv:1.9) Gecko/20080705 Firefox/3.0 Kapiko/3.0",
		"Mozilla/5.0 (X11; Linux i686; U;) Gecko/20070322 Kazehakase/0.4.5",
		"Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.0.8) Gecko Fedora/1.9.0.8-1.fc10 Kazehakase/0.5.6",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_3) AppleWebKit/535.20 (KHTML, like Gecko) Chrome/19.0.1036.7 Safari/535.20",
		"Mozilla/5.0 (Linux; Android 4.1.1; Nexus 7 Build/JRO03D) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.166  Safari/535.19",
		"Mozilla/5.0 (Linux; U; Android 4.0.4; en-gb; GT-I9300 Build/IMM76D) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
		"Mozilla/5.0 (Linux; U; Android 2.2; en-gb; GT-P1000 Build/FROYO) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
		"Mozilla/5.0 (Android; Mobile; rv:14.0) Gecko/14.0 Firefox/14.0",
		"Mozilla/5.0 (Android; Tablet; rv:14.0) Gecko/14.0 Firefox/14.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.8; rv:21.0) Gecko/20100101 Firefox/21.0",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:21.0) Gecko/20130331 Firefox/21.0",
		"Mozilla/5.0 (Windows NT 6.2; WOW64; rv:21.0) Gecko/20100101 Firefox/21.0",
		"Mozilla/5.0 (Linux; Android 4.0.4; Galaxy Nexus Build/IMM76B) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.133 Mobile Safari/535.19",
		"Mozilla/5.0 (Linux; Android 4.1.2; Nexus 7 Build/JZ054K) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.166 Safari/535.19",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.93 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/535.11 (KHTML, like Gecko) Ubuntu/11.10 Chromium/27.0.1453.93 Chrome/27.0.1453.93 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.94 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 6_1_4 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) CriOS/27.0.1453.10 Mobile/10B350 Safari/8536.25",
		"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0)",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)",
		"Mozilla/5.0 (compatible; WOW64; MSIE 10.0; Windows NT 6.2)",
		"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.9.168 Version/11.52",
		"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; fr) Presto/2.9.168 Version/11.52",
		"Opera/9.80 (Windows NT 6.1; WOW64; U; en) Presto/2.10.229 Version/11.62",
		"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_6; en-US) AppleWebKit/533.20.25 (KHTML, like Gecko) Version/5.0.4 Safari/533.20.27",
		"Mozilla/5.0 (iPad; CPU OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3",
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))

	return USER_AGENTS[r.Intn(len(USER_AGENTS))]

}

type Feed struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data struct {
		Products []struct {
			ID string `json:"id"`
		} `json:"products"`
	} `json:"data"`
}
