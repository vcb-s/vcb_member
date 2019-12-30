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
	{
		// 登录
		admin.GET("/login", service.Login)
		// 关联登录
		admin.GET("/loginWithWP", service.LoginWithWPCode)
	}

	// 带登录验证的部分
	adminWithAuth := root.Group("/admin")
	adminWithAuth.Use(service.AuthMiddleware)
	{
		// 重置自己的密码
		adminWithAuth.GET("/resetPass", service.ResetPass)
		// 重置他人的密码
		adminWithAuth.GET("/resetPassForSA", service.ResetPassForSuperAdmin)
		// @TODO 绑定主站账号
		// adminWithAuth.GET("/createWPBind", service.ResetPassForSuperAdmin)
	}
}
