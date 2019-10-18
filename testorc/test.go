package main

import (
	"fmt"
	"yokitalk.com/docservice/server/util"
)

func main() {
	file := "cache/test.png"
	latex, err := util.Mathpix(file)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(latex)
	}
}