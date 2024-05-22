package net

import (
	"encoding/json"
	"errors"
	"log"
	"ms_sg_back/utils"
	"sync"
	"time"

	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

// ws服务
type wsServer struct {
	wsConnect *websocket.Conn
	router    *Router
	// 回复信息通过通道的形式发送给管道, 然后在管道中进一步处理, 本质上是一个 [队列]
	outChan      chan *WsMsgRes
	Seq          int64
	property     map[string]interface{} // 属性, 是一个kv形式, k是string
	propertyLock sync.RWMutex           // 属性锁, 写入的时候先上锁, 在进行写入操作
	needSecret   bool                   // 是否需要加密 (仅网关与客户端交流需要加密)
}

var cid int64 = 0 // 暂存客户端id(自增)
func NewWsServer(wsConnect *websocket.Conn, needSecret bool) *wsServer {
	s := &wsServer{
		wsConnect:  wsConnect,                    // 连接
		Seq:        0,                            // 序列号
		outChan:    make(chan *WsMsgRes, 1000),   // 管道
		property:   make(map[string]interface{}), // 属性
		needSecret: needSecret,                   // 是否加密标识
	}
	// 客户端ic自增
	cid++
	// 设置当前服务对应的客户端cid, 防止在网关服务中获取cid取不到
	s.SetProperty("cid", cid)
	return s
}

// 设置router
func (w *wsServer) Router(router *Router) {
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
	// 拿不到则抛错, 通过ok判断属性是否存在
	if value, ok := w.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("属性不存在")
	}
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

func (w *wsServer) Write(msg *WsMsgRes) {
	// msg.Body转换为json
	data, err := json.Marshal(msg.Body)
	if err != nil {
		log.Panicln("ws服务端写入数据格式非法: ", err)
	}
	secretKey, err := w.GetProperty("secretKey")
	if err == nil {
		// 加密
		key := secretKey.(string)
		// 对data加密
		data, _ = utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
	}
	// 压缩
	if data, err := utils.Zip(data); err == nil {
		// 将数据写回去
		w.wsConnect.WriteMessage(websocket.BinaryMessage, data)
	} else {
		log.Println("ws服务端写数据出错", err)
	}
	d, _ := json.Marshal(msg)
	log.Println("ws服务端写数据", string(d))
}

// 启动读写数据的处理逻辑
func (w *wsServer) Start() {
	// 通道一旦建立, 那么收发消息需要一直监听
	// 创建两个协程, 一个读, 一个发, 读完了就写
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
			log.Println("ws服务端捕捉异常: ", err)
			w.Close()
		}
	}()
	for {
		_, data, err := w.wsConnect.ReadMessage()
		if err != nil {
			log.Println("ws服务收消息出现错误", err)
			break
		}
		// fmt.Println("收到了", data)
		// 收到消息后, 需要对消息进行解析, 前端发过来的格式是json
		// 1. 解压data, unzip
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("解压数据出错, 非法格式: ", err)
		}
		// 2. 如果消息是需要加密的, 则需要进行解密(仅游戏客户端与网关交流需要如下操作, 网关与游戏服务端无需加解密, 全程服务器完成)
		if w.needSecret {
			secretKey, err := w.GetProperty("secretKey")
			if err == nil {
				// 有加密
				key := secretKey.(string)
				// 利用秘钥进行解密
				d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
				if err != nil {
					log.Println("数据格式错误, 解密失败: ", err)
					// TODO 出错, 发握手(正常是游戏客户端在初始化连接时发起握手)
					// w.Handshake()
					// 出错之后, 同样需要重新握手, 更换秘钥
					w.Handshake()
				} else {
					data = d
				}
			} else {
				log.Println("获取密钥出错: ", err)
			}
		}
		// 3. data -> body
		body := &ReqBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			log.Println("ws服务器json解析错误, 格式非法: ", err)
		} else {
			// 获取到前端传递的数据了, 需要利用这些数据处理具体的业务
			req := &WsMsgReq{Conn: w, Body: body}
			res := &WsMsgRes{Body: &ResBody{Name: body.Name, Seq: body.Seq}}
			if req.Body.Name == "heartbeat" {
				// 心跳回复
				h := &Heartbeat{}
				// 将msg解析为心跳
				mapstructure.Decode(req.Body.Msg, h)
				// ns -> ms
				h.STime = time.Now().UnixNano() / 1e6
				// 设回去
				res.Body.Msg = h
			} else {
				if w.router != nil {
					w.router.Run(req, res)
				}
			}
			// res写入管道outChan
			w.outChan <- res
		}
	}
	// 跳出循环说明中间结束或者出现问题, 同样执行关闭
	w.Close()
}

func (w *wsServer) writeMsgLoop() {
	for {
		select {
		// 从管道中获取数据, 先进先出, 从队尾获取
		case msg := <-w.outChan:
			w.Write(msg)
		}
	}
}

// 关闭ws连接
func (server *wsServer) Close() {
	_ = server.wsConnect.Close()
}

// 握手相关常量
const HandshakeMsg = "handshake"

// 当游戏客户端发送请求的时候, 会先进行握手协议
// 后端会发送对应的加密key给客户端
// 客户端在发送数据的时候, 会使用此key进行加密处理
// ? 一旦断开链接, 重新生成链接时, 需要重新进行握手, 生成加密解密数据的秘钥 key
func (server *wsServer) Handshake() {
	// 首先获取secretKey
	key := ""
	secretKey, err := server.GetProperty("secretKey")
	if err == nil {
		// 转string
		key = secretKey.(string)
	} else {
		// 报错说明没有值, 随机生成一个key
		key = utils.RandSeq(16)
	}
	// 封装一个message方法, 用于握手
	handshake := &Handshake{Key: key}

	body := &ResBody{
		Name: HandshakeMsg,
		Msg:  handshake,
	}
	// 先处理为json, 在进行压缩传递
	if data, err := json.Marshal(body); err == nil {
		// 防止key发生变化, 需要重新设置
		if key != "" {
			server.SetProperty("secretKey", key)
		} else {
			server.RemoveProperty("secretKey")
		}
		// data, _ = utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
		// 压缩
		if data, err = utils.Zip([]byte(data)); err == nil {
			// 发送握手消息(实际上就是写入数据)
			server.wsConnect.WriteMessage(websocket.BinaryMessage, data)
		}
	}
}
