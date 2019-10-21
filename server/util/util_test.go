package util

import (
	"fmt"
	"testing"
)

func TestImageToLatex(t *testing.T) {
	file := "../../cache/WechatIMG36.jpeg"
	latex, err := ImageToLatex(file)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(latex)
	}
}
