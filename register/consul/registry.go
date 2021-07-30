package consul

import (
	"log"
)



//RegisterService 注册服务到注册中心
func (c *Config) RegisterService() (err error) {
	if c.Conf.Address == "" {
		log.Println("-----WARNING:没有获取到注册中心配置------")
		return nil
	}
	err = c.NewClient()
	if err != nil {
		return err
	}
	err = c.client.Agent().ServiceRegister(c.Reg)
	if err != nil {
		return err
	}
	log.Printf("Register Service success addr: %s:%d\n", c.Reg.Address, c.Reg.Port)
	return nil
}

//DeregisterService 反注册
func (c *Config) DeregisterService() {
	if c.Conf.Address != "" {
		e := c.deregister()
		if e != nil {
			log.Printf("Service:%s Deregister fail msg: %s \n", c.Reg.ID, e)
		} else {
			log.Printf("Service:%s Deregister success \n", c.Reg.ID)
		}
	}
}

//deregister 反注册,从 consul 中删除服务
func (c *Config) deregister() error {
	return c.client.Agent().ServiceDeregister(c.Reg.ID)
}
