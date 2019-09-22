package main

import (
	"fmt"

	"yokitalk.com/docservice/server/dbinstance"
	"yokitalk.com/docservice/server/repository"
)

func main() {
	mysqlManager := dbinstance.GetMysqlInstance()

	defer mysqlManager.Destroy()

	db := mysqlManager.DB

	userDBRepository := repository.NewUserRepository(db)

	user, err := userDBRepository.Find(6)

	fmt.Println(user)
	fmt.Println(err)
}