package main

import (
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"yokitalk.com/docservice/server/middlewares"
	"yokitalk.com/docservice/server/service"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "yokitalk",
		Subsystem: "doc_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "yokitalk",
		Subsystem: "doc_service",
		Name:      "request_latency_microsenconds",
		Help:      "Total duration of requests in microsenconds",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "yokitalk",
		Subsystem: "doc_service",
		Name:       "count_result",
		Help:        "The reuslt of each count method",
	}, []string{})

	var ds service.Service
	ds = service.DocService{}
	ds = middlewares.LoggingMiddleware{logger, ds}
	ds = middlewares.InstrumentingMiddleware{requestCount, requestLatency, countResult, ds}

	importHandler := httptransport.NewServer(
		service.MakeImportEndpoint(ds),
		service.DecodeImportRequest,
		service.EncodeResponse,
	)

	exportHandler := httptransport.NewServer(
		service.MakeExportEndpoint(ds),
		service.DecodeExportRequest,
		service.EncodeResponse,
	)

	http.Handle("/import", importHandler)
	http.Handle("/export", exportHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", "8080")
	logger.Log("err", http.ListenAndServe(":8080", nil))
}