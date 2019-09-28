// Package service 定义 服务相关定义及函数
// 微服务的数据传输相关函数

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func DecodeAuthRequest(_ context.Context, r *http.Request) (interface{}, error)  {
	var request authRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func DecodeImportRequest(_ context.Context, r *http.Request) (interface{}, error)  {
	var request importRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func DecodeExportRequest(_ context.Context, r *http.Request) (interface{}, error)  {
	fmt.Println("test")
	var request exportRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func DecodeUploadRequest(_ context.Context, r *http.Request) (interface{}, error)  {
	var request uploadRequest
	request.R = r

	return request, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func EncodeDownloadResponse(_ context.Context, w http.ResponseWriter, response interface{}) error{
	return nil
}

func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter)  {
	code := http.StatusUnauthorized
	msg  := err.Error()

	w.WriteHeader(code)

	json.NewEncoder(w).Encode(serviceRespose{V: "", Err: msg})
}


func AuthErrorEncoder(_ context.Context, err error, w http.ResponseWriter)  {
	code := http.StatusUnauthorized
	msg  := err.Error()

	w.WriteHeader(code)

	json.NewEncoder(w).Encode(authResponse{Token: "", Err: msg})
}

type authRequest struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type importRequest struct {
	S string `json:"s"`
}

type exportRequest struct {
	S string `json:"s"`
}

type uploadRequest struct {
	R *http.Request
}

type authResponse struct {
	Token   string `json:"token,omitempty"`
	Err   string `json:"error,omitempty"`
} 

type importResponse struct {
	V string `json:"v"`
	Err string `json:"err, omitempty"`
}

type exportResponse struct {
	V string `json:"v"`
	Err string `json:"err, omitempty"`
}

type uploadResponse struct {
	V string `json:"v"`
	Err string `json:"err, omitempty"`
}

type serviceRespose struct {
	V string `json:"v"`
	Err string `json:"err, omitempty"`
}