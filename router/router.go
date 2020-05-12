package router

import (
	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"vcb_member/service"
)

// Router 全局路由处理
var Router *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	Router = gin.Default()

	// config := cors.DefaultConfig()
	// config.AllowWildcard = true
	// config.AllowOrigins = []string{
	// 	"http://localhost:*",
	// }
	// config.AllowHeaders = []string{
	// 	"*",
	// }
	// Router.Use(cors.New(config))

	root := Router.Group("vcbs_member_api")
	// 前台
	{
		root.GET("/user/list", service.UserCardList)
		root.GET("/user/list/tiny", service.TinyUserCardList)
		root.GET("/group/list", service.GroupList)
	}

	// 后台
	{
		admin := root.Group("/admin")
		// 登录
		admin.POST("/login", service.Login)
		// 关联登录
		admin.POST("/loginWithWP", service.LoginFromWP)
	}

	// 带登录验证的部分
	{
		adminWithAuth := root.Group("/admin")
		adminWithAuth.Use(AuthMiddleware)

		// 查询用户信息及其卡片列表
		adminWithAuth.GET("/personInfo", service.PersonInfo)

		// 重置密码
		adminWithAuth.POST("/resetPass", service.ResetPass)
		// 绑定主站账号
		adminWithAuth.POST("/createWPBind", service.CreateBindForWP)
		// 解绑主站账号
		adminWithAuth.POST("/deleteWPBind", service.DeleteWPBind)
		// 修改自己或他人的信息
		adminWithAuth.POST("/updateUserCard", service.UpdateUserCard)
		// 踢出
		adminWithAuth.POST("/kickoff", service.KickOff)
	}
}
