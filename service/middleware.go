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
		j.Message = err.Error()
		j.FailAuth(c)
		return
	}
	c.Request.Header.Set("uid", uid)

	c.Next()
}
