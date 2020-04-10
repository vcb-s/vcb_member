package service

import (
	"vcb_member/helper"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 登录检查以及token重签
func AuthMiddleware(c *gin.Context) {
	var j JSONData
	originToken := []byte(c.GetHeader("X-Token"))

	if cap(originToken) == 0 {
		j.Unauthorized(c)
		return
	}

	uid, err := helper.CheckToken(originToken)
	if err != nil {
		// 不是过期错误就不需要检查refreshToken
		if err != nil {
			j.FailAuth(c)
			return
		}
	}
	c.Request.Header.Set("uid", uid)

	c.Next()

	if err == nil && !c.IsAborted() {
		newToken, err := helper.GenToken(uid)
		if err != nil {
			j.ServerError(c, err)
			return
		}
		c.Writer.Header().Set("token", newToken)
	}

}
