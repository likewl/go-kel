// Description 服务器相关的api
// Author LiKe
// Date 22:04 2021/7/29

package service

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/likewl/go-kel/register"
	"os"
	"os/signal"
	"syscall"
)

//Run 完成配置初始化工作,并运行服务
func Run(ser Service, c register.Config) error {
	if ser == nil || ser.GetPort() == ":0" {
		return errors.New("service should not empty or port is null")
	}
	if c.GetReg() == nil {
		c.SetReg(new(api.AgentServiceRegistration))
	}
	if c.GetConf() == nil {
		c.SetConf(new(api.Config))
	}
	c.Default(ser.GetPort())
	//注册服务
	err := c.RegisterService()
	if err != nil {
		return err
	}
	//监听信号
	errc := make(chan error)
	s := make(chan os.Signal)
	go func() {
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-s)
	}()
	//运行服务
	go func() {
		errc <- ser.Run()
	}()
	err = <-errc
	//发生错误 反注册
	c.DeregisterService()
	return errors.New("system has been stopped by signal : " + err.Error())
}
