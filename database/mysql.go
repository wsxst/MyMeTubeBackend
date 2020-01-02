package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
	"github.com/Unknwon/goconfig"
)

// ORM gorm引擎的实例
var ORM *gorm.DB

func Init() {
	cfg, err := goconfig.LoadConfigFile("mysql.ini")
	if err != nil{
		panic("错误")
	}
	sec, err := cfg.GetSection("mysql")
	var (
		userName = sec["userName"]
		password = sec["password"]
		ip = sec["ip"]
		port = sec["port"]
		dbName = sec["dbName"]
	)
	path := strings.Join([]string{userName, ":", password, "@tcp(",ip, ":", port, ")/", dbName, "?charset=utf8&parseTime=True&loc=Local"}, "")

	//连接数据库
	db, err := gorm.Open("mysql", path)
	if err != nil {
		fmt.Println(err)
		panic("连接数据库失败")
	}

	ORM = db
	// 打印sql详情
	db.LogMode(true)
}
