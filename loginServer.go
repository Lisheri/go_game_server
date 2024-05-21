// 游戏服务入口
package main

import (
	"ms_sg_back/config"
	"ms_sg_back/net"
	"ms_sg_back/server/login"
)

func main() {
	host := config.File.MustValue("login_server", "host", "127.0.0.1")
	port := config.File.MustValue("login_server", "port", "8003")

	s := net.NewServer(host + ":" + port)
	login.Init()
	s.Router(login.Router)
	// 路由指令
	// 启动服务
	s.Start()
}
