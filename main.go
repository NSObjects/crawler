package main

import (
	"crawler/src/ini"
	"crawler/src/router"
	"crawler/src/util"

	_ "github.com/sevenNt/echo-pprof"

	"crawler/src/global"

	"crawler/src/controller"
	_ "net/http/pprof"

	"github.com/labstack/echo"
)

func main() {
	ServeBackGround()
	controller.Setup()
	e := echo.New()
	util.Wrap(e)
	g := e.Group("/api")
	router.RegisterRoutes(g)
	e.Logger.Fatal(e.Start(":2596"))
}

func init() {
	ini.Setup()
}

func ServeBackGround() {
	go global.CacheWishId()
	go global.CacheSalesGreaterThanWishId()
	util.LoopTimer(9, global.CacheWeekSalesGreaterThanZeroWishId)
}
