// 网关服务入口
package main

import (
	"ms_sg_back/config"
	"ms_sg_back/net"
	"ms_sg_back/server/gate"
)

/**
1. 登录功能 account.login 需要通过网关 转发到 登录服务
2. 网关(ws对应的客户端) 需要和登录服务器(ws服务端)交互
3. 网关需要和游戏客户端交互, 因此网关也是一个 ws的服务端, 接收游戏客户端推送的数据
4. ws的服务端(登录服务)实际上已经有了初步实现
5. ws的客户端(网关)需要实现
6. 网关: 实际上是一个代理服务器,
	a. 需要代理地址(可能有多个), 主要实现代理连接通道, 用于转发游戏客户端到登录服务的连接请求
	b. 同时作为ws登录服务的客户端, 需要实现面向登录服务的推送请求以及ws连接
7. 网关路由: 接收所有的请求(*通配符), 作为网关的ws服务端功能
8. 握手协议, 检测第一次建立连接时收发是否通畅
*/

func main() {
	host := config.File.MustValue("gate_server", "host", "127.0.0.1")
	port := config.File.MustValue("gate_server", "port", "8004")

	s := net.NewServer(host + ":" + port)
	s.NeedSecret(true)
	gate.Init()
	s.Router(gate.Router)
	// 网关需要加密, 再传递给客户端
	// 路由指令
	// 启动服务
	s.Start()
}
