package util

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"yokitalk.com/docservice/server/ocr"
)

// 执行系统命令函数
func ExecCommand(timeout time.Duration, name string, args ...string) (result []byte, err error) {

	cmd := exec.Command(name, args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}
    log.Println(cmd.String())
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	outBuf := bytes.NewBuffer(nil)
	errBuf := bytes.NewBuffer(nil)

	if err != nil {
		return
	}

	err = cmd.Start()

	if err != nil {
		return
	}

	stdin.Close()

	go io.Copy(outBuf, stdout)
	go io.Copy(errBuf, stderr)

	ch := make(chan error)

	go func(cmd *exec.Cmd) {
		defer close(ch)
		ch <- cmd.Wait()
	}(cmd)

	select {
	case err = <-ch:
	case <-time.After(timeout):
		_ = cmd.Process.Kill()
		err = errors.New("execute timeout")
		return
	}

	if err != nil {
		errStr := errBuf.String()
		return nil, errors.New(errStr)
	}

	if outBuf.Len() > 0 {
		return outBuf.Bytes(), nil
	}

	return
}

// md5加密函数
func Md5str(str string) string {
	h := md5.New()
	h.Write([]byte(str))

	return hex.EncodeToString(h.Sum(nil))
}

// 生成随机token函数
func RandToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// 外部调研
func ImageToLatex(image string) (string, error) {
	return mathpix(image)
}

// 调研mathpix将图片转为latex函数
func mathpix(image string) (string, error) {
	url := "https://api.mathpix.com/v3/latex"
	imageBytes, err := ioutil.ReadFile(image)
	if err != nil {
		log.Fatal("文件错误:", err)
	}
	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)
	data := strings.NewReader("{ \"src\": \"data:image/jpeg;base64,'" + imageBase64 + "'\", \"formats\": [\"latex_normal\"] }")
	req, _ := http.NewRequest("POST", url, data)
	req.Header.Add("app_id", "hxdo_qq_com")
	req.Header.Add("app_key", "7ca422fa18672a50189a")
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	result, _ := ioutil.ReadAll(res.Body)
	d := &ocr.Detecion{}
	err = json.Unmarshal(result, d)
	if err != nil {
		return "", err
	}

	return d.LatexNormal, nil
}
