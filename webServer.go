// 账号服务入口
package main

import (
	"log"
	"ms_sg_back/config"
	"ms_sg_back/server/web"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	host := config.File.MustValue("web_server", "host", "127.0.0.1")
	port := config.File.MustValue("web_server", "port", "8003")

	router := gin.Default()
	// 路由配置
	web.Init(router)
	s := &http.Server{
		Addr:           host + ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	log.Println(err)
}
