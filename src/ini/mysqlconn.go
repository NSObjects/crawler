/*
 * Created  mysqlconn.go on 17-12-4 下午3:42
 * Copyright (c) 2017  dyt.Co.Ltd All right reserved
 * Author lintao
 * Last modified 17-12-4 下午2:22
 */

package ini

import (
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func initEngine() (err error) {
	local, err := time.LoadLocation("Asia/Shanghai")

	if err != nil {
		panic(err)
	}
	time.Local = local
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/source?charset=utf8&parseTime=true&loc=Asia%2FShanghai", 30, 30)

	return
}

type mysql struct {
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	User     string `ini:"user"`
	Password string `ini:"password"`
	Dbname   string `ini:"dbname"`
	MaxIdle  int    `ini:"max_idle"`
	MaxConn  int    `ini:"max_conn"`
}
