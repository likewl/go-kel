// Description web服务
// Author LiKe
// Date 20:41 2021/7/29

package web

import (
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
	"strings"
)

type WebService struct {
	Port    int
	Handler http.Handler
	Grpc    *grpc.Server
}

func (w WebService) GetType() string {
	return "webService"
}

func (w WebService) GetPort() string {
	return ":" + strconv.Itoa(w.Port)
}
func (w WebService) Run() error {
	fmt.Println("service run success port:", w.Port)
	return http.ListenAndServe(w.GetPort(), w.grpcHandlerFunc())
}

//grpcHandlerFunc 方法是对restful和grpc请求的一个分流判断,
//使用 h2c.NewHandler 方法进行了特殊处理，并返回一个 http.Handler,
//主要的内部逻辑是拦截了所有 h2c 流量，然后根据不同的请求流量类型将其劫持并重定向到相应的 Handler 中去处理。
func (w *WebService) grpcHandlerFunc() http.Handler {
	mux := http.NewServeMux()
	handlerFunc := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("content-type", "application/json; charset=utf-8")
		resp.Write([]byte("{\"statue\":\"ok\"}"))
	})
	mux.HandleFunc("/check/health", handlerFunc)
	return h2c.NewHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			w.Grpc.ServeHTTP(rw, r)
		} else if r.RequestURI == "/check/health" {
			mux.ServeHTTP(rw, r)
		} else {
			w.Handler.ServeHTTP(rw, r)
		}
	}), &http2.Server{})
}
