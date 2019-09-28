// Package dbinstance 数据库初始化及单例模式实例化相关
package dbinstance

import (
	"fmt"
	"log"
	"os"
	"sync"
	"yokitalk.com/docservice/server/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var mysqlManager *MysqlManager
var mysqlOnce sync.Once

func GetMysqlInstance() *MysqlManager {
	mysqlOnce.Do(func() {
		mysqlManager = new(MysqlManager)
		mysqlManager.init()
	})

	return mysqlManager
}

type MysqlManager struct {
	DB *gorm.DB
	ErrorMsg error
}

func (this *MysqlManager) init() (*gorm.DB, error) {

	if err := config.Init(""); err != nil {
		panic(err)
	}
	conf := config.Config{}

	dbConf := conf.GetStringMapString("common.db")
	dbConn := dbConf["username"] + ":" + dbConf["password"] + "@tcp(" + dbConf["addr"] + ":" + dbConf["port"] + ")/" + dbConf["name"] + "?charset=utf8&parseTime=True&loc=Local"

	this.DB, this.ErrorMsg = gorm.Open(dbConf["driver"], dbConn)

	fmt.Println(this.ErrorMsg)

	if this.ErrorMsg != nil {
		log.Fatal(this.ErrorMsg)
	} else {
		//this.DB.SingularTable(true) // 表名单数
		this.DB.LogMode(true) // 输出日志
		this.DB.SetLogger(log.New(os.Stdout, "\r\n", 0)) // 输出到控制台
	}

	return this.DB, this.ErrorMsg
}

func (this *MysqlManager) Destroy()  {
	this.DB.Close()
}