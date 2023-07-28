package router

import (
	"github.com/gin-gonic/gin"

	"vcb_member/helper"
	"vcb_member/service"
)

// AuthMiddleware 登录检查以及token重签
func AuthMiddleware(c *gin.Context) {
	var j service.JSONData

	tokenStr := c.GetHeader("X-Token")
	if tokenStr == "" {
		j.Unauthorized(c)
		return
	}

	uid, err := helper.CheckToken(tokenStr)
	if err != nil {
		j.Message = err.Error()
		j.FailAuth(c)
		return
	}
	c.Request.Header.Set("uid", uid)

	c.Next()
}
