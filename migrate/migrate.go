package main


import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"yokitalk.com/docservice/server/model"
)

func main() {
	db, err := gorm.Open("mysql", "root:123qwe@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		log.Fatal("数据库连接失败")
	}
	if db.HasTable(&model.Question{}) { // 检查模型`User` 表是否存在
		log.Println("模型`Question` 表已存在")
	} else {
		// 为模型`User` 创建表
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&model.Question{})
		log.Println("模型`Question` 表创建成功")
		// db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.User{}, &model.Language{}, &model.Address{}, &model.CreditCard{}, &model.Email{})
	}
}