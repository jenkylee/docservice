package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"yokitalk.com/docservice/server/util"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("cwd", cwd)
	dir := cwd + "/cache/"
	fileName := "1bffac319a0547a92d1799a9"
	imageDir := dir + "tex/" + fileName + "/image"
	wordXmlDir := dir + "tex/" + fileName + "/wordxml/"
	wordFile := dir + "upload/" + "1bffac319a0547a92d1799a9.docx"
	err = util.DeCompress(wordFile, wordXmlDir)
	if err != nil {
		log.Fatal(err)
	}
	wmfDir := wordXmlDir + "word/media"
	files, err := GetWmfFiles(wmfDir)

	var pngFiles []string

	for i := range files {
		name := filepath.Base(files[i])
		name = strings.Replace(name, ".wmf", "", 1)

		pngFile := imageDir + "/" + name + ".png"
		pngFiles = append(pngFiles, pngFile)
	}
	_, err = util.ExecCommand(time.Second * 300, "mogrify", "-density 100", "-units PixelsPerInch", "-background white", "-path "+imageDir, "-format png", wmfDir + "/*wmf")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pngFiles)
}

func GetWmfFiles(path string) (files []string, err error) {
	dirFiles, err := ioutil.ReadDir(path)

	if err != nil {
		return  nil, err
	}

	for _, fi := range  dirFiles{
		if fi.IsDir() {
			continue
		}

		ok := strings.HasSuffix(fi.Name(), ".wmf")

		if ok {
			files = append(files, path + "/" + fi.Name())
		}
	}

	return files, nil
}
