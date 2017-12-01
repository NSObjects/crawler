package main

import (
	"crawler/src/controller"
	"crawler/src/ini"
	"crawler/src/router"
	"crawler/src/util"

	"github.com/labstack/echo"

	"crawler/src/global"

	_ "net/http/pprof"
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
