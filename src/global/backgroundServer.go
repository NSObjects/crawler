package global

import (
	"crawler/src/ini"
	"crawler/src/model"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/jinzhu/now"
)

var (
	WeekSalesCacheLenght            int
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
func CacheWeekSalesGreaterThanZeroWishId() {

	ini.RedisClient.Del(WEEK_SALES_GREATER_THAN_ZERO).Result()
	start := now.BeginningOfWeek()
	end := now.EndOfDay()
	results, err := ini.AppWish.Query("select DISTINCT product_id "+
		"from t_incremental "+
		"where product_id "+
		"in (select product_id "+
		"from t_incremental where created>=? and created<=? group by product_id"+
		" having sum(product_id)>0)", start, end)

	if err != nil || len(results) <= 0 {
		if err != nil {
			log.WithFields(logrus.Fields{
				"backgroundServer.go": "56",
			}).Error(err)
		}
		return
	}

	for index, r := range results {

		if index > 1000 {
			WeekSalesCacheLenght = len(results)
		}
		if productId, err := strconv.Atoi(string(r["product_id"])); err == nil {
			var product model.TProduct

			if _, err := ini.AppWish.Id(productId).Get(&product); err == nil {
				if product.WishId != "" {
					if err := ini.RedisClient.RPush(WEEK_SALES_GREATER_THAN_ZERO, product.WishId).Err(); err != nil {
						log.WithFields(logrus.Fields{
							"backgroundServer.go": "56",
						}).Error(err)
					}
				}
			} else {
				log.WithFields(logrus.Fields{
					"backgroundServer.go": "56",
				}).Error(err)
			}

		} else {
			log.WithFields(logrus.Fields{
				"backgroundServer.go": "56",
			}).Error(err)
		}
	}

}

func CacheWishId() {

	loadPage := &model.TLoadPage{Id: 1}
	_, err := ini.AppWish.Get(loadPage)
	if err != nil {
		log.WithFields(logrus.Fields{
			"backgroundServer.go": "79",
		}).Error(err)
	}

	page := loadPage.AllWishidCachePage
	for {
		AllWishIdCacheLenght, _ = ini.RedisClient.LLen(ALL_WISH_ID_CACHE).Result()
		//缓存内数据少于40000条时，新增缓存
		if AllWishIdCacheLenght < 40000 {
			results, err := ini.AppWish.Query("select wish_id from `t_wish_id` order by id limit 10000 offset ?", page*10000)

			if err != nil || len(results) == 0 {
				page = 0
				_, err := ini.AppWish.Exec("update t_load_page set all_wishid_cache_page=0")
				if err != nil {
					log.WithFields(logrus.Fields{
						"backgroundServer.go": "94",
					}).Error(err)
				}
				continue
			}

			for _, r := range results {
				err = ini.RedisClient.RPush(ALL_WISH_ID_CACHE, string(r["wish_id"])).Err()
				if err != nil {
					log.WithFields(logrus.Fields{
						"backgroundServer.go": "103",
					}).Error(err)
				}
			}
			page++

			_, err = ini.AppWish.Exec("update t_load_page set all_wishid_cache_page=?", page)
			if err != nil {
				log.WithFields(logrus.Fields{
					"backgroundServer.go": "111",
				}).Error(err)
			}
		}

	}

}

func CacheSalesGreaterThanWishId() {

	var loadPage model.TLoadPage
	_, err := ini.AppWish.Id(1).Get(&loadPage)

	if err != nil {
		log.WithFields(logrus.Fields{
			"backgroundServer.go": "130",
		}).Error(err)
		return
	}

	page := loadPage.SalesGtZeroPage

	for {
		SalesGreaterThanZeroCacheLenght, _ = ini.RedisClient.LLen(SALES_GREATER_THAN_ZERO).Result()
		if SalesGreaterThanZeroCacheLenght < 40000 {
			results, err := ini.AppWish.Query("select wish_id from t_product where num_bought > 0 order by id limit 10000 offset ?", page*10000)

			if err != nil || len(results) == 0 {
				page = 0
				_, err := ini.AppWish.Exec("update t_load_page set sales_gt_zero_page=0")
				if err != nil {
					log.WithFields(logrus.Fields{
						"backgroundServer.go": "146",
					}).Error(err)
				}
				continue
			}
			for _, r := range results {
				err = ini.RedisClient.RPush(SALES_GREATER_THAN_ZERO, r["wish_id"]).Err()
				if err != nil {
					log.WithFields(logrus.Fields{
						"backgroundServer.go": "156",
					}).Error(err)
				}
			}
			page++
			_, err = ini.AppWish.Exec("update t_load_page set sales_gt_zero_page=?", page)
			if err != nil {
				log.WithFields(logrus.Fields{
					"backgroundServer.go": "164",
				}).Error(err)
			}
		}

	}
}
