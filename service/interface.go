// Description 服务器的通用接口
// Author LiKe
// Date 20:38 2021/7/29

package service

type Service interface {
	GetType() string //获取服务类型
	GetPort() string //获取端口号
	Run() error      //启动服务
}
