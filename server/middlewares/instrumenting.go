package middlewares

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"

	"yokitalk.com/docservice/server/service"
)

type InstrumentingMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	CountResult    metrics.Histogram
	Next           service.Service
}

func (mw InstrumentingMiddleware) Login(name, pwd string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "login", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.Next.Login(name, pwd)
	return
}

func (mw InstrumentingMiddleware) Import(s string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "import", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.Next.Import(s)
	return
}

func (mw InstrumentingMiddleware) Export(s string) (n int) {
	defer func(begin time.Time) {
		lvs := []string{"method", "export", "error", "false"}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		mw.CountResult.Observe(float64(n))
	}(time.Now())

	n = mw.Next.Export(s)
	return
}