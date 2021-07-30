// Description TODO
// Author LiKe
// Date 13:05 2021/7/30

package register

import "github.com/hashicorp/consul/api"

type Config interface {
	//Default 初始化配置
	Default(string)
	GetConf() *api.Config
	SetConf(*api.Config)
	GetReg() *api.AgentServiceRegistration
	SetReg(*api.AgentServiceRegistration)
	RegisterService() error
	DeregisterService()
}
