package dbinstance

import (
	"fmt"
	"log"
	"os"
	"sync"

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
	this.DB, this.ErrorMsg = gorm.Open("mysql", "root:123qwe@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local")

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