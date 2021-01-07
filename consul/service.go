package consul

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
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

func RunRPCServiceAndRegistry(service *grpc.Server, port string, config *api.Config, reg *api.AgentServiceRegistration) (err error) {
	var p int
	p, err = strconv.Atoi(strings.TrimPrefix(port, ":"))
	if err != nil {
		err = errors.New("port error")
		return
	}
	Port = p
	err = RegisterService(config, reg)
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
			DeregisterService(reg)
			log.Fatalf("net.Listen err: %v", err2)
		}
		fmt.Println("service run success")
		errc <- service.Serve(lis)
	}()
	err = <-errc
	DeregisterService(reg)
	return
}
func RunWebServiceAndRegistry(service http.Handler, port string, config *api.Config, reg *api.AgentServiceRegistration) (err error) {
	var p int
	p, err = strconv.Atoi(strings.TrimPrefix(port, ":"))
	if err != nil {
		err = errors.New("port error")
		return
	}
	Port = p
	if reg.Check == nil {
		var path string
		if reg.Address == "" {
			path = "localhost"
		}
		reg.Check = &api.AgentServiceCheck{
			Interval: "5s",
			HTTP:     "http://" + path + port+"/consul_check/health",
		}
	}
	err = RegisterService(config, reg)
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
		fmt.Println("service run success")
		errc <- http.ListenAndServe(port, service)
	}()
	err = <-errc
	DeregisterService(reg)
	return
}
