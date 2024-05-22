// 游戏服务本体
package main

import (
	"ms_sg_back/config"
	"ms_sg_back/net"
	"ms_sg_back/server/game"
)

/**
1. 进入游戏服务, 登录已经完成
2. 根据用户id查询角色, 若没有角色, 则需要创建角色
3. 设置资源(木材, 铁, 令牌, 金钱, 主城, 武将等)数据初始化
4. 地图初始化, 标识玩家自身的地图资源等
5. 资源, 军队, 城池, 武将等
*/

func main() {
	host := config.File.MustValue("game_server", "host", "127.0.0.1")
	port := config.File.MustValue("game_server", "port", "8001")

	s := net.NewServer(host + ":" + port)
	// 游戏服务为内部服务, 无需进行加密, 设置标识为false
	s.NeedSecret(false)
	game.Init()
	s.Router(game.Router)
	// 路由指令
	// 启动服务
	s.Start()
}
