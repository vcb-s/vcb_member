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
		root.GET("/user/list", service.UserList)
		root.GET("/user-card/list", service.UserCardList)
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

		// 创建用户
		adminWithAuth.POST("/user/info", service.PersonInfo)

		// 创建用户
		adminWithAuth.POST("/user/create", service.CreateUser)
		// 更新用户
		adminWithAuth.POST("/user/update", service.UpdateUser)
		// 拉组
		adminWithAuth.POST("/user/group/add", service.AddGroup)
		// 踢出
		adminWithAuth.POST("/user/kickoff", service.KickOff)

		// 创建用户卡片
		adminWithAuth.POST("/user-card/create", service.CreateUserCard)
		// 更新用户卡片
		adminWithAuth.POST("/user-card/update", service.UpdateUserCard)

		// 重置密码
		adminWithAuth.POST("/password/reset", service.ResetPass)
		// 绑定主站账号
		adminWithAuth.POST("/bind-wp/create", service.CreateBindForWP)
		// 解绑主站账号
		adminWithAuth.POST("/bind-wp/delete", service.DeleteWPBind)

	}
}
