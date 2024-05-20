package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 利用中间件处理跨域

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		// 添加Access-Control-Allow-Origin 响应头
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		// 携带cookie
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 继续处理请求
		c.Next()
	}
}
