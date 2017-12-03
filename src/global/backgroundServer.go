package global

import (
	"crawler/src/ini"
	"crawler/src/model"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/astaxie/beego/orm"
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
	o := orm.NewOrm()
	var list orm.ParamsList
	ini.RedisClient.Del(WEEK_SALES_GREATER_THAN_ZERO).Result()
	start := now.BeginningOfWeek()
	end := now.EndOfDay()
	_, err := o.Raw("select DISTINCT product_id "+
		"from t_incremental "+
		"where product_id "+
		"in (select product_id "+
		"from t_incremental where created>=? and created<=? group by product_id"+
		" having sum(product_id)>0)", start, end).ValuesFlat(&list)

	if err != nil || len(list) <= 0 {
		if err != nil {
			log.WithFields(logrus.Fields{
				"backgroundServer.go": "56",
			}).Error(err)
		}
		return
	}

	for index, pid := range list {

		if index > 1000 {
			WeekSalesCacheLenght = len(list)
		}

		if id, ok := pid.(string); ok == true {
			if pid, err := strconv.Atoi(id); err == nil {
				product := model.TProduct{Id: uint32(pid)}
				if err := o.Read(&product); err == nil {
					if err := ini.RedisClient.RPush(WEEK_SALES_GREATER_THAN_ZERO, product.WishId).Err(); err != nil {
						log.WithFields(logrus.Fields{
							"backgroundServer.go": "76",
						}).Error(err)
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
		lenght, _ := ini.RedisClient.LLen(ALL_WISH_ID_CACHE).Result()
		if lenght < 400000 {
			var list orm.ParamsList
			_, err := o.Raw("select wish_id from `wish_id` order by id limit 10000 offset ?", page*10000).ValuesFlat(&list)

			if err != nil || len(list) <= 0 {
				page = 0
				_, err := o.Raw("update load_page set all_wishid_cache_page=?", page).Exec()
				if err != nil {
					log.WithFields(logrus.Fields{
						"backgroundServer.go": "110",
					}).Error(err)
				}
				continue
			}

			for _, wishId := range list {
				if id, ok := wishId.(string); ok == true {
					if err := ini.RedisClient.RPush(ALL_WISH_ID_CACHE, id).Err(); err != nil {
						log.WithFields(logrus.Fields{
							"backgroundServer.go": "120",
						}).Error(err)
					}
				}
			}

			page++
			_, err = o.Raw("update load_page set all_wishid_cache_page=?", page).Exec()
			if err != nil {
				log.WithFields(logrus.Fields{
					"backgroundServer.go": "130",
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
		log.Fatal(err)
	}

	page := loadPage.SalesGtZeroPage

	for {
		lenght, _ := ini.RedisClient.LLen(SALES_GREATER_THAN_ZERO).Result()
		if lenght < 40000 {
			var list orm.ParamsList
			_, err := o.Raw("select wish_id from product where num_bought > 0 order by id limit 1000 offset ?", page*1000).ValuesFlat(&list)

			if err != nil || len(list) <= 0 {
				page = 0
				_, err := o.Raw("update load_page set sales_gt_zero_page=?", page).Exec()
				if err != nil {
					log.WithFields(logrus.Fields{
						"backgroundServer.go": "160",
					}).Error(err)
				}
				continue
			}

			for _, wishId := range list {
				if id, ok := wishId.(string); ok == true {
					if err := ini.RedisClient.RPush(SALES_GREATER_THAN_ZERO, id).Err(); err != nil {
						log.WithFields(logrus.Fields{
							"backgroundServer.go": "170",
						}).Error(err)
					}
				}
			}

			page++
			_, err = o.Raw("update load_page set sales_gt_zero_page=?", page).Exec()
			if err != nil {
				log.WithFields(logrus.Fields{
					"backgroundServer.go": "180",
				}).Error(err)
			}
		}
	}
}
