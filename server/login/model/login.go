package model

import "time"

// 两个状态
const (
	Login = iota // 登录
	Logout // 登出
)

// 登录历史
type LoginHistory struct {
	Id       int       `xorm:"id pk autoincr"` // id自增
	UId      int       `xorm:"uid"`
	CTime    time.Time `xorm:"ctime"`
	Ip       string    `xorm:"ip"`
	State    int8      `xorm:"state"`
	Hardware string    `xorm:"hardware"`
}

// 上一次登录
type LoginLast struct {
	Id         int       `xorm:"id pk autoincr"`
	UId        int       `xorm:"uid"`
	LoginTime  time.Time `xorm:"login_time"`
	LogoutTime time.Time `xorm:"logout_time"`
	Ip         string    `xorm:"ip"`
	Session    string    `xorm:"session"`
	IsLogout   int8      `xorm:"is_logout"`
	Hardware   string    `xorm:"hardware"`
}

// xorm中可以自行指定表明, 结构体拥有 TableName() string 的成员方法, 那么此方法的返回值即是该结构体对应的数据库表名

func (*LoginHistory) TableName() string {
	return "login_history"
}

func (*LoginLast) TableName() string {
	return "login_last"
}

