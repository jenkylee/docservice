package main

import (
	"crypto/subtle"
	"github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	kitlog "github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	kithttp "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"net/http"
	"os"
	"yokitalk.com/docservice/server/dbinstance"
	"yokitalk.com/docservice/server/middlewares"
	"yokitalk.com/docservice/server/service"
)

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

	fieldKeys := []string{"method", "client","error"}
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

	jwtKeyFunc := func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	}

	jwtOptions := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(service.ErrorEncoder),
		kithttp.ServerErrorLogger(logger),
	}

	importHandler := kithttp.NewServer(
		kitjwt.NewParser(jwtKeyFunc, jwt.SigningMethodHS256, service.CustomClaimsFactory)(service.MakeImportEndpoint(ds)),
		//service.MakeImportEndpoint(ds),
		service.DecodeImportRequest,
		service.EncodeResponse,
		append(jwtOptions, kithttp.ServerBefore(kitjwt.HTTPToContext()))...,
	)

	exportHandler := kithttp.NewServer(
		kitjwt.NewParser(jwtKeyFunc, jwt.SigningMethodHS256, service.CustomClaimsFactory)(service.MakeExportEndpoint(ds)),
		//service.MakeExportEndpoint(ds),
		service.DecodeExportRequest,
		service.EncodeResponse,
		append(jwtOptions, kithttp.ServerBefore(kitjwt.HTTPToContext()))...,
	)

	uploadHandler := kithttp.NewServer(
		kitjwt.NewParser(jwtKeyFunc, jwt.SigningMethodHS256, service.CustomClaimsFactory)(service.MakeUploadEndpint(ds)),
		service.DecodeUploadRequest,
		service.EncodeResponse,
		append(jwtOptions, kithttp.ServerBefore(kitjwt.HTTPToContext()))...,
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

	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(service.AuthErrorEncoder),
		kithttp.ServerErrorLogger(logger),
	}

	authHandler := kithttp.NewServer(
		service.MakeAuthEndpoint(auth),
		service.DecodeAuthRequest,
		service.EncodeResponse,
		options...,
	)

	http.Handle("/auth", methodControl("POST", authHandler))

	http.Handle("/import", methodControl("POST", importHandler))
	http.Handle("/export", methodControl("POST", exportHandler))
	http.Handle("/upload", methodControl("POST", uploadHandler))
	http.Handle("/metrics", basicAuth(basicAuthUser, basicAuthPass, promhttp.Handler()))

	http.HandleFunc("/download", downloadFileHandler())
	//fs := http.FileServer(http.Dir(downloadPath))
	//http.Handle("/files/", http.StripPrefix("/files", fs))

	logger.Log("msg", "HTTP", "addr", "8080")
	logger.Log("err", http.ListenAndServe(":8080", nil))
}

func downloadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//文件下载只允许GET方法
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 文件名
		fileName := r.FormValue("filename")

		if fileName == "" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		log.Println("fileName:" + fileName)

		// 打开文件
		file, err := os.Open(downloadPath + "/" +fileName)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		// 结束后关闭文件
		defer file.Close()

		// 设置响应的header头
		w.Header().Add("Content-type", "application/octet-stream")
		w.Header().Add("content-disposition", "attachment; filename=\""+fileName+"\"")
		// 将从文件写入resposeBody
		_, err = io.Copy(w, file)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
	})
}
