package proto

// 登录响应数据
type LoginRes struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Session  string `json:"session"`
	UId      int    `json:"uid"`
}

// 请求参数
type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
	Hardware string `json:"hardware"`
}
