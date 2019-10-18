package main

import (
	"fmt"
	"github.com/guylaor/goword"
	"github.com/prometheus/common/log"
)

func main() {
	text, err := goword.Parse("../cache/upload/2532b6844d2e0add83a52fce.docx")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", text)
}