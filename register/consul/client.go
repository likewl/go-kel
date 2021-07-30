package consul

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"math/rand"
	"strconv"
)

//NewClient 从 config 中配置中获取一个注册中心
func (c *Config) NewClient() (err error) {
	c.client, err = api.NewClient(c.Conf)
	return
}

//getService 从consul获取以随机的方式获取服务
func (c *Config) getService(serviceName string, tag string, q *api.QueryOptions) (*api.CatalogService, error) {
	result, _, err := c.client.Catalog().Service(serviceName, tag, q)
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
func GetClientConn(c Config,serviceName string, tag string, q *api.QueryOptions, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	client, err := c.getService(serviceName, tag, q)
	if err != nil {
		return nil, err
	}
	if len(opts) == 0 {
		opts = append(opts, grpc.WithInsecure())
	}
	return grpc.Dial(client.ServiceAddress+":"+strconv.Itoa(client.ServicePort), opts...)
}
