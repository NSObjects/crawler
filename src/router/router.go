/*
 * Created  router.go on 17-12-4 下午3:42
 * Copyright (c) 2017  dyt.Co.Ltd All right reserved
 * Author lintao
 * Last modified 17-12-3 下午2:55
 */

package router

import (
	"crawler/src/controller"

	"github.com/labstack/echo"
)

func RegisterRoutes(g *echo.Group) {
	new(controller.ProductCrawlerController).RegisterRoute(g)
	new(controller.MerchantController).RegisterRoute(g)
	//new(controller.SnapshotController).RegisterRoute(g)
}
