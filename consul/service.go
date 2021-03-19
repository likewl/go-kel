package consul

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

//RunRPCServiceAndRegistry 运行并注册服务到注册中心
//关闭服务时，运行反注册方法，从注册中心删除
func RunRPCServiceAndRegistry(service *grpc.Server, port string, config *api.Config, reg *api.AgentServiceRegistration) (err error) {
	var p int
	p, err = strconv.Atoi(strings.TrimPrefix(port, ":"))
	if err != nil {
		err = errors.New("port error")
		return
	}
	Port = p
	err = registerService(config, reg)
	if err != nil {
		return
	}

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()
	go func() {
		lis, err2 := net.Listen("tcp", port)
		if err2 != nil {
			deregisterService(reg)
			log.Fatalf("net.Listen err: %v", err2)
		}
		fmt.Println("service run success")
		errc <- service.Serve(lis)
	}()
	err = <-errc
	deregisterService(reg)
	return
}
func RunWebServiceAndRegistry(r *gin.Engine,grpcSvc *grpc.Server, port string, config *api.Config, reg *api.AgentServiceRegistration) (err error) {
	var p int
	p, err = strconv.Atoi(strings.TrimPrefix(port, ":"))
	if err != nil {
		err = errors.New("port error")
		return
	}
	Port = p
	if reg.Check == nil {
		var ip string
		if reg.Address == "" {
			ip = "localhost"
		}else{
			ip = reg.Address
		}
		router := "http://" + ip + port+"/consul_check/health"
		reg.Check = &api.AgentServiceCheck{
			Interval: "5s",
			HTTP:     router,
		}
	}
	err = registerService(config, reg)
	if err != nil {
		return
	}

	//初始化健康检测路由
	r.GET("/consul_check/health", func(c *gin.Context) {
		c.JSON(200,gin.H{
			"statue":"ok",
		})
	})

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()
	go func() {
		
		fmt.Println("service run success")
		errc <- http.ListenAndServe(port, grpcHandlerFunc(grpcSvc, r))
		
	}()
	err = <-errc
	deregisterService(reg)
	return
}
//grpcHandlerFunc 方法是对restful和grpc请求的一个分流判断，
//使用 h2c.NewHandler 方法进行了特殊处理，并返回一个 http.Handler ，
//主要的内部逻辑是拦截了所有 h2c 流量，然后根据不同的请求流量类型将其劫持并重定向到相应的 Hander 中去处理。
func grpcHandlerFunc(grpcServer *grpc.Server, httpHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			httpHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}