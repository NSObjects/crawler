package ini

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	AppWish *xorm.Engine
)

func initEngine() (err error) {
	mysql := new(mysql)

	err = IniFile.Section("mysql").MapTo(mysql)
	if err != nil {
		panic(err)
	}
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?&parseTime=True&loc=Local",
		mysql.User,
		mysql.Password,
		mysql.Host,
		mysql.Port,
		mysql.Dbname)

	AppWish, err = xorm.NewEngine("mysql", dns)
	AppWish.SetMaxIdleConns(mysql.MaxIdle)
	AppWish.SetMaxOpenConns(mysql.MaxConn)
	AppWish.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	//showSQL := ConfigFile.MustBool("xorm", "show_sql", false)
	//logLevel := ConfigFile.MustInt("xorm", "log_level", 1)
	//
	//MasterDB.ShowSQL(showSQL)
	//MasterDB.Logger().SetLevel(core.LogLevel(logLevel))

	// 启用缓存
	// cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	// MasterDB.SetDefaultCacher(cacher)

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
