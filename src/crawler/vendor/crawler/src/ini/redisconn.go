package ini

import (
	"fmt"

	"github.com/go-ini/ini"
	redis "gopkg.in/redis.v5"
)

var RedisClient *redis.Client

func initRedis() (err error) {
	var sct *ini.Section
	sct, err = IniFile.GetSection("redis")
	if err != nil {
		panic(err)
	}
	host := sct.Key("host").String()
	port := sct.Key("port").String()

	addr := fmt.Sprintf("%s:%s", host, port)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       1,
	})

	_, err = RedisClient.Ping().Result()

	return
}
