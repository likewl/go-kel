# go-kel

一个简易grpc+web服务器微服务搭建脚手架

# 工具特点
1. 支持web和grpc服务，同时兼容两种服务用一个端口运行
    - 目前只要涉及到web只支持http1.0服务
2. 支持像指定注册中心注册和注销服务
    - 目前只支持consul
   
#工具使用方法
```go
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
	//工具配置文件
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

```

功能还在添加中。。。