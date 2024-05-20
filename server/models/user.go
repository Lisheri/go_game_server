package models

import "time"

// 实体类
type User struct {
	UId      int       `xorm:"uid pk autoincr"`
	Username string    `xorm:"username" validate:"min=4,max=20,regexp=^[a-zA-Z0-9_]*$"`
	PassCode string    `xorm:"passcode"`
	Passwd   string    `xorm:"passwd"`
	Hardware string    `xorm:"hardware"`
	Status   int       `xorm:"status"`
	Ctime    time.Time `xorm:"ctime"`
	Mtime    time.Time `xorm:"mtime"`
	IsOnline bool      `xorm:"-"`
}

func (*User) TableName() string {
	// 表名
	return "user"
}
