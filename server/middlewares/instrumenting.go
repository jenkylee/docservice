package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
		lvs := []string{"method", "Auth", "error", fmt.Sprint(err != nil)}
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
		//custCl, _ := ctx.Value(jwt.JWTClaimsContextKey).(*service.CustomClaims)
		//lvs := []string{"method", "import", "client", custCl.ClientID, "error", fmt.Sprint(err != nil)}
		lvs := []string{"method", "import", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.Next.Import(ctx, s)
	return
}

func (mw InstrumentingMiddleware) Export(ctx context.Context, s string) (output string, err error) {
	defer func(begin time.Time) {
		//custCl, _ := ctx.Value(jwt.JWTClaimsContextKey).(*service.CustomClaims)
		//lvs := []string{"method", "export", "client", custCl.ClientID, "error", "false"}
		lvs := []string{"method", "export", "error", "false"}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds()) 	//mw.CountResult.Observe(float64(n))
	}(time.Now())

	output, err = mw.Next.Export(ctx, s)
	return
}

func (mw InstrumentingMiddleware) Upload(ctx context.Context, r *http.Request) (output string, err error) {
	defer func(begin time.Time) {
		//custCl, _ := ctx.Value(jwt.JWTClaimsContextKey).(*service.CustomClaims)
		//lvs := []string{"method", "export", "client", custCl.ClientID, "error", "false"}
		lvs := []string{"method", "upload", "error", "false"}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds()) 	//mw.CountResult.Observe(float64(n))
	}(time.Now())

	output, err = mw.Next.Upload(ctx, r)
	return
}