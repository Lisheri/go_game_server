package web

import (
	"ms_sg_back/db"
	"ms_sg_back/server/web/controller"
	"ms_sg_back/server/web/middleware"

	"github.com/gin-gonic/gin"
)

// go模块入口
func Init(router *gin.Engine) {
	db.TestAndInitDB()
	// 路由初始化
	initRouter(router)
}

func initRouter(router *gin.Engine) {
	// account := router.Group("/account")
	router.Use(middleware.Cors())
	router.POST("/account/register", controller.DefaultAccountController.Register)
}
