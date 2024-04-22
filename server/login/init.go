package login

import (
	"ms_sg_back/net"
	"ms_sg_back/server/login/controller"
)

var Router = net.NewRouter()

// 模块初始化
func Init() {
	initRouter()
}

func initRouter() {
	controller.DefaultAccount.Router(Router)
}
