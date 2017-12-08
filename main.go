/*
 * Created  main.go on 17-12-4 下午3:51
 * Copyright (c) 2017  dyt.Co.Ltd All right reserved
 * Author lintao
 * Last modified 17-12-4 下午3:51
 */

package main

import (
	"crawler/src/controller"
	"crawler/src/global"
	"crawler/src/ini"
	"crawler/src/util"
	"log"
	_ "net/http/pprof"

	"crawler/src/router"
	"net/http"
	_ "net/http/pprof"

	"github.com/labstack/echo"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
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
	go global.CacheWeekSalesGreaterThanZeroWishId()

	util.LoopTimer(0, clearCache)
}

func clearCache() {
	ini.RedisClient.Del(global.SNAPSHOT_IDS)
}
