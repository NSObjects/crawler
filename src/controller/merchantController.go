package controller

import (
	"crawler/src/model"
	"net/http"

	"github.com/labstack/echo"
)

type MerchantController struct {
	Data model.MerchanApiInterface
}

func (this *MerchantController) RegisterRoute(g *echo.Group) {
	g.GET("/merchantCrawler", this.GetMerchant)
	g.POST("/merchantCrawler", this.Post)
	this.Data = new(model.Merchant)
}

func (this *MerchantController) GetMerchant(ctx echo.Context) error {
	value, err := this.Data.GetMerchantName()
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, value)
}

func (this *MerchantController) Post(ctx echo.Context) error {
	if err := this.Data.MerchantInfoHandler(ctx); err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	return ctx.String(http.StatusOK, "ok")
}
