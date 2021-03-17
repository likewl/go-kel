package consul

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"math/rand"
	"strconv"
)

//getService 从consul获取以随机的方式获取服务
func getService(serviceName string, tag string, q *api.QueryOptions) (*api.CatalogService, error) {
	result, _, err := Client.Catalog().Service(serviceName, tag, q)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("没有查询到服务")
	}
	service := result[rand.Int()%len(result)]
	return service, nil
}

//GetClientConn 从consul获取客户端连接
func GetClientConn(serviceName string, tag string, q *api.QueryOptions, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	client, err := getService(serviceName, tag, q)
	if err != nil {
		return nil, err
	}
	if len(opts) == 0 {
		opts = append(opts, grpc.WithInsecure())
	}
	return grpc.Dial(client.ServiceAddress+":"+strconv.Itoa(client.ServicePort), opts...)
}
