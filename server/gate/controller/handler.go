package controller

import (
	"fmt"
	"log"
	"ms_sg_back/config"
	"ms_sg_back/constant"
	"ms_sg_back/net"
	"strings"
	"sync"
)

var GateHandler = &Handler{
	proxyMap: make(map[string]map[int64]*net.ProxyClient),
}

// 这里需要设置转发相关的属性
type Handler struct {
	// 代理锁
	proxyMutex sync.RWMutex
	// key是一个string, 表示代理地址, value还是一个map, key是游戏客户端id, value是当前客户端id对应的连接
	proxyMap map[string]map[int64]*net.ProxyClient
	// 账号服务代理
	loginProxy string
	// 游戏服务代理
	gameProxy string
}

func (h *Handler) Router(r *net.Router) {
	h.loginProxy = config.File.MustValue("gate_server", "login_proxy", "ws://127.0.0.1:8003")
	h.gameProxy = config.File.MustValue("gate_server", "game_proxy", "ws://127.0.0.1:8001")
	// 这里处理所有请求
	group := r.Group("*")
	group.AddRouter("*", h.all)
}

func (h *Handler) all(req *net.WsMsgReq, res *net.WsMsgRes) {
	fmt.Println("网关处理器执行中...")
	// 转发请求, 比如account等
	// 1. 创建并初始化代理客户端
	name := req.Body.Name
	proxySrc := ""
	if isAccount(name) {
		// 满足登录服务
		proxySrc = h.loginProxy
	}
	if proxySrc == "" {
		res.Body.Code = constant.ProxyNotInConnect
		return
	}
	// 正确连接
	// 上锁读取代理服务器
	h.proxyMutex.RLock()
	_, ok := h.proxyMap[proxySrc]
	if !ok {
		h.proxyMap[proxySrc] = make(map[int64]*net.ProxyClient)
	}
	h.proxyMutex.RUnlock()
	// 获取客户端id
	originCid, err := req.Conn.GetProperty("cid")
	if err != nil {
		log.Println("cid未找到")
		res.Body.Code = constant.InvalidParam
		return
	}
	cid := originCid.(int64)
	// 从proxyClientMap中获取对应的代理客户端实例
	proxyClient := h.proxyMap[proxySrc][cid]
	if proxyClient == nil {
		// 代理为空说明是首次连接代理
		proxyClient = net.NewProxyClient(proxySrc)
		err := proxyClient.Connect()
		if err != nil {
			h.proxyMutex.Lock()
			delete(h.proxyMap[proxySrc], cid)
			h.proxyMutex.Unlock()
			res.Body.Code = constant.ProxyConnectError
			return
		}
		h.proxyMutex.Lock()
		h.proxyMap[proxySrc][cid] = proxyClient
		h.proxyMutex.Unlock()
		// 存储连接
		proxyClient.SetProperty("cid", cid)
		proxyClient.SetProperty("proxy", proxySrc)
		proxyClient.SetProperty("gateConn", req.Conn)
		proxyClient.SetOnPush(h.onPush)
	}

	// 存在proxyClient则向ws服务器推送消息, 实际上就是透传游戏客户端传过来的消息
	res.Body.Seq = req.Body.Seq
	res.Body.Name = req.Body.Name
	r, _ := proxyClient.Send(req.Body.Name, req.Body.Msg)
	if r != nil {
		res.Body.Code = r.Code
		res.Body.Msg = r.Msg
	} else {
		res.Body.Code = constant.ProxyConnectError
		return
	}
}

func isAccount(name string) bool {
	// 字符串是否为account.开头
	return strings.HasPrefix(name, "account.")
}

func (h *Handler) onPush(connect *net.ClientConnect, body *net.ResBody) {
	gc, err := connect.GetProperty("gateConn")
	if err != nil {
		log.Panicln("onPush gateConn 报错: ", err)
	}
	// 转换为gateConnect, 然后推送给ws服务器
	gateConnect := gc.(net.WsConnect)
	gateConnect.Push(body.Name, body.Msg)
}
