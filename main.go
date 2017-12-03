package main

import (
	"crawler/src/ini"
	"crawler/src/util"

	"crawler/src/global"

	"crawler/src/controller"
	"crawler/src/router"
	_ "net/http/pprof"

	"github.com/labstack/echo"
)

func main() {

	ServeBackGround()
	controller.Setup()
	e := echo.New()

	g := e.Group("/api")
	router.RegisterRoutes(g)
	e.Logger.Fatal(e.Start(":2597"))
}

func init() {
	ini.Setup()
}

func ServeBackGround() {
	go global.CacheWishId()
	go global.CacheSalesGreaterThanWishId()
	util.LoopTimer(9, global.CacheWeekSalesGreaterThanZeroWishId)
	util.LoopTimer(0, clearCache)
}

func clearCache() {
	ini.RedisClient.Del(global.SNAPSHOT_IDS)
}
