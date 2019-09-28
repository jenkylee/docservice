// Package middlewares 微服务的中间见定义及函数
// 微服务的指标中间件函数

package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/metrics"

	"yokitalk.com/docservice/server/service"
)

type InstrumentingAuthMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	Next           service.AuthService
}

func (mw InstrumentingAuthMiddleware) Auth(clientID string, clientSecret string) (token string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Auth", "client", clientID, "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	token, err = mw.Next.Auth(clientID, clientSecret)
	return
}

type InstrumentingMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	CountResult    metrics.Histogram
	Next           service.DocService
}

func (mw InstrumentingMiddleware) Import(ctx context.Context, s string) (output string, err error) {
	defer func(begin time.Time) {
		custCl, _ := ctx.Value(jwt.JWTClaimsContextKey).(*service.CustomClaims)
		lvs := []string{"method", "import", "client", custCl.ClientID, "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.Next.Import(ctx, s)
	return
}

func (mw InstrumentingMiddleware) Export(ctx context.Context, s string) (output string, err error) {
	defer func(begin time.Time) {
		custCl, _ := ctx.Value(jwt.JWTClaimsContextKey).(*service.CustomClaims)
		lvs := []string{"method", "export", "client", custCl.ClientID, "error", "false"}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds()) 	//mw.CountResult.Observe(float64(n))
	}(time.Now())

	output, err = mw.Next.Export(ctx, s)
	return
}

func (mw InstrumentingMiddleware) Upload(ctx context.Context, r *http.Request) (output string, err error) {
	defer func(begin time.Time) {
		custCl, _ := ctx.Value(jwt.JWTClaimsContextKey).(*service.CustomClaims)
		lvs := []string{"method", "export", "client", custCl.ClientID, "error", "false"}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds()) 	//mw.CountResult.Observe(float64(n))
	}(time.Now())

	output, err = mw.Next.Upload(ctx, r)
	return
}