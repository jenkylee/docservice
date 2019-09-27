package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"

	"yokitalk.com/docservice/server/service"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next service.DocService
}

func (mw LoggingMiddleware) Import(ctx context.Context, s string) (output string, err error) {
	defer func(begin time.Time) {
		//custCl, _ := ctx.Value(jwt.JWTClaimsContextKey).(*service.CustomClaims)
		_ = mw.Logger.Log(
				"method", "import",
				//"client", custCl.ClientID,
				"input", s,
				"output", output,
				"err", err,
				"took", time.Since(begin),
			)
	}(time.Now())

	output, err = mw.Next.Import(ctx, s)
	return
}

func (mw LoggingMiddleware) Export(ctx context.Context, s string) (output string, err error) {
	defer func(begin time.Time) {
		//custCl, _ := ctx.Value(jwt.JWTClaimsContextKey).(*service.CustomClaims)
		_ = mw.Logger.Log(
			"method", "export",
			//"client", custCl.ClientID,
			"input", s,
			"output", output,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.Next.Export(ctx, s)
	return
}

func (mw LoggingMiddleware) Upload(ctx context.Context, r *http.Request) (output string, err error) {
	defer func(begin time.Time) {
		//custCl, _ := ctx.Value(jwt.JWTClaimsContextKey).(*service.CustomClaims)
		_ = mw.Logger.Log(
			"method", "upload",
			//"client", custCl.ClientID,
			"input", "",
			"output", output,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.Next.Upload(ctx, r)
	return
}

type LoggingAuthMiddleware struct {
	Logger log.Logger
	Next service.AuthService
}

func (mw LoggingAuthMiddleware) Auth(clientID string, clientSecret string) (token string, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "auth",
			"clientID", clientID,
			"token", token,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	token, err = mw.Next.Auth(clientID, clientSecret)
	return
}

