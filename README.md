# go-kel
一个基于consul的简易grpc+web服务器搭建脚手架

# 本工具是一个快速开发web和rpc服务的脚手架
## web服务 是基于gin

并默认开始随机给服务命名和check功能
### 启动服务
consul.RunWebServiceAndRegistry(service *gin.Engine, port string, config *api.Config, reg *api.AgentServiceRegistration) (err error)

## rpc服务 是基于grpc

并默认开始随机给服务命名
### 启动服务
consul.RunRPCServiceAndRegistry(service *grpc.Server, port string, config *api.Config, reg *api.AgentServiceRegistration) (err error)

功能还在添加中。。。
