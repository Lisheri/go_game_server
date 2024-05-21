package net

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"ms_sg_back/utils"
	"sync"
	"time"

	"ms_sg_back/constant"

	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

type syncCtx struct {
	// Goroutine(协程)的上下文, 包含Goroutine运行状态, 环境, 现场等信息
	ctx    context.Context
	cancel context.CancelFunc
	// 一旦从代理服务器接收到数据, 就会丢到outChan中
	outChan chan *ResBody
}

func NewSyncCtx() *syncCtx {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	return &syncCtx{
		ctx:    ctx,
		cancel: cancel,
		// 无需缓冲, 一直等待接收即可, 有一个就处理一个
		outChan: make(chan *ResBody),
	}
}

// 实现接收方法
func (s *syncCtx) waitMsg() *ResBody {
	select {
	case msg := <-s.outChan:
		return msg
	case <-s.ctx.Done():
		// 超时
		log.Println("代理服务器响应超时")
		return nil
	}
}

// 客户端连接通道
type ClientConnect struct {
	wsConnect *websocket.Conn
	// 监听当前客户端是否为关闭状态
	isClosed bool
	// 暂存属性
	property map[string]interface{}
	// 存取属性时专用锁
	propertyLock sync.RWMutex
	// 序列号
	Seq              int64
	handshake        bool
	handshakeChannel chan bool
	// 向代理服务器发送消息
	onPush func(connect *ClientConnect, body *ResBody)
	// 连接关闭钩子
	onClose func(connect *ClientConnect)
	// 客户端id -> 读取协程, 需要协程上下文, 还是为了处理代理服务器发送到通道时的读写超时设置
	syncCtxMap map[int64]*syncCtx
	// 发送消息时专用锁
	syncCtxLock sync.RWMutex
}

func (c *ClientConnect) Start() bool {
	// 开启消息接收
	// 同时等待握手消息的返回
	c.handshake = false
	// 开启读消息协程
	go c.wsReadLoop()
	return c.waitHandShake()
}

func (c *ClientConnect) waitHandShake() bool {
	// 等待握手的成功, 也就是等待握手的消息(通过channel通知)
	// TODO 需要超时处理
	// 这里设置10s超时, context.WithTimeout专用于协程内部处理超时连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// 调用完成后, 直接调用cancel取消超时监听
	defer cancel()
	select {
	case _ = <-c.handshakeChannel:
		// 收到了
		log.Println("握手成功")
		return true
	case <-ctx.Done():
		// 一旦超过超时设置, 就会向这个Done() channel 中写入对应值
		log.Println("握手超时了")
		return false
	}
}

// 需要对wsServer实现 WsConnect 这个接口
func (c *ClientConnect) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()         // 上锁
	defer c.propertyLock.Unlock() // 解锁
	// 设置属性
	c.property[key] = value
}

func (c *ClientConnect) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()         // 上锁, 这里要用读取锁
	defer c.propertyLock.RUnlock() // 解锁
	// 拿不到则抛错, 通过ok判断属性是否存在
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("属性不存在")
	}
}

func (c *ClientConnect) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	// 移除
	delete(c.property, key)
}

func (c *ClientConnect) Addr() string {
	// 直接获取ws的远程地址
	return c.wsConnect.RemoteAddr().String()
}

func (c *ClientConnect) Push(name string, data interface{}) {
	// 封装回复的消息
	res := &WsMsgRes{Body: &ResBody{Name: name, Msg: data, Seq: 0}}
	// 将res对象添加到管道中进一步处理
	// 需要加密写入
	c.write(res.Body)
	// c.outChan <- res
}

func (c *ClientConnect) write(body interface{}) error {
	// msg.Body转换为json
	data, err := json.Marshal(body)
	if err != nil {
		log.Panicln("数据格式非法: ", err)
		return err
	}
	secretKey, err := c.GetProperty("secretKey")
	if err == nil {
		// 加密
		key := secretKey.(string)
		// 对data加密
		data, err = utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
		if err != nil {
			log.Println("加密失败")
			return err
		}
	}
	// 压缩
	if data, err := utils.Zip(data); err == nil {
		// 将数据写回去
		err := c.wsConnect.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			log.Println("写入消息失败")
			return err
		}
	} else {
		log.Println("压缩数据失败")
		return err
	}
	return nil
}

func (c *ClientConnect) wsReadLoop() {
	// for {
	// 	_, data, err := c.wsConnect.ReadMessage()
	// 	fmt.Println(data)
	// 	fmt.Println(err)
	// // 收到握手消息, 发送一个channel
	// c.handshake = true
	// // 发送消息
	// c.handshakeChannel <- true
	// 	// 读取消息时, 可能会收到很多消息, 比如 握手, 心跳, 请求信息(account.login)
	// 	// 服务端 写消息
	// }
	defer func() {
		// 出现问题也需要关闭
		if err := recover(); err != nil {
			log.Println("捕捉异常: ", err)
			c.Close()
		}
	}()
	for {
		_, data, err := c.wsConnect.ReadMessage()
		if err != nil {
			log.Println("收消息出现错误", err)
			break
		}
		fmt.Println("收到了", data)
		// 收到消息后, 需要对消息进行解析, 前端发过来的格式是json
		// 1. 解压data, unzip
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("解压数据出错, 非法格式: ", err)
		}
		// 2. 前端消息是加密的, 需要进行解密
		secretKey, err := c.GetProperty("secretKey")
		if err == nil {
			// 有加密
			key := secretKey.(string)
			d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
			if err != nil {
				log.Println("数据格式错误, 解密失败: ", err)
			} else {
				data = d
			}
		} else {
			log.Println("获取密钥出错: ", err)
		}
		// 3. data -> body
		body := &ResBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			log.Println("数据格式有误, 格式非法: ", err)
		} else {
			// 获取到前端传递的数据了, 需要利用这些数据处理具体的业务
			// 区分握手 和 其他请求
			if body.Seq == 0 {
				// 第一次握手
				if body.Name == HandshakeMsg {
					// 获取秘钥
					hs := &Handshake{}
					mapstructure.Decode(body.Msg, hs)
					if hs.Key != "" {
						// 存在秘钥
						c.SetProperty("secretKey", hs.Key)
					} else {
						// 秘钥不存在, 需要删除老的, 防止握手错误
						c.RemoveProperty("secretKey")
					}
					// 收到握手消息, 发送一个channel
					c.handshake = true
					// 发送消息
					c.handshakeChannel <- true
				} else {
					// 处理其他请求
					if c.onPush != nil {
						c.onPush(c, body)
					}
				}
			} else {
				// 其他请求(此时已经获取到信息)
				c.syncCtxLock.RLock()
				// 获取客户端id(每次序号会自增, 因此body.Seq就代表客户端id)
				ctx, ok := c.syncCtxMap[body.Seq]
				c.syncCtxLock.RUnlock()
				if ok {
					ctx.outChan <- body
				} else {
					log.Println("客户端id对应的syncCtx不存在")
				}
			}
		}
	}
	// 跳出循环说明中间结束或者出现问题, 同样执行关闭
	c.Close()
}

func (c *ClientConnect) SetOnPush(pushHook func(connect *ClientConnect, body *ResBody)) {
	// 设置onPush方法
	c.onPush = pushHook
}

func (c *ClientConnect) Close() {
	// 关闭连接
	_ = c.wsConnect.Close()
}

func (c *ClientConnect) Send(name string, msg interface{}) (*ResBody, error) {
	// 把请求发送给代理服务器 登录服务器
	c.syncCtxLock.Lock()
	c.Seq++
	seq := c.Seq
	sc := NewSyncCtx()
	// 构建一个request请求, 然后向登录服务器推送消息
	req := &ReqBody{Seq: seq, Name: name, Msg: msg}
	c.syncCtxMap[seq] = sc
	c.syncCtxLock.Unlock()
	res := &ResBody{Name: name, Seq: seq, Code: constant.OK}
	err := c.write(req)
	if err != nil {
		sc.cancel()
	} else {
		r := sc.waitMsg()
		if r == nil {
			res.Code = constant.ProxyConnectError
		} else {
			res = r
		}
	}

	c.syncCtxLock.Lock()
	delete(c.syncCtxMap, seq)
	c.syncCtxLock.Unlock()
	return res, nil
}

func NewClientConnect(wsConnect *websocket.Conn) *ClientConnect {
	return &ClientConnect{
		wsConnect:        wsConnect,
		handshakeChannel: make(chan bool), // 初始化channel
	}
}
