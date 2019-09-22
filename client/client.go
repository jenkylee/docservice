package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

const filePath  = "./cache/"

func main() {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	fileName := "test.docx"
	fileWriter, _ := bodyWriter.CreateFormFile("file", fileName)
	fmt.Println(filePath + fileName)
	file, _ := os.Open(filePath + fileName)
	defer file.Close()

	io.Copy(fileWriter, file)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, _ := http.Post("http://127.0.0.1:8080/upload", contentType, bodyBuffer)
	defer resp.Body.Close()

	resp_body, _ := ioutil.ReadAll(resp.Body)

	log.Println(resp.Status)
	log.Println(string(resp_body))
}
