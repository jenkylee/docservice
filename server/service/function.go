// Package service 定义 服务相关定义及函数
// 通用函数相关定义
package service

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"syscall"
	"time"
)

func execCommand(timeout time.Duration, name string, args ...string) (result []byte, err error) {

	cmd := exec.Command(name, args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}

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

func md5str(str string) string {
	h := md5.New()
	h.Write([]byte(str))

	return hex.EncodeToString(h.Sum(nil))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}