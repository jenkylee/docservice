package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"yokitalk.com/docservice/server/middlewares"
	"yokitalk.com/docservice/server/service"
)

const maxUploadSize = 100 * 1024 * 1024
const uploadPath  = "./cache/upload"

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
	http.HandleFunc("/upload", uploadFileHandler())

	fs := http.FileServer(http.Dir(uploadPath))
	http.Handle("/files/", http.StripPrefix("/files", fs))

	logger.Log("msg", "HTTP", "addr", "8080")
	logger.Log("err", http.ListenAndServe(":8080", nil))
}

func uploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// validate file size
		/*r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}*/

		// parse and validate file and post parameters
		file, handler, err := r.FormFile("file")
		defer file.Close()

		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		if handler.Size > maxUploadSize {
			renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}

		fileExt := strings.ToLower(path.Ext(handler.Filename))

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		// check file type, detectcontenttype only needs the first 512 bytes
		filetype := http.DetectContentType(fileBytes)
		fmt.Println("文件类型", filetype)
		switch filetype {
		case "image/jpeg", "image/jpg":
		case "image/gif", "image/png":
		case "application/pdf":
		case "application/zip":
			break
		default:
			renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
			return
		}
		fileName := randToken(12)
		// fileEndings, err := mime.ExtensionsByType(fileExt)
		/*if err != nil {
			renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
			return
		}*/
		newPath := filepath.Join(uploadPath, fileName+fileExt)
		fmt.Printf("FileType: %s, File: %s\n", fileExt, newPath)

		// write file
		newFile, err := os.Create(newPath)
		if err != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		defer newFile.Close() // idempotent, okay to call twice
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("SUCCESS"))
	})
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}