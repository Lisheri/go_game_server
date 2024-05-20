package controller

import (
	"fmt"
	"log"
	"ms_sg_back/constant"
	"ms_sg_back/server/web/common"
	"ms_sg_back/server/web/model"
	"ms_sg_back/server/web/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var DefaultAccountController = &AccountController{}

type AccountController struct {
}

func (account *AccountController) Register(ctx *gin.Context) {
	fmt.Println("register")
	// 1. 解析参数
	// 2. 根据用户名 查询数据局库是否重复, 未重复则进行注册
	// 3. 回调成功结果
	req := &model.RegisterReq{}

	// 绑定请求
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Println("参数格式不合法", err)
		ctx.JSON(http.StatusOK, common.Error(constant.InvalidParam, "参数不合法"))
		return
	}

	// 一般web服务, 错误格式会自定义
	err = service.DefaultAccountService.Register(req)
	if err != nil {
		// 产生错误
		log.Println("注册业务出错", err)
		// 转换为自定义的Error, 前面Code是自定义的, 但是后面的Error自带一个, 只是我们带有了自定义实现, 因此无需对error进行转换也可以使用
		ctx.JSON(http.StatusOK, common.Error(err.(*common.MyError).Code(), err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, common.Success(constant.OK, nil, "注册成功"))
}
