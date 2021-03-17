package consul

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/likewl/go-kel/util/ip"
	"github.com/likewl/go-kel/util/rand"
)

var (
	Client *api.Client
	Port   = 8000
)
//newClient 从 config 中配置中获取一个注册中心
func newClient(config *api.Config) (*api.Client, error) {

	//如果没有设置consul地址
	//默认初始化 ::8500
	if config.Address == "" {
		localIP, ok := ip.GetIP()
		if !ok {
			return nil, errors.New("net.Interfaces failed")
		}
		config.Address = localIP + ":8500"
	}
	var err error
	Client, err = api.NewClient(config)
	return Client, err
}
//registerService 注册服务到注册中心
func registerService(config *api.Config, reg *api.AgentServiceRegistration) (err error) {
	Client, err = newClient(config)
	if err != nil {
		return err
	}
	//初始化参数
	//if reg.Address == "" {
	//	localIP, _ := ip.GetIP()
	//	reg.Address = localIP
	//}

	if reg.ID == "" {
		if reg.Name != "" {

			reg.ID = reg.Name + "-"
		}
		reg.ID = reg.ID + rand.SplitStringToName(rand.GenerateString("ID"))
	}
	if len(reg.Tags) == 0 {
		var tag string
		if reg.Name != "" {
			tag = reg.Name + "-"
		}
		tag = tag + rand.GenerateString("Tags")
		reg.Tags = append(reg.Tags, tag)
	}
	if reg.Name == "" {
		reg.Name = rand.GenerateString("Name")
	}
	if reg.Port == 0 {
		reg.Port = Port
	}
	//注册服务
	err = Client.Agent().ServiceRegister(reg)
	if err != nil {
		return err
	}
	fmt.Printf("Register Service success addr: %s:%d\n", reg.Address, reg.Port)
	return nil
}
//DeregisterService 反注册，从 consul 中删除服务
func deregisterService(reg *api.AgentServiceRegistration) {
	Client.Agent().ServiceDeregister(reg.ID)
	fmt.Printf("Service:%s Deregister success",reg.ID)
}
