package model

import (
	"bytes"
	"crawler/src/ini"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/labstack/gommon/log"
	uuid "github.com/satori/go.uuid"
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

type LoginData struct {
	NewUser        bool   `json:"new_user"`
	AlreadyHadApp  bool   `json:"already_had_app"`
	User           string `json:"user"`
	SignupFlowType string `json:"signup_flow_type"`
	SessionToken   string `json:"session_token"`
}

type LoginInfo struct {
	Msg         string    `json:"msg"`
	Code        int       `json:"code"`
	Data        LoginData `json:"data"`
	SweeperUuid string    `json:"sweeper_uuid"`
	NotiCount   int       `json:"noti_count"`
}

func WishLoginWith(user *TUser) error {
	// Wish登录 (POST http://www.wish.com/api/email-login)

	body := requestBody(user.Password, user.Email)

	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://www.wish.com/api/email-login", body)
	if err != nil {
		return err
	}
	req = requestHeader(req)

	// Fetch Request
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()

	}
	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		return err
	}

	var loginInfo LoginInfo
	err = json.Unmarshal(respBody, &loginInfo)
	if err != nil {
		return err
	}

	if loginInfo.Code != 0 {
		return err
	}

	user.UserId = loginInfo.Data.User

	for _, cookie := range resp.Cookies() {

		switch cookie.Name {
		case "bsid":
			user.Baid = cookie.Value
		case "sweeper_session":
			user.SweeperSession = cookie.Value
		}

	}

	o := orm.NewOrm()
	if _, err := o.Update(user); err != nil {
		return err
	}

	return nil
}

func RegistIdWith() (user TUser) {
	// 注册 (POST https://www.wish.com/api/email-signup)

	firstNames := []string{
		"Aaron", "Abbott", "Abel", "Abner", "Abraham", "Adair", "Adam", "Addison",
		"Adolph", "Adonis", "Adrian", "Ahern", "Alan", "Albert", "Aldrich", "Alexander",
		"Alfred", "Alger", "Algernon", "Allen", "Alston", "Alva", "Alvin", "Alvis", "Amos", "Andre",
		"Andrew", "Andy", "Angelo", "Augus", "Ansel", "Antony", "Bevis", "Bill", "Bishop", "Blair", "Blake",
		"Bob", "Clarence", "Clark", "Claude", "Clyde", "Colin", "Dana", "Darnell", "Darcy", "Dempsey", "Dominic",
		"Edwiin", "Edward", "Elvis", "Fabian", "Frank", "Gale", "Gilbert", "Goddard", "Grover", "Hayden",
		"Hogan", "Hunter", "Isaac", "Ingram", "Isidore", "Jacob", "Jason", "Jay", "Jeff", "Jeremy", "Jesse",
		"Jerry", "Jim", "Jonathan", "Joseph", "Joshua", "Julian", "Julius", "Ken", "Kennedy", "Kent",
		"Kerr", "Kerwin", "Kevin", "Kirk", "King", "Lance", "Larry", "Leif", "Leonard", "Leopold", "Lewis",
		"Lionel", "Lucien", "Lyndon", "Magee", "Malcolm", "Mandel", "Marico", "Marsh", "Marvin", "Maximilian",
		"Meredith", "Merlin", "Mick", "Michell", "Monroe", "Montague", "Moore", "Mortimer", "Moses", "Nat",
		"Nathaniel", "Neil", "Nelson", "Newman", "Nicholas", "Nick", "Noah", "Noel", "Norton", "Ogden",
		"Oliver", "Omar", "Orville", "Osborn", "Oscar", "Osmond", "Oswald", "Otis", "Otto", "Owen", "Page", "Parker",
		"Paddy", "Patrick", "Paul", "Payne", "Perry", "Pete", "Peter", "Philip", "Phil",
		"Primo", "Quennel", "Quincy", "Quintion", "Rachel", "Ralap", "Randolph", "Robin", "Rodney", "Ron",
		"Roy", "Rupert", "Ryan", "Sampson", "Samuel", "Simon", "Stan", "Stanford", "Steward",
	}

	lastNames := []string{
		"Baker", "Hunter", "Carter", "Smith", "Cook", "Turner", "Baker", "Miller", "Smith", "Turner", "Hall",
		"Hill", "Lake", "Field", "Green", "Wood", "Well", "Brown", "Longman", "Short", "White", "Sharp",
		"Hard", "Yonng", "Sterling", "Hand", "Bull", "Fox", "Hawk", "Bush", "Stock", "Cotton", "Reed",
		"George", "Henry", "David", "Clinton", "Macadam", "Abbot", "Abraham", "Acheson", "Ackerman", "Adam",
		"Addison", "Adela", "Adolph", "Agnes", "Albert", "Alcott", "Aldington", "Alerander", "Alick", "Amelia",
		"Adams",
	}

	emailType := []string{
		"@qq.com", "@126.com", "@163.com", "@vip.sina.com", "@sina.com", "@tom.com", "@263.com", "@189.com",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	firstName := firstNames[r.Intn(len(firstNames))]
	lastName := lastNames[r.Intn(len(lastNames))]
	email := fmt.Sprintf("%d%s", time.Now().Unix(), emailType[r.Intn(len(emailType))])

	params := url.Values{}
	params.Set("_app_type", "wish")
	params.Set("_version", "3.20.0")
	params.Set("_client", "iosapp")

	params.Set("_xsrf", "1")
	params.Set("app_device_model", "iPhone9,2")
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
	params.Set("_capabilities[]", "47")
	params.Set("_capabilities[]", "6")
	params.Set("_capabilities[]", "7")
	params.Set("_capabilities[]", "8")
	params.Set("_capabilities[]", "9")

	AdvertiserID := strings.ToUpper(uuid.NewV4().String())
	params.Set("advertiser_id", AdvertiserID)

	riskifiedSessionToken := strings.ToUpper(uuid.NewV4().String())
	params.Set("_riskified_session_token", riskifiedSessionToken)

	key := []byte(strings.ToUpper(uuid.NewV4().String()))
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(time.Now().String()))
	appDeviceID := fmt.Sprintf("%x", mac.Sum(nil))
	params.Set("app_device_id", appDeviceID)

	params.Set("first_name", firstName)
	params.Set("last_name", lastName)
	params.Set("password", "1234567890")
	params.Set("email", email)
	body := bytes.NewBufferString(params.Encode())

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://www.wish.com/api/email-signup", body)

	// Headers
	req.Header.Add("Cookie", "_xsrf=1; _appLocale=zh-Hans-CN; _timezone=8")
	req.Header.Add("User-Agent", "Wish/3.20.0 (iPhone; iOS 10.3.1; Scale/3.00)")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Print(e)
		}
	}

	var loginInfo LoginInfo
	err = json.Unmarshal(respBody, &loginInfo)
	if err != nil {
		log.Print(err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Print(e)
		}
	}

	if loginInfo.Code != 0 {
		fmt.Println(loginInfo, email)
		return
		//utility.Errorln("登录错误", "Email:", user.Email, "Password:", user.Password, loginInfo)
	}

	user.UserId = loginInfo.Data.User
	user.Email = email
	user.Password = "1234567890"
	user.AppDeviceID = appDeviceID
	user.AdvertiserId = AdvertiserID
	user.RiskifiedSessionToken = riskifiedSessionToken
	user.FullName = firstName + " " + lastName
	for _, cookie := range resp.Cookies() {
		switch cookie.Name {
		case "bsid":
			user.Baid = cookie.Value
		case "sweeper_session":
			user.SweeperSession = cookie.Value
		}
	}
	return

}

