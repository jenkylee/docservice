package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"github.com/jinzhu/gorm"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"yokitalk.com/docservice/server/model"
	"yokitalk.com/docservice/server/repository"
)

var testingType map[string]string // 创建试题类型集合

func init()  {
	testingType = make(map[string]string)

	/* map插入对应试题类型及对应关键字， 便于分析tex试题文本*/
	testingType["single"] = "{[}单选题{]}"
	testingType["multiple"] = "{[}多选题{]}"
	testingType["indefinite"] = "{[}不定项选择题{]}"
	testingType["judgement"] = "{[}判断题{]}"
	testingType["completion"] = "{[}填空题{]}"
	testingType["questionanswer"] = "{[}问答题{]}"
}

func main() {
	/*cwd, err := os.Getwd()
	if err != nil {
		return
	}
	//fmt.Println("cwd", cwd)
	dir := cwd + "/cache/"
	s := "test.docx"
	fileExt := strings.ToLower(path.Ext(s))

	sFile := dir + "upload/" + s
	tFile := dir + "tex/" + strings.Replace(s, fileExt, ".tex", 1)
	imageDir := dir + "tex/image"

	if !isTexFileExist(sFile) {
		panic("文件不存在")
	}

	if !isTexFileExist(tFile) {
		_, err = execCommand(time.Second * 300, "pandoc", sFile, "-o", tFile, "--extract-media=" + imageDir)
		if err != nil {
			fmt.Printf("ret %v\n", err.Error())
		}
	}*/

/*	s := ""

	mysqlManager := dbinstance.GetMysqlInstance()

	defer mysqlManager.Destroy()

	db := mysqlManager.DB
*/
	//osIoutil(tFile, db)
}

func isTexFileExist(f string) bool {
	_, err := os.Stat(f)

	return err == nil || os.IsExist(err)
}

func osIoutil(name string, db *gorm.DB) {

	if fileObj, err := os.Open(name); err == nil {
		defer fileObj.Close()

		if err != nil {
			panic(err)
		}

		rd := bufio.NewReader(fileObj)
		var questionMap map[int]string
		questionMap = make(map[int]string)
		isTestStart := false
		isTestEnd := false
		i := 0
		t_l := 0
		tType := ""
		oType := ""
		for{
			line, err := rd.ReadString('\n')
			if err != nil || io.EOF == err {
				isTestStart = false
				isTestEnd = true
				questionParse(questionMap, tType, db)
				break
			} else {
				for k, v := range testingType {
					if strings.Index(line, v) > -1 {
						isTestStart = true
						if i > 0 {
							isTestEnd = true
						}
						i++
						tType = k
						if oType == "" {
							oType = tType
						}
						break
					}
				}

				if isTestEnd {
					questionParse(questionMap, oType, db)
					oType = tType
					isTestEnd = false
					questionMap = make(map[int]string)
					t_l = 0
				}

				if isTestStart {
					questionMap[t_l] = line
					t_l++
				}
			}
		}
	}
}

func questionParse(tMap map[int]string, tType string, db *gorm.DB) error {

	question := model.Question{}
	question.Type = tType

	tType = testingType[tType]

	keys := []int{}

	for k := range tMap {
		keys = append(keys, k)
	}

	sort.Sort(sort.IntSlice(keys))

	isContent := false
	isAnalysis := false
	for _, key := range keys  {
		lineStr := strings.Replace(tMap[key], "\n", "<br />", 1)

		if strings.Index(lineStr, tType) > -1 {
			isContent = true
			question.Content = strings.Replace(lineStr, tType, "", 1)
			continue
		}

		if strings.Index(lineStr, string("{[}答案{]}")) > -1 {
			tmpStr := strings.Replace(lineStr, string("{[}答案{]}"), "", 1);
			strArr := strings.Split(tmpStr, string("{[}分数{]}"))
			question.Answer = strArr[0]
			tmpStr = strArr[1]
			if strings.Index(tmpStr, string("{[}所有空无序{]}")) > -1 {
				tmpStr = strings.Replace(tmpStr, string("{[}所有空无序{]}"), "", 1)
				question.BlankDisorder = 1;
			} else {
				question.BlankDisorder = 0;
			}

			strArr = strings.Split(tmpStr, string("{[}分类{]}"))
			reg := regexp.MustCompile(`([0-9])`)

			result := reg.FindAllStringSubmatch(strArr[0],-1) //匹配

			mark, err := strconv.Atoi(result[0][1])
			if err == nil {
				question.Mark = mark
			}
			reg = regexp.MustCompile("^([a-zA-Z0-9\u4e00-\u9fa5,，]+).*")
			tmpStr = strArr[1]
			if strings.Index(tmpStr, string("{[}标签{]}")) > -1 {
				strArr = strings.Split(tmpStr, string("{[}标签{]}"))
				result = reg.FindAllStringSubmatch(strArr[0],-1) //匹配
				question.Class = result[0][1];
				result = reg.FindAllStringSubmatch(strArr[1],-1) //匹配
				question.Tag = result[0][1];
			} else {
				result := reg.FindAllStringSubmatch(tmpStr,-1) //匹配
				question.Class = result[0][1];
			}
			isContent = false
			continue
		}
		if strings.Index(lineStr, string("{[}解析{]}")) > -1 {
			question.Analysis = strings.Replace(lineStr, string("{[}解析{]}"), "", 1);
			isAnalysis = true
			isContent = false
			continue
		}
		if isContent {
			question.Content += lineStr
		}
		if isAnalysis {
			question.Analysis += lineStr
		}
	}

	question.ID = md5str(question.Content)
    //fmt.Println(question)
	questionRepository := repository.NewQuestionRepository(db)
	err := questionRepository.Create(&question)

	return err
}

func md5str(str string) string {
	h := md5.New()
	h.Write([]byte(str))

	return hex.EncodeToString(h.Sum(nil))
}