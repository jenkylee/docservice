package service

import (
	"context"
	"encoding/json"
	"net/http"
	
	"github.com/go-kit/kit/endpoint"
)

func MakeImportEndpoint(ds Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(importRequest)
		v, err := ds.Import(req.S)
		if err != nil {
			return importResponse{v, err.Error()}, nil
		}

		return importResponse{v, ""}, nil
	}
}

func MakeExportEndpoint(ds Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(exportRequest)
		v := ds.Export(req.S)
		if err != nil {
			return exportResponse{v}, nil
		}

		return exportResponse{v}, nil
	}
}

func DecodeImportRequest(_ context.Context, r *http.Request) (interface{}, error)  {
	var request importRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func DecodeExportRequest(_ context.Context, r *http.Request) (interface{}, error)  {
	var request exportRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type importRequest struct {
	S string `json:"s"`
}

type importResponse struct {
	V string `json:"v"`
	Err string `json:"err, omitempty"`
}

type exportRequest struct {
	S string `json:"s"`
}

type exportResponse struct {
	V int `json:"v"`
} 