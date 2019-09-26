package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/subtle"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	kitlog "github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"yokitalk.com/docservice/server/dbinstance"
	"yokitalk.com/docservice/server/middlewares"
	"yokitalk.com/docservice/server/service"
)

const maxUploadSize = 100 * 1024 * 1024
const uploadPath  = "./cache/upload"
const downloadPath    = "./cache/download"

const (
	basicAuthUser = "prometheus"
	basicAuthPass = "password"
)

func methodControl(method string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			h.ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func basicAuth(username string, password string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="metrics"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorised\n"))
			return
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	logger := kitlog.NewLogfmtLogger(os.Stderr)
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller)

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

	mysqlManager := dbinstance.GetMysqlInstance()

	defer mysqlManager.Destroy()

	db := mysqlManager.DB

	var ds service.DocService
	ds = service.NewDocService(db)

	ds = middlewares.LoggingMiddleware{Logger: logger, Next: ds}
	ds = middlewares.InstrumentingMiddleware{RequestCount: requestCount, RequestLatency: requestLatency, CountResult: countResult, Next: ds}

	//privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	//fmt.Println(privateKey)
	//
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//publicKey := privateKey.PublicKey
    //fmt.Println(publicKey)
	var secretKey = []byte("abcd1234!@#$")

	/*jwtKeyFunc := func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	}*/

	/*jwtOptions := []httptransport.ServerOption {
		httptransport.ServerErrorEncoder(service.AuthErrorEncoder),
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerBefore(kitjwt.ContextToHTTP()),
	}*/

	importHandler := httptransport.NewServer(
		service.MakeImportEndpoint(ds), //kitjwt.NewParser(jwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory)(service.MakeImportEndpoint(ds)),
		service.DecodeImportRequest,
		service.EncodeResponse,
		//jwtOptions...,
	)

	exportHandler := httptransport.NewServer(
		service.MakeExportEndpoint(ds), //kitjwt.NewParser(jwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory)(service.MakeExportEndpoint(ds)),
		service.DecodeExportRequest,
		service.EncodeResponse,
		//jwtOptions...,
	)

	authFieldKeys := []string{"method", "error"}
	requestAuthCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "auth_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, authFieldKeys)
	requestAuthLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "auth_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, authFieldKeys)

	// API clients database
	var clients = map[string]string{
		"mobile": "m_secret",
		"web":    "w_secret",
	}

	var auth service.AuthService
	auth = service.NewAuthService(secretKey, clients)

	auth = middlewares.LoggingAuthMiddleware{Logger: logger, Next: auth}
	auth = middlewares.InstrumentingAuthMiddleware{RequestCount: requestAuthCount, RequestLatency: requestAuthLatency, Next: auth}

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(service.AuthErrorEncoder),
		httptransport.ServerErrorLogger(logger),
	}

	authHandler := httptransport.NewServer(
		service.MakeAuthEndpoint(auth),
		service.DecodeAuthRequest,
		service.EncodeResponse,
		options...,
	)

	http.Handle("/auth", methodControl("POST", authHandler))

	http.Handle("/import", methodControl("POST", importHandler))
	http.Handle("/export", methodControl("POST", exportHandler))
	http.Handle("/metrics", basicAuth(basicAuthUser, basicAuthPass, promhttp.Handler()))
	http.HandleFunc("/upload", uploadFileHandler())

	fs := http.FileServer(http.Dir(downloadPath))
	http.Handle("/files/", http.StripPrefix("/files", fs))

	logger.Log("msg", "HTTP", "addr", "8080")
	logger.Log("err", http.ListenAndServe(":8080", nil))
}

func getKey() (*ecdsa.PrivateKey, error) {
	prk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		return prk, err
	}

	return prk, nil
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
		w.Write([]byte(fileName+fileExt))
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