/*
 * Created  backgroundServer.go on 17-12-4 下午3:42
 * Copyright (c) 2017  dyt.Co.Ltd All right reserved
 * Author lintao
 * Last modified 17-12-4 下午2:50
 */

package global

import (
	"crawler/src/ini"
	"crawler/src/model"
	"crawler/src/util"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/astaxie/beego/orm"
)

var (
	WeekSalesCacheLenght            int64
	AllWishIdCacheLenght            int64
	SalesGreaterThanZeroCacheLenght int64
)

var log = logrus.New()

func init() {
	file, err := os.OpenFile("err.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	local, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
	}
	time.Local = local
}

/*
	缓存一周销量大于0的WishId
*/

var list orm.ParamsList

func cacheList() {
	o := orm.NewOrm()

	ini.RedisClient.Del(WEEK_SALES_GREATER_THAN_ZERO).Result()
	_, err := o.Raw("select DISTINCT product_id from t_incremental where product_id in (select product_id from t_incremental group by product_id having sum(product_id)>0)").ValuesFlat(&list)

	if err != nil {
		log.WithFields(logrus.Fields{
			"backgroundServer.go": "61",
		}).Error(err)
	}
}

func CacheWeekSalesGreaterThanZeroWishId() {
	o := orm.NewOrm()
	util.LoopTimer(1, cacheList)
	for {
		WeekSalesCacheLenght, _ = ini.RedisClient.LLen(WEEK_SALES_GREATER_THAN_ZERO).Result()

		if WeekSalesCacheLenght < int64(len(list)/2) {
			for _, pid := range list {
				if id, ok := pid.(string); ok == true {
					if pid, err := strconv.Atoi(id); err == nil {
						product := model.TProduct{Id: uint32(pid)}
						if err := o.Read(&product); err == nil {
							if err := ini.RedisClient.RPush(WEEK_SALES_GREATER_THAN_ZERO, product.WishId).Err(); err != nil {
								log.WithFields(logrus.Fields{
									"backgroundServer.go": "82",
								}).Error(err)
							}
						}
					}

				}

			}
		}
	}

}

func CacheWishId() {

	o := orm.NewOrm()
	loadPage := model.TLoadPage{Id: 1}
	err := o.Read(&loadPage)
	if err != nil {
		log.Fatal(err)
	}

	page := loadPage.AllWishIdCachePage

	for {
		AllWishIdCacheLenght, _ = ini.RedisClient.LLen(ALL_WISH_ID_CACHE).Result()
		if AllWishIdCacheLenght < 400000 {
			var list orm.ParamsList
			_, err := o.Raw("select wish_id from `t_wish_id` order by id limit 10000 offset ?", page*10000).ValuesFlat(&list)

			if err != nil || len(list) <= 0 {
				page = 0
				_, err := o.Raw("update t_load_page set all_wishid_cache_page=?", page).Exec()
				if err != nil {
					log.WithFields(logrus.Fields{
						"backgroundServer.go": "118",
					}).Error(err)
				}
				continue
			}

			for _, wishId := range list {
				if id, ok := wishId.(string); ok == true {
					if err := ini.RedisClient.RPush(ALL_WISH_ID_CACHE, id).Err(); err != nil {
						log.WithFields(logrus.Fields{
							"backgroundServer.go": "128",
						}).Error(err)
					}
				}
			}

			page++
			_, err = o.Raw("update t_load_page set all_wishid_cache_page=?", page).Exec()
			if err != nil {
				log.WithFields(logrus.Fields{
					"backgroundServer.go": "138",
				}).Error(err)
			}
		}
	}

}

func CacheSalesGreaterThanWishId() {

	o := orm.NewOrm()
	loadPage := model.TLoadPage{Id: 1}
	err := o.Read(&loadPage)
	if err != nil {
		log.WithFields(logrus.Fields{
			"backgroundServer.go": "153",
		}).Error(err)
	}

	page := loadPage.SalesGtZeroPage

	for {
		SalesGreaterThanZeroCacheLenght, _ = ini.RedisClient.LLen(SALES_GREATER_THAN_ZERO).Result()
		if SalesGreaterThanZeroCacheLenght < 400000 {
			var list orm.ParamsList
			_, err := o.Raw("select wish_id from t_product where num_bought > 0 order by id limit 10000 offset ?", page*10000).ValuesFlat(&list)

			if err != nil || len(list) <= 0 {
				page = 0
				_, err := o.Raw("update t_load_page set sales_gt_zero_page=?", page).Exec()
				if err != nil {
					log.WithFields(logrus.Fields{
						"backgroundServer.go": "170",
					}).Error(err)
				}
				continue
			}

			for _, wishId := range list {
				if id, ok := wishId.(string); ok == true {
					if err := ini.RedisClient.RPush(SALES_GREATER_THAN_ZERO, id).Err(); err != nil {
						log.WithFields(logrus.Fields{
							"backgroundServer.go": "180",
						}).Error(err)
					}
				}
			}

			page++
			_, err = o.Raw("update t_load_page set sales_gt_zero_page=?", page).Exec()
			if err != nil {
				log.WithFields(logrus.Fields{
					"backgroundServer.go": "190",
				}).Error(err)
			}
		}
	}
}
