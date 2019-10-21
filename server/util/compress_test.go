package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestIsZip(t *testing.T) {
	f := "../../cache/template.docx"
	if IsZip(f) {
		fmt.Println("ok")
	} else {
		fmt.Println("no")
	}
}

func TestCompress(t *testing.T) {
	f, err := os.Open("../../cache/test")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var files = []*os.File{f}

	dest := "../../cache/test2.zip"

	err = Compress(files,  dest)
	if err != nil {
		t.Fatal(err)
	}

	fileData, err := ioutil.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile("../../cache/test2.docx", fileData, 0755)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("ok")
}

func TestDeCompress(t *testing.T) {
	f := "../../cache/template.docx"
	err := DeCompress(f, "../../cache/test/")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("ok")
}