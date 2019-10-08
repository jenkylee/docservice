package main

import (
	"fmt"
	"log"
	"os"

	"yokitalk.com/docservice/server/wkhtmltox"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	dir := cwd + "/cache/"
	in := "http://blog.sina.com.cn/s/blog_4660f1870101hwbk.html"
	out := dir + "test.pdf"
	pdf := wkhtmltox.NewPdf()
	path, err := pdf.OutFile(in, out)

	if err != nil {
		log.Println(err)
	}

	fmt.Println(path)
}