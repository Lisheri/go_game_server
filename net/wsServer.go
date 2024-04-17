package net

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// ws服务
type wsServer struct {
	wsConnect *websocket.Conn
	router    *router
	// 回复信息通过通道的形式发送给管道, 然后在管道中进一步处理, 本质上是一个 [队列]
	outChan      chan *WsMsgRes
	Seq          int64
	property     map[string]interface{} // 属性, 是一个kv形式, k是string
	propertyLock sync.RWMutex           // 属性锁, 写入的时候先上锁, 在进行写入操作
}

func NewWsServer(wsConnect *websocket.Conn) *wsServer {
	return &wsServer{
		wsConnect: wsConnect,                    // 连接
		Seq:       0,                            // 序列号
		outChan:   make(chan *WsMsgRes, 1000),   // 管道
		property:  make(map[string]interface{}), // 属性
	}
}

// 设置router
func (w *wsServer) Router(router *router) {
	w.router = router
}

// 需要对wsServer实现 WsConnect 这个接口
func (w *wsServer) SetProperty(key string, value interface{}) {
	w.propertyLock.Lock()         // 上锁
	defer w.propertyLock.Unlock() // 解锁
	// 设置属性
	w.property[key] = value
}

func (w *wsServer) GetProperty(key string) (interface{}, error) {
	w.propertyLock.RLock()         // 上锁, 这里要用读取锁
	defer w.propertyLock.RUnlock() // 解锁
	return w.property[key], nil
}

func (w *wsServer) RemoveProperty(key string) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	// 移除
	delete(w.property, key)
}

func (w *wsServer) Addr() string {
	// 直接获取ws的远程地址
	return w.wsConnect.RemoteAddr().String()
}

func (w *wsServer) Push(name string, data interface{}) {
	// 封装回复的消息
	res := &WsMsgRes{Body: &ResBody{Name: name, Msg: data, Seq: 0}}
	// 将res对象添加到管道中进一步处理
	w.outChan <- res
}

// 启动读写数据的处理逻辑
func (w *wsServer) Start() {
	// 通道一旦建立, 那么收发消息需要一直监听
	// 创建两个协程, 一个读, 一个发
	go w.readMsgLoop()
	go w.writeMsgLoop()
}

func (w *wsServer) readMsgLoop() {
	// 先获取到客户端发送的数据, 然后处理数据, 最后恢复消息
	// 回复消息需要先经过路由, 也就是实际处理程序, 然后才能经过回复
	// res := &WsMsgRes{Body: &ResBody{Name: name, Msg: data, Seq: 0}}
	// 将res对象添加到管道中进一步处理
	// w.outChan <- res
	defer func() {
		// 出现问题也需要关闭
		if err := recover(); err != nil {
			log.Fatal(err)
			w.Close()
		}
	}()
	for {
		_, data, err := w.wsConnect.ReadMessage()
		if err != nil {
			log.Println("收消息出现错误", err)
			break
		}
		fmt.Println("收到了", data)
	}
	// 跳出循环说明中间结束或者出现问题, 同样执行关闭
	w.Close()
}

func (w *wsServer) writeMsgLoop() {
	for {
		select {
		case msg := <-w.outChan:
			// TODO 暂不处理
			fmt.Println(msg)
		}
	}
}

// 关闭ws连接
func (w *wsServer) Close() {
	_ = w.wsConnect.Close()
}
