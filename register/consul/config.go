// Description TODO
// Author LiKe
// Date 20:52 2021/7/29

package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/likewl/go-kel/util/rand"
	"strconv"
	"strings"
)

type Config struct {
	client *api.Client
	port   string
	Conf   *api.Config
	Reg    *api.AgentServiceRegistration
}

func (c *Config) SetConf(config *api.Config) {
	c.Conf = config
}

func (c *Config) SetReg(registration *api.AgentServiceRegistration) {
	c.Reg = registration
}

func (c *Config) GetConf() *api.Config {
	return c.Conf
}
func (c *Config) GetReg() *api.AgentServiceRegistration {
	return c.Reg
}

//defaultCheck 默认检查地址
func (c *Config) defaultCheck() {
	if c.Reg.Check == nil {
		var ip string
		if c.Reg.Address == "" {
			ip = "localhost"
		} else {
			ip = c.Reg.Address
		}
		router := "http://" + ip + c.port + "/check/health"
		c.Reg.Check = &api.AgentServiceCheck{
			Interval: "5s",
			HTTP:     router,
		}
	}
}

//defaultID 默认ID
func (c *Config) defaultID() {
	if c.Reg.ID == "" {
		c.defaultName()
		c.Reg.ID = c.Reg.Name + "-" + rand.SplitStringToName(rand.GenerateString("ID"))
	}
}

//defaultTags 默认Tags
func (c *Config) defaultTags() {
	if len(c.Reg.Tags) == 0 {
		var tag string
		if c.Reg.Name != "" {
			tag = c.Reg.Name + "-"
		}
		tag = tag + rand.GenerateString("Tags")
		c.Reg.Tags = append(c.Reg.Tags, tag)
	}
}

//defaultName 默认名字
func (c *Config) defaultName() {
	if c.Reg.Name == "" {
		c.Reg.Name = rand.GenerateString("Name")
	}
}

//defaultPort 默认端口
func (c *Config) defaultPort(port string) {
	if c.port == "" {
		c.port = port
	}
	if c.Reg.Port == 0 {
		p, _ := strconv.Atoi(strings.TrimPrefix(c.port, ":"))
		c.Reg.Port = p
	}
}

//Default 初始化配置
func (c *Config) Default(port string) {
	c.defaultName()
	c.defaultID()
	c.defaultTags()
	c.defaultPort(port)
	c.defaultCheck()
}
