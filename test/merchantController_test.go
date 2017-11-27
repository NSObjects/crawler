package test

import (
	"crawler/src/controller"
	"crawler/src/model"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"

	"bytes"
	_ "crawler/src/ini"

	"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {

}

type mockMerchantData struct{}

func (this *mockMerchantData) GetMerchantName() (wishIdJSON model.MerchantNameJSON, err error) {
	wishIdJSON.MerchntName = "hello"
	wishIdJSON.Users = model.TUser{
		UserId:   "9527",
		Email:    "lt@pyl.com",
		Password: "123456",
	}
	return wishIdJSON, nil
}

func (this *mockMerchantData) MerchantInfoHandler(ctx echo.Context) error {
	var dat model.MerchantJSON
	b, err := ioutil.ReadAll(ctx.Request().Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &dat)
	if err != nil {
		return err
	}

	if len(dat.WishIds) <= 0 {
		return errors.New("wish id error")
	}

	if dat.WishIds[0] != "558a6a0e84c2e807b76f923d" {
		return errors.New("wish id is not equal 558a6a0e84c2e807b76f923d")
	}

	if dat.MerchantInfo.Name != "hello" {
		return errors.New("name is not equal hello")
	}

	if dat.MerchantInfo.DisplayName != "hello world" {
		return errors.New("display name is not equal hello world")
	}

	return nil
}

func TestMerchantController_GetMerchant(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)

	controller.Setup()
	Convey("/api/merchantCrawler", t, func() {
		context.SetPath("/api/merchantCrawler")
		this := controller.MerchantController{
			Data: new(mockMerchantData),
		}
		Convey("Get", func() {
			So(this.GetMerchant(context), ShouldBeNil)

			Convey("Get results", func() {
				So(rec.Code, ShouldEqual, http.StatusOK)
				var merchantName model.MerchantNameJSON
				json.Unmarshal(rec.Body.Bytes(), &merchantName)

				So(merchantName.MerchntName, ShouldEqual, "hello")
				So(merchantName.Users.Email, ShouldEqual, "lt@pyl.com")
				So(merchantName.Users.Password, ShouldEqual, "123456")
				So(merchantName.Users.UserId, ShouldEqual, "9527")
			})
		})
	})
}

func TestMerchantController_Post(t *testing.T) {

	e := echo.New()

	m := model.MerchantJSON{
		WishIds: []string{"558a6a0e84c2e807b76f923d"},
		MerchantInfo: model.MerchantInfo{
			Name:        "hello",
			DisplayName: "hello world",
		},
	}

	b, _ := json.Marshal(&m)
	body := bytes.NewBuffer(b)
	req := httptest.NewRequest(echo.POST, "/", body)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)

	Convey("/api/merchantCrawler", t, func() {
		context.SetPath("/api/merchantCrawler")
		this := controller.MerchantController{
			Data: new(mockMerchantData),
		}
		Convey("Post", func() {
			So(this.Post(context), ShouldBeNil)
			So(rec.Body.String(), ShouldEqual, "ok")
		})
	})
}
