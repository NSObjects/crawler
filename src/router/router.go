package router

import (
	"crawler/src/controller"

	"github.com/labstack/echo"
)

func RegisterRoutes(g *echo.Group) {
	new(controller.ProductCrawlerController).RegisterRoute(g)
	new(controller.MerchantController).RegisterRoute(g)
	new(controller.SnapshotController).RegisterRoute(g)
}
