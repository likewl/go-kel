/**
 * @Description TODO
 * @Author LiKe
 * @Date 13:38 2021/3/19
 * @Version  1.0
 */
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/likewl/go-kel/register/consul"
	"github.com/likewl/go-kel/service"
	"github.com/likewl/go-kel/service/web"
)

func main() {
	//consul配置文件
	Config := api.Config{
		Address: "127.0.0.1:8500",
	}
	//consul注册信息配置文件
	Reg := api.AgentServiceRegistration{
		Name: "test",
	}
	//本工具的consul配置文件
	conf := consul.Config{
		Reg:  &Reg,
		Conf: &Config,
	}
	//初始化路由
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Next()
	})
	r.GET("a", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"a": "a",
		})
	})
	//注册并运行服务
	err := service.Run(web.WebService{
		Handler: r,
		Port:    8080,
	}, &conf)
	if err != nil {
		fmt.Println(err)
	}
}
