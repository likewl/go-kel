// Description grpc服务
// Author LiKe
// Date 21:47 2021/7/29

package grpc

import (
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

type GrpcService struct {
	Port int //端口号
	grpc *grpc.Server
}

func (g GrpcService) GetType() string {
	return "grpcService"
}

func (g GrpcService) GetPort() string {
	return ":" + strconv.Itoa(g.Port)
}
func (g GrpcService) Run() error {
	lis, err := net.Listen("tcp", g.GetPort())
	if err != nil {
		return errors.New("net.Listen err: " + err.Error())
	}
	fmt.Println("service run success port: ", g.Port)
	return g.grpc.Serve(lis)
}
