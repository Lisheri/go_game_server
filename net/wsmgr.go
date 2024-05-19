package net

import "sync"

var Mgr = &WsMgr{
	userCache: make(map[int]WsConnect),
}

// 管理用户登录缓存
type WsMgr struct {
	uc sync.RWMutex
	// key为uid, value为ws链接实例
	userCache map[int]WsConnect
}

func (m *WsMgr) UserLogin(conn WsConnect, uid int, token string) {
	// 先上锁
	m.uc.Lock()
	// 异步释放
	// defer是Go语言提供的一种用于注册延迟调用的机制：
	// 目的是让函数或语句可以在当前函数执行完毕后（包括通过return正常结束或者panic导致的异常结束）执行, 用defer可以保证函数执行完毕后释放锁
	defer m.uc.Unlock()
	oldConnect := m.userCache[uid]
	if oldConnect != nil {
		// 存在用户登录
		if conn != oldConnect {
			// 通知旧的客户端有用户正在抢登
			oldConnect.Push("robLogin", nil)
		}
	}
	m.userCache[uid] = conn
	conn.SetProperty("uid", uid)
	conn.SetProperty("token", token)
}
