package game

import (
	"ms_sg_back/db"
	"ms_sg_back/net"
	"ms_sg_back/server/game/controller"
)

var Router = net.NewRouter()

// 模块初始化
func Init() {
	// 初始化数据库
	db.TestAndInitDB()
	initRouter()
}

func initRouter() {
	controller.DefaultRoleController.Router(Router)
}
