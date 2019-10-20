package util

import (
	"fmt"
	"testing"
)

func TestIsZip(t *testing.T) {
	f := "../../cache/test.docx"
	if isZip(f) {
		fmt.Println("ok")
	} else {
		fmt.Println("no")
	}
}

func TestZip(t *testing.T) {
	f := "../../cache/test.docx"
	unzip(f, "../../cache/test")
}