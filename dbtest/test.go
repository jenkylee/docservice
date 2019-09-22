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

	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.User{}, &model.Language{}, &model.Address{}, &model.CreditCard{}, &model.Email{})

	if err != nil {
		log.Fatal("fail")
	}
}