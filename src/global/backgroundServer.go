package global

import (
	"crawler/src/ini"
	"crawler/src/model"
	"fmt"
	"strconv"

	"github.com/jinzhu/now"
)

var (
	WeekSalesCacheLenght            int
	AllWishIdCacheLenght            int64
	SalesGreaterThanZeroCacheLenght int64
)

/*
	缓存一周销量大于0的WishId
*/
func CacheWeekSalesGreaterThanZeroWishId() {

	ini.RedisClient.Del(WEEK_SALES_GREATER_THAN_ZERO).Result()
	start := now.BeginningOfWeek()
	end := now.EndOfDay()
	results, err := ini.AppWish.Query("select DISTINCT product_id "+
		"from incremental "+
		"where product_id "+
		"in (select product_id "+
		"from incremental where created>=? and created<=? group by product_id"+
		" having sum(product_id)>0)", start, end)

	if err != nil {
		panic(err)
	}
	for _, r := range results {
		if productId, err := strconv.Atoi(string(r["product_id"])); err == nil {
			var product model.Product
			if _, err := ini.AppWish.Id(productId).Get(&product); err == nil {
				if product.WishId != "" {
					if err := ini.RedisClient.RPush(WEEK_SALES_GREATER_THAN_ZERO, product.WishId).Err(); err != nil {
						fmt.Println(err)
					}
				}

			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	}

	WeekSalesCacheLenght = len(results)

}

func CacheWishId() {
	loadPage := &model.LoadPage{Id: 1}
	_, err := ini.AppWish.Get(loadPage)
	if err != nil {
		panic(err)
	}

	page := loadPage.AllWishidCachePage
	for {
		AllWishIdCacheLenght, _ = ini.RedisClient.LLen(ALL_WISH_ID_CACHE).Result()
		if AllWishIdCacheLenght < 400000 {
			results, err := ini.AppWish.Query("select wish_id from `wish_id` order by id limit 10000 offset ?", page*10000)
			if err != nil || len(results) == 0 {
				page = 0
				_, err := ini.AppWish.Exec("update load_page set all_wishid_cache_page=0")
				if err != nil {
					panic(err)
				}
				continue
			}
			for _, r := range results {
				err = ini.RedisClient.RPush(ALL_WISH_ID_CACHE, string(r["wish_id"])).Err()
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		page++
		_, err := ini.AppWish.Exec("update load_page set all_wishid_cache_page=?", page)
		if err != nil {
			panic(err)
		}
	}

}

func CacheSalesGreaterThanWishId() {
	var loadPage model.LoadPage
	_, err := ini.AppWish.Id(1).Get(&loadPage)
	if err != nil {
		panic(err)
	}
	page := loadPage.SalesGtZeroPage
	for {
		SalesGreaterThanZeroCacheLenght, _ = ini.RedisClient.LLen(SALES_GREATER_THAN_ZERO).Result()
		if SalesGreaterThanZeroCacheLenght < 400000 {
			results, err := ini.AppWish.Query("select wish_id from product where num_bought > 0 order by id limit 10000 offset ?", page*10000)
			if err != nil || len(results) == 0 {
				page = 0
				_, err := ini.AppWish.Exec("update load_page set sales_gt_zero_page=0")
				if err != nil {
					panic(err)
				}
				continue
			}
			for _, r := range results {
				err = ini.RedisClient.RPush(SALES_GREATER_THAN_ZERO, r["wish_id"]).Err()
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		page++
		_, err := ini.AppWish.Exec("update load_page set sales_gt_zero_page=?", page)
		if err != nil {
			panic(err)
		}
	}
}
