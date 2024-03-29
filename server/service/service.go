// Package service 定义 服务相关定义及函数
// word试题相关的服务函数定义

package service

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"yokitalk.com/docservice/server/model"
	"yokitalk.com/docservice/server/util"
)

// 提供DocService 操作
type DocService interface {
	Import(context.Context, string) (string, error)
	Export(context.Context, string) (string, error)
	Upload(context.Context, *http.Request) (string, error)
}

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

type docService struct {
	db   *gorm.DB
}

func NewDocService(db *gorm.DB) DocService {
	service := docService{}
	service.db = db

	return service
}

func (doc docService) Import(ctx context.Context, s string) (string, error){
	if s == "" {
		return "", ErrEmpty
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	//fmt.Println("cwd", cwd)
	dir := cwd + "/cache/"

	fileExt := strings.ToLower(path.Ext(s))
	_, file := path.Split(s)
	fileName := strings.ToLower(file);
	fileName = strings.Replace(fileName, fileExt, "", 1)
	fmt.Println(fileName)

	sFile := dir + "upload/" + s
	tFile := dir + "tex/" + fileName + "/" + strings.Replace(s, fileExt, ".tex", 1)
	imageDir := dir + "tex/" + fileName + "/image"
    log.Println(sFile)
	if !doc.isFileExist(sFile) {
		return "", errors.New("源文件不存在")
	}

	if !doc.isFileExist(tFile) {
		_, err := util.ExecCommand(time.Second * 300, "pandoc", sFile, "-o", tFile, "--extract-media=" + imageDir)
		if err != nil {
			return "", err
		}
	}

	//doc.osIoutil(tFile, doc.db)

	return tFile, nil
}

func (doc docService) Export(ctx context.Context, s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	//fmt.Println("cwd", cwd)
	dir := cwd + "/cache/"

	fileExt := strings.ToLower(path.Ext(s))
	_, file := path.Split(s)
	fileName := strings.ToLower(file);
	fileName = strings.Replace(fileName, fileExt, "", 1)
	fmt.Println(fileName)
	sFile := dir + "tex/" + fileName + "/" + s
	tFile := dir + "download/" + strings.Replace(s, fileExt, ".docx", 1)
	imageDir := dir + "tex//" + fileName + "/image"
	log.Println(sFile)
	if !doc.isFileExist(sFile) {
		return "", errors.New("源文件不存在")
	}

	if !doc.isFileExist(tFile) {
		// "--reference-doc="+dir+"moban.docx"
		_, err := util.ExecCommand(time.Second * 300, "pandoc",  sFile, "-o", tFile, "--extract-media=" + imageDir)
		if err != nil {
			return "", err
		}
	}
	return tFile, nil
}

func (doc docService) Upload(ctx context.Context, r *http.Request) (string, error) {

	file, handler, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		return "INVALID_FILE", err
	}

	if handler.Size > maxUploadSize {
		return "FILE_TOO_BIG", nil
	}

	fileName := handler.Filename

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "INVALID_FILE", err
	}

	// check file type, detectcontenttype only needs the first 512 bytes
	filetype := http.DetectContentType(fileBytes)
	fmt.Println("文件类型", filetype)
	/*switch filetype {
	case "image/jpeg", "image/jpg":
	case "image/gif", "image/png":
	case "application/pdf":
	case "application/zip":
		break
	default:
		return "INVALID_FILE_TYPE", nil
	}*/

	fileExt := strings.ToLower(path.Ext(fileName))
	if fileExt != ".docx" {
		return "INVALID_FILE_TYPE", nil
	}
	newFilName := util.RandToken(12)

	newPath := filepath.Join(uploadPath, newFilName+fileExt)
	fmt.Printf("FileType: %s, File: %s\n", fileExt, newPath)

	// 创建文件
	newFile, err := os.Create(newPath)
	if err != nil {
		return "CANT_WRITE_FILE", nil
	}
	defer newFile.Close() // idempotent, okay to call twice

	//将文件写到本地
	//_, err = io.Copy(newFile, file)
	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		return "CANT_WRITE_FILE", nil
	}

	return newFilName+fileExt, nil
}

func (docService) isFileExist(f string) bool {
	_, err := os.Stat(f)

	return err == nil || os.IsExist(err)
}

func (doc docService) handleWmfToPng(imgDir, wordFile, latexFile string) error {

	return nil
}

func (doc docService) handleImagToLatex(latexFile string) error {

	return  nil
}

func (doc docService) osIoutil(name string, db *gorm.DB) error {
	fileObj, err := os.Open(name);
	if  err == nil {
		return err
	}
	defer fileObj.Close()

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
			doc.questionParse(questionMap, tType, db)
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
				doc.questionParse(questionMap, oType, db)
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

	return nil
}

func (docService) questionParse(tMap map[int]string, tType string, db *gorm.DB) error {

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

		if idx := strings.Index(lineStr, tType); idx > -1 {
			isContent = true
			question.Title = lineStr[idx + len(tType): len(lineStr) -1]
			fmt.Println("Title", question.Title)
			continue
		}

		if strings.Index(lineStr, string("{[}答案{]}")) > -1 {
			tmpStr := strings.Replace(lineStr, string("{[}答案{]}"), "", 1);
			strArr := strings.Split(tmpStr, string("{[}分数{]}"))
			question.Answer = strArr[0]
			tmpStr = strArr[1]

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
				question.Subject = result[0][1];
				result = reg.FindAllStringSubmatch(strArr[1],-1) //匹配
				question.Difficulty = result[0][1];
			} else {
				result := reg.FindAllStringSubmatch(tmpStr,-1) //匹配
				question.Subject = result[0][1];
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
			question.Title += lineStr
		}
		if isAnalysis {
			question.Analysis += lineStr
		}
	}

	//question.ID = util.Md5str(question.Content)
	fmt.Println(question)
	//questionRepository := repository.NewQuestionRepository(db)
	//err := questionRepository.Create(&question)

	//return err
	return nil
}

var ErrEmpty = errors.New("empty strin")