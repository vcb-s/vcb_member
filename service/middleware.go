package service

import (
	"github.com/gin-gonic/gin"
	"vcb_member/helper"
)

var sercertUIDHeaderKey string

func init() {
	sercertUIDHeaderKey = helper.GenID()
}

// AuthMiddleware 登录检查以及token重签
func AuthMiddleware(c *gin.Context) {
	var j JSONData
	originToken := []byte(c.GetHeader("token"))
	originRefreshToken := []byte(c.GetHeader("refreshToken"))

	if cap(originToken) == 0 {
		j.Unauthorized(c)
		return
	}

	uid, err := helper.CheckToken(originToken)
	if err != nil {
		// 不是过期错误就不需要检查refreshToken
		if err.Error() != helper.ErrorExpired || cap(originRefreshToken) == 0 {
			j.FailAuth(c)
			return
		}
		// 过期token检查一下是否可以重发
		uid, err = helper.CheckRefreshToken(originRefreshToken)
		if err != nil {
			j.FailAuth(c)
			return
		}
	}
	c.Request.Header.Set(sercertUIDHeaderKey, uid)

	c.Next()

	/** token重签间隔：每次刷新 */
	/** @TODO refreshToken重签间隔：常规token的过期时间 * 2 */

	var newRefreshToken string
	newToken, err := helper.GenToken(uid)
	if err != nil {
		j.ServerError(c)
		return
	}
	if cap(originRefreshToken) > 0 {
		newRefreshToken, err = helper.ReGenRefreshToken(originRefreshToken)
		if err != nil {
			j.ServerError(c)
			return
		}
	} else {
		newRefreshToken, err = helper.GenRefreshToken(uid)
		if err != nil {
			j.ServerError(c)
			return
		}
	}

	c.Writer.Header().Set("token", newToken)
	c.Writer.Header().Set("refreshToken", newRefreshToken)
}
