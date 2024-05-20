package service

import (
	"log"
	"ms_sg_back/constant"
	"ms_sg_back/db"
	"ms_sg_back/server/models"
	"ms_sg_back/server/web/common"
	"ms_sg_back/server/web/model"
	"ms_sg_back/utils"
	"time"
)

var DefaultAccountService = &AccountService{}

type AccountService struct {
}

func (account AccountService) Register(req *model.RegisterReq) error {
	// 1. 解析参数
	// 2. 根据用户名 查询数据局库是否重复, 未重复则进行注册
	// 3. 回调成功结果
	username := req.Username
	user := &models.User{}
	ok, err := db.Engine.Table(user).Where("username=?", username).Get(user)
	if err != nil {
		log.Println("注册查询失败", err)
		return common.New(constant.DBError, "数据库异常")
	}

	if ok {
		return common.New(constant.UserExist, "用户已存在")
	} else {
		// 注册用户
		user.Mtime = time.Now()
		user.Ctime = time.Now()

		user.Username = req.Username
		user.PassCode = utils.RandSeq(6)
		user.Passwd = utils.Password(req.Password, user.PassCode)
		user.Hardware = req.Hardware
		// 插入数据表
		_, err := db.Engine.Table(user).Insert(user)
		if err != nil {
			log.Println("插入用户信息失败", err)
			return common.New(constant.DBError, "数据库异常")
		}
		return nil
	}
}
