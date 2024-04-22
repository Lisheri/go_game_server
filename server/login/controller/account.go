package controller

import (
	"ms_sg_back/net"
	"ms_sg_back/server/login/proto"
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
	res.Body.Code = 0
	loginRes := &proto.LoginRes{}
	loginRes.UId = 1
	loginRes.Username = "admin"
	loginRes.Session = "as"
	loginRes.Password = ""
	res.Body.Msg = loginRes
}
