package service

import (
	"context"

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
