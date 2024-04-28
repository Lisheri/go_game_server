package controller

import (
	"log"
	"ms_sg_back/constant"
	"ms_sg_back/db"
	"ms_sg_back/net"
	"ms_sg_back/server/login/model"
	"ms_sg_back/server/login/proto"
	"ms_sg_back/utils"
	"time"

	"github.com/mitchellh/mapstructure"
)

var DefaultAccount = &Account{}

type Account struct {
}

func (a *Account) Router(router *net.Router) {
	group := router.Group("account")
	group.AddRouter("login", a.login)
}

// 这个login, 就是传递的handler
func (a *Account) login(req *net.WsMsgReq, res *net.WsMsgRes) {
	// 1. 用户名 2. 密码 3. 硬件id
	// 2. 根据用户名查询 user表 获取数据
	// 3. 密码比对, 若密码正确则登录成功
	// 4. 保存用户登录记录
	// 5. 保存用户的最后一次登录信息(便于跟踪用户相关信息，活跃度等)
	// 6. 客户端需要一个 token, 利用jwt生成一个加密字符串(token)(本周上是一个加密算法)
	// 7. 后续所有请求, 均需要带上token判断是否合法
	// * 这里需要将请求的msg转换为对应的loginReq和loginRes, 需要借助第三方库 mapstructure
	// * mapstructure用于将通用的map[string]interface{}解码到对应的 Go 结构体中，或者执行相反的操作。
	loginReq := &proto.LoginReq{}
	loginRes := &proto.LoginRes{}

	// 将interface转换为结构体
	mapstructure.Decode(req.Body.Msg, loginReq)
	// 获取到用户名, 然后查询数据
	username := loginReq.Username
	user := &model.User{}
	ok, err := db.Engine.Table(user).Where("username=?", username).Get(user)
	if err != nil {
		log.Println("用户表查询出错", err)
	}
	if !ok {
		// 没有查处数据来, 说明用户不存在
		res.Body.Code = constant.UserNotExist
		return
	}
	// 用户存在
	pwd := utils.Password(loginReq.Password, user.PassCode)
	if pwd != user.Passwd {
		// 说明没有验证通过, 用户名密码错误
		res.Body.Code = constant.PwdIncorrect
		return
	}
	// 密码正确, 验证通过, 登录成功
	// A B C 组成 jwt, A为加密算法， B为放入数据, C根据秘钥+A和B生成的加密字符串, 只要秘钥不丢, token就是安全的
	// 1. 基于uid生成Token(jwt), 7天过期
	token, _ := utils.Award(user.UId)

	res.Body.Code = constant.OK
	loginRes.UId = user.UId
	loginRes.Username = username
	loginRes.Session = token
	loginRes.Password = ""
	res.Body.Msg = loginRes

	// 保存用户登录记录
	ul := &model.LoginHistory{
		UId: user.UId, CTime: time.Now(), Ip: loginReq.Ip,
		Hardware: loginReq.Hardware, State: model.Login,
	}

	// 插入数据
	db.Engine.Table(ul).Insert(ul)

	// 最后一次登录状态记录
	ll := &model.LoginLast{}
	ok, _ = db.Engine.Table(ll).Where("uid=?", user.UId).Get(ll)
	if ok {
		// 有数据, 更新
		ll.IsLogout = 0
		ll.Ip = loginReq.Ip
		ll.LoginTime = time.Now()
		ll.Session = token
		ll.Hardware = loginReq.Hardware
		ll.UId = user.UId
		db.Engine.Table(ll).Insert(ll)
	}

	// 缓存当前用户和当前ws的连接(用户其他地方登录需要断开当前连接)
}
