package net

// 请求体
type ReqBody struct {
	Seq   int64       `json:"seq"`   // 序号标识
	Name  string      `json:"name"`  // 路由
	Msg   interface{} `json:"msg"`   // 消息信息
	Proxy string      `json:"proxy"` // 用于多线程调度, 服务需要调用服务, 因此需要一个Proxy做代理
}

// 响应体
type ResBody struct {
	Seq  int64       `json:"seq"`  // 序号
	Name string      `json:"name"` // 发送消息名称(也是类似路由)
	Code int         `json:"code"` // 状态码
	Msg  interface{} `json:"msg"`  // 回复内容
}

// 请求
type WsMsgReq struct {
	Body *ReqBody
	Conn WsConnect
}

// 响应
type WsMsgRes struct {
	Body *ResBody
}

// 存储数据等
// 可以理解为 request请求, 会有参数, 主要用于请求中存取参数
type WsConnect interface {
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
	Addr() string
	Push(name string, data interface{})
}

type Handshake struct {
	Key string `json:"key"`
}

// 心跳检测相关
type Heartbeat struct {
	CTime int64 `json:"ctime"`
	STime int64 `json:"stime"`
}
