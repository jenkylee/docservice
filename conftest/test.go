package main

import (
	"fmt"
	"yokitalk.com/docservice/server/config"
)

func main() {
	if err := config.Init(""); err != nil {
		panic(err)
	}
	conf := config.Config{}

	dbConf := conf.GetStringMapString("common.db")

	fmt.Println(dbConf)
}
