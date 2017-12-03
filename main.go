package main

import (
	"crawler/src/ini"
	"crawler/src/util"

	"crawler/src/global"

	"fmt"
	_ "net/http/pprof"
)

func main() {
	fmt.Println(util.FNV("59a539a6a02d3719e700eae2"))
	//ServeBackGround()
	//controller.Setup()
	//e := echo.New()
	//
	//g := e.Group("/api")
	//router.RegisterRoutes(g)
	//e.Logger.Fatal(e.Start(":2597"))
}

func init() {
	//ini.Setup()
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
