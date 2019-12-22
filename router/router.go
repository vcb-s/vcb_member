package router

import (
	"vcb_member/service"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	Router = gin.Default()
	// Router.GET("/check", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{"code": 0, "msg": "健康检测", "data": nil})
	// })

	// public := Router.Group("/admin")
	// public.GET("/", service.list)

	admin := Router.Group("/admin")
	// 登录验证中间件
	admin.Use(service.AuthMiddleware)
	// 是否登录
	// admin.GET("/isLogin", service.IsLogin)
}
