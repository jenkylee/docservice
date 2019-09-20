package middlewares

import (
	"time"

	"github.com/go-kit/kit/log"
	"yokitalk.com/docservice/server/service"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next service.Service
}

func (mw LoggingMiddleware) Import(s string) (output string, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
				"method", "import",
				"input", s,
				"output", output,
				"err", err,
				"took", time.Since(begin),
			)
	}(time.Now())

	output, err = mw.Next.Import(s)
	return
}

func (mw LoggingMiddleware) Export(s string) (n int) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "export",
			"input", s,
			"n", n,
			"took", time.Since(begin),
		)
	}(time.Now())

	n = mw.Next.Export(s)
	return
}
