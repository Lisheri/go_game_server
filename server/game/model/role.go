package model

// 进入游戏请求参数设计
type EnterServerReq struct {
	Session string `json:"session"`
}

type EnterServerRes struct {
	Role         Role         `json:"role"`     // 角色信息
	RoleResource RoleResource `json:"role_res"` // 角色资源
	Time         int64        `json:"time"`
	Token        string       `json:"token"` // 根据角色id生成一个token返回给客户端使用, 用于认证
}

type Role struct {
	RId      int    `json:"rid"`
	UId      int    `json:"uid"`
	NickName string `json:"nickName"`
	Sex      int8   `json:"sex"`
	Balance  int    `json:"balance"`
	HeadId   int16  `json:"headId"`  // 头像
	Profile  string `json:"profile"` // 属性
}

// 角色固有的资源属性
type RoleResource struct {
	Wood          int `json:"wood"`
	Iron          int `json:"iron"`
	Stone         int `json:"stone"`
	Grain         int `json:"grain"`      // 粮食
	Gold          int `json:"gold"`       // 金币
	Decree        int `json:"decree"`     // 令牌
	WoodYield     int `json:"wood_yield"` // 增长量
	IronYield     int `json:"iron_yield"`
	StoneYield    int `json:"stone_yield"`
	GrainYield    int `json:"grain_yield"`
	GoldYield     int `json:"gold_yield"`
	DepotCapacity int `json:"depot_capacity"` // 仓库容量

}
