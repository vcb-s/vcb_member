package router

import (
	"github.com/gin-gonic/gin"

	"vcb_member/service"
)

// Router 全局路由处理
var Router *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	Router = gin.Default()
	root := Router.Group("vcbs_member_api")
	{
		root.GET("/user/list", service.UserList)
		root.GET("/group/list", service.GroupList)
	}

	// 后台
	admin := root.Group("/admin")
	// 登录验证中间件
	admin.Use(service.AuthMiddleware)
	{
		admin.GET("/login", service.Login)
	}

	// 带登录验证的部分
	adminWithAuth := root.Group("/admin")
	adminWithAuth.Use(service.AuthMiddleware)
	{
		adminWithAuth.GET("/resetPassForSA", service.ResetPassForSuperAdmin)
	}

}
