package net

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type server struct {
	addr   string  // 地址
	router *router // 路由
}

// 用于初始化server
func NewServer(addr string) *server {
	return &server{
		addr: addr,
	}
}

// 启动服务
func (s *server) Start() {
	// 用于处理请求
	http.HandleFunc("/", s.wsHandler)
	// 监听服务并且启动
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		panic(err)
	}
}

// 升级http到ws(需要用到 github.com/gorilla/websocket)
var wsUpgrader = websocket.Upgrader{
	// 全部允许跨域跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *server) wsHandler(w http.ResponseWriter, r *http.Request) {
	// 这是一个webSocket请求
	// 1. 基于http初始化webSocket
	// 返回一个websockaet链接
	wsConnect, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		// * 日志模块会导致终止
		log.Fatal("ws服务链接出错")
	}
	// log.Fatal("ws服务连接成功")
	// 消息需要约定格式, 用于处理收发信息
	// err = wsConnect.WriteMessage(websocket.BinaryMessage, []byte("hello"))
	// fmt.Println(err);
	wsServer := NewWsServer(wsConnect)
	wsServer.Router(s.router)
	wsServer.Start()
	// 读取信息
	// 客户端格式约束:  { Name: "account.login" }, 收到之后进行解析

}
