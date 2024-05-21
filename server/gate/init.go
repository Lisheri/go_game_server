package gate

import (
	"ms_sg_back/net"
	"ms_sg_back/server/gate/controller"
)

var Router = &net.Router{}

// 模块初始化
func Init() {
	// 初始化数据库
	// 网关同样是初始化路由
	initRouter()
}

func initRouter() {
	controller.GateHandler.Router(Router)
}
