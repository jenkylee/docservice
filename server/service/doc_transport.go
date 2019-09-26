package service

import (
	"context"
	"encoding/json"
	"net/http"
	
	"github.com/go-kit/kit/endpoint"
)

func MakeAuthEndpoint(as AuthService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(authRequest)
		token, err := as.Auth(req.ClientID, req.ClientSecret)
		if err != nil {
			return nil, err
		}

		return authResponse{token, ""}, nil
	}
}

func MakeImportEndpoint(ds DocService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(importRequest)
		v, err := ds.Import(ctx, req.S)
		if err != nil {
			return importResponse{v, err.Error()}, nil
		}

		return importResponse{v, ""}, nil
	}
}

func MakeExportEndpoint(ds DocService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(exportRequest)
		v, err := ds.Export(ctx, req.S)
		if err != nil {
			return exportResponse{v, ""}, nil
		}

		return exportResponse{v, ""}, nil
	}
}

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
	var request exportRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
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