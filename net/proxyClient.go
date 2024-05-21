package net

import (
	"errors"
	"time"

	"github.com/gorilla/websocket"
)

// 结构体本身保持私有
type ProxyClient struct {
	 // 代理地址
	proxy string
	// 连接通道(客户端专用)
	connect *ClientConnect
}

func NewProxyClient(proxySrc string) *ProxyClient {
	return &ProxyClient{
		proxy: proxySrc,
	}
}

func (c *ProxyClient) Connect() error {
	// 连接对应的 ws 服务端
	// 通过 Dialer 连接 ws 服务器
	var dialer = websocket.Dialer {
		Subprotocols: []string{"p1", "p2"},
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		HandshakeTimeout: 30 * time.Second,
	}

	// 返回一个ws客户端连接以及错误信息等
	ws, _, err := dialer.Dial(c.proxy, nil)
	if err == nil {
		c.connect = NewClientConnect(ws)
		if !c.connect.Start() {
			return errors.New("握手失败")
		}
	}
	return err
}

func (c *ProxyClient) SetProperty(key string, date interface{}) {
	if c.connect != nil {
		c.connect.SetProperty(key, date)
	}
}

func (c *ProxyClient) SetOnPush(push func(connect *ClientConnect, body *ResBody)) {
	if c.connect != nil {
		c.connect.SetOnPush(push)
	}
}

func (c *ProxyClient) Send(name string, msg interface{}) (*ResBody, error) {
	if c.connect != nil {
		// 交给真实的连接推送消息
		return c.connect.Send(name, msg)
	}
	return nil, errors.New("连接未找到")
}
