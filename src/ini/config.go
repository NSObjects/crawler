/*
 * Created  config.go on 17-12-4 下午3:42
 * Copyright (c) 2017  dyt.Co.Ltd All right reserved
 * Author lintao
 * Last modified 17-12-3 下午2:55
 */

package ini

import (
	"fmt"
	"time"

	"github.com/go-ini/ini"
)

var IniFile *ini.File

const IniPaht = "conf/config.ini"

func init() {
	local, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
	}
	time.Local = local

}

func Setup() {
	var err error

	IniFile, err = ini.Load(IniPaht)
	if err != nil {
		panic(err)
	}

	if err := initEngine(); err != nil {
		panic(err)
	}

	if err := initRedis(); err != nil {
		panic(err)
	}
}