func UpdateUser(user TUser, gender string) {
	// 更新用户信息 (POST https://www.wish.com/api/profile/update)

	params := url.Values{}
	params.Set("_app_type", "wish")
	params.Set("app_device_model", "iPhone9,2")
	params.Set("advertiser_id", "FF22EA82-D474-4F28-9279-264E5F81946C")
	params.Set("_client", "iosapp")
	params.Set("_xsrf", "1")
	params.Set("transform", "true")
	params.Set("dob_month", "5")
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
	params.Set("_capabilities[]", "47")
	params.Set("_capabilities[]", "6")
	params.Set("_capabilities[]", "7")
	params.Set("_capabilities[]", "8")
	params.Set("_capabilities[]", "9")
	params.Set("_version", "3.21.0")
	params.Set("_riskified_session_token", "B4B6612A-7A78-4460-86A0-68104A1E9369")
	params.Set("app_device_id", "d5c712a51c5ef40cc341a5fcda73d5fc5b64de7d")
	params.Set("dob_day", "28")
	params.Set("_threat_metrix_session_token", "7625-D0EB623D-4049-4AA8-B4B0-1444FCDA96E5")
	params.Set("gender", gender)
	params.Set("dob_year", "1992")
	body := bytes.NewBufferString(params.Encode())

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "https://www.wish.com/api/profile/update", body)

	// Headers
	cookie := fmt.Sprintf("_xsrf=1; _timezone=8; _appLocale=zh-Hans-CN; sweeper_session=\"%s\"; bsid=%s", user.SweeperSession, user.Baid)
	req.Header.Add("Cookie", cookie)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	o := orm.NewOrm()
	user.Gender = gender
	if _, err := o.Insert(&user); err != nil {
		log.Print(err)
	}

}

func requestHeader(req *http.Request) *http.Request {
	// Headers
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Cookie", "_xsrf=1; _timezone=8; _appLocale=zh-Hans-CN;")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "Wish/3.18.5 (iPhone; iOS 10.2.1; Scale/3.00)")
	req.Header.Add("Accept-Language", "zh-Hans-CN;q=1, en-CN;q=0.9")
	return req
}

func requestBody(password string, email string) *bytes.Buffer {
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
	params.Set("password", password)
	params.Set("_app_type", "wish")
	params.Set("_version", "3.18.2")
	params.Set("_xsrf", "1")
	params.Set("advertiser_id", "00000000-0000-0000-0000-000000000000")
	params.Set("email", email)
	params.Set("app_device_model", "iPhone7,1")
	params.Set("_client", "iosapp")
	params.Set("app_device_id", "34328efa59d113403ea83c649c68808e9988dc90")
	body := bytes.NewBufferString(params.Encode())
	return body
}
