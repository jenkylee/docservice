package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	fmt.Println("cwd", cwd)
	dir := cwd + "/cache/"
	sFile := dir + "test.docx"
	tFile := dir + "test.tex"

	if !isTexFileExist(sFile) {
		panic("文件不存在")
	}

	if !isTexFileExist(tFile) {
		_, err = execCommand(time.Second * 300, "pandoc", sFile, "-o", tFile, "--extract-media=cache/image")
		if err != nil {
			fmt.Printf("ret %v\n", err.Error())
		}
	}

	OsIoutil(tFile)
}

func isTexFileExist(f string) bool {
	_, err := os.Stat(f)

	return err == nil || os.IsExist(err)
}

func OsIoutil(name string) {

	var testingType map[string]string // 创建试题类型集合

	testingType = make(map[string]string)

	/* map插入对应试题类型及对应关键字， 便于分析tex试题文本*/
	testingType["single"] = "{[}单选题{]}"
	testingType["multiple"] = "{[}多选题{]}"
	testingType["indefinite"] = "{[}不定项选择题{]}"
	testingType["judgement"] = "{[}判断题{]}"
	testingType["completion"] = "{[}填空题{]}"
	testingType["questionanswer"] = "{[}问答题{]}"

	if fileObj, err := os.Open(name); err == nil {
		defer fileObj.Close()

		if err != nil {
			panic(err)
		}

		rd := bufio.NewReader(fileObj)
		var testingMap map[int]string
		testingMap = make(map[int]string)
		isTestStart := false
		isTestEnd := false
		i := 0
		t_l := 0
		for{
			line, err := rd.ReadString('\n')
			if err != nil || io.EOF == err {
				isTestStart = false
				isTestEnd = true
				TestingParse(testingMap)
				break
			} else {
				for _, v := range testingType {
					if strings.Index(line, v) > 0 {
						isTestStart = true
						if i > 0 {
							isTestEnd = true
						}
						i++
						break
					}
				}

				if isTestEnd {
					TestingParse(testingMap)
					isTestEnd = false
					testingMap = make(map[int]string)
					t_l = 0
				}

				if isTestStart {
					testingMap[t_l] = line
					t_l++
				}
			}
		}
	}
}

func TestingParse(tMap map[int]string)  {
	keys := []int{}

	for k := range tMap {
		keys = append(keys, k)
	}

	sort.Sort(sort.IntSlice(keys))

	for _, key := range keys  {
		fmt.Println(key, tMap[key])
	}
}