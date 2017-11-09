package main

import (
	"crawlerController/src/router"
	"crawlerController/src/util"

	"crawlerController/src/global"

	"crawlerController/src/controller"

	"github.com/labstack/echo"
)

func main() {

	ServeBackGround()
	controller.Setup()
	e := echo.New()
	g := e.Group("/api")
	router.RegisterRoutes(g)
	e.Logger.Fatal(e.Start(":2596"))
}

func ServeBackGround() {

	go global.CacheWishId()
	go global.CacheSalesGreaterThanWishId()
	util.LoopTimer(9, global.CacheWeekSalesGreaterThanZeroWishId)
}
