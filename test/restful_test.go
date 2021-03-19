/**
 * @Description TODO
 * @Author LiKe
 * @Date 13:38 2021/3/19
 * @Version  1.0
 */
package test

import (
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/likewl/go-kel/consul"
	"testing"
)
var Config api.Config
var Reg api.AgentServiceRegistration
func TestRestful(t *testing.T) {
	//consul配置文件
	Config = api.Config{
		Address: "localhost:8500",
	}
	//服务配置文件
	Reg = api.AgentServiceRegistration{
		Name: "test",
	}
	//初始化路由
	r := gin.New()
	r.GET("/a", func(c *gin.Context) {
		c.JSON(200,gin.H{
			"a":"a",
		})
	})
	//运行服务
	consul.RunWebServiceAndRegistry(r,nil, ":8500", &Config, &Reg)
}
