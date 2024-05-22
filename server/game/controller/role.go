package controller

import "ms_sg_back/net"

var DefaultRoleController = &RoleController{}

type RoleController struct {}

func (r *RoleController) Router(router *net.Router) {
	g := router.Group("role")
	g.AddRouter("enterServer", r.enterServer)
}

// 进入游戏的逻辑
func (r *RoleController) enterServer(req *net.WsMsgReq, res *net.WsMsgRes) {
	// 1. 验证session是否合法, 合法时获取出登录的用户id
	// 2. 根据uid查询对应的游戏角色, 若无角色, 则提示无角色, 然后等待接收请求进行创建
	// 3. 根据角色查询资源, 若有资源则返回, 若是新角色则初始化角色对应的资源并入库
}

