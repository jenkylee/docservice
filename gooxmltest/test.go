package main

import (
	"fmt"
	"log"
	"os"

	"baliance.com/gooxml/document"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	dir := cwd + "/cache/upload/"
	file := dir + "2532b6844d2e0add83a52fce.docx"
	doc, err := document.Open(file)
	if err != nil {
		log.Fatal("error opening document: %s", err)
	}

	for i, para := range doc.Paragraphs()  {
		fmt.Println("----第", i, "段--------")
		for _, run := range para.Runs() {
			fmt.Print(run.Text())
		}
		fmt.Println()
	}


	for ref := range doc.Images {
		fmt.Println(ref)
	}
}