package net

import (
	"log"
	"strings"
)

// 处理器(具体业务逻辑)
type HandlerFunc func(req *WsMsgReq, res *WsMsgRes)

// 路由分组
type group struct {
	prefix     string
	handlerMap map[string]HandlerFunc
}

type Router struct {
	group []*group
}

func NewRouter() *Router {
	return &Router{}

}

// 为group实现exec
func (group *group) exec(name string, req *WsMsgReq, res *WsMsgRes) {
	handler := group.handlerMap[name]
	if handler != nil {
		handler(req, res)
	} else {
		// 降级处理网关路由
		group.execGateWay(req, res)
	}
}

func (group *group) execGateWay(req *WsMsgReq, res *WsMsgRes) {
	handler := group.handlerMap["*"]
	if handler != nil {
		handler(req, res)
	} else {
		// 未找到路由
		log.Panicln("路由未定义")
	}
}

// 添加路由
func (group *group) AddRouter(name string, handler HandlerFunc) {
	group.handlerMap[name] = handler
}

// 构造Group
func (router *Router) Group(prefix string) *group {
	group := &group{
		prefix:     prefix,
		handlerMap: make(map[string]HandlerFunc),
	}

	router.group = append(router.group, group)
	return group
}

// 执行入口
func (router *Router) Run(req *WsMsgReq, res *WsMsgRes) {
	// 处理前端传递的数据
	// req.Body.Name就是路径, 以登录为例, name就是 account.login, account为group标识, login是具体的路由标识

	// 1. 先将路径上的group和具体的路由区分开, 利用split方法转换为数组
	strs := strings.Split(req.Body.Name, ".")
	prefix := ""
	name := ""
	if len(strs) == 2 {
		// 是否符合标准
		prefix = strs[0]
		name = strs[1]
	}

	for _, group := range router.group {
		if group.prefix == prefix {
			// 执行对应的逻辑
			group.exec(name, req, res)
		} else if group.prefix == "*" {
			// 网关服务匹配
			group.execGateWay(req, res)
		}
	}
}
