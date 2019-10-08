package wkhtmltox

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
)


var (
	argsError     = errors.New("no input file or out path")
	fileTypeError = errors.New("the file must be in pdf format")
)

type HtmlToPdf struct {
	Command string
	in      string
	out     string
	argMap  map[string]string
	params  []string
}

func NewPdf() *HtmlToPdf {
	args := map[string]string{
		"--load-error-handling": "ignore",
		"--footer-center":       "第[page]页/共[topage]页",
		"--footer-font-size":    "8",
		"-B":                    "31",
		"-T":                    "32",
	}

	return &HtmlToPdf{
		Command: "wkhtmltopdf",
		argMap: args,
	}
}

func (this *HtmlToPdf) OutFile(input string, outPath string) (string, error) {
	// 输入输出参数不能为空
	if input == "" || outPath == "" {
		return outPath, argsError
	}

	// 判断是否生成pdf文件
	ext := filepath.Ext(outPath)
	if ext != ".pdf" {
		return outPath, fileTypeError
	}
	this.in = input
	this.out = outPath

	// 构建参数
	this.buildParams()
	// 执行命令
	bytes, err := this.doExec()
	if err != nil {
		return outPath, err
	}

	log.Println("[wkhtmltopdf-stdout]: %s", string(bytes))

	return outPath, nil
}

func (this *HtmlToPdf) doExec() ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx, this.Command, this.params...)
	stdout, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	defer stdout.Close()

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(bytes))
	fmt.Println("htmltopdf退出程序中", cmd.Process.Pid)
	cmd.Wait()

	return bytes, err
}

func (this *HtmlToPdf) buildParams() {
	for key, val := range this.argMap {
		this.params = append(this.params, key, val)
	}
	// 添加输入输出参数
	this.params = append(this.params, this.in, this.out)
}