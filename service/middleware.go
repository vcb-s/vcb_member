package service

import (
	"github.com/gin-gonic/gin"
	"vcb_member/helper"
)

// AuthMiddleware 登录检查以及token重签
func AuthMiddleware(c *gin.Context) {
	var j JSONData
	originToken := []byte(c.GetHeader("X-Token"))
	originRefreshToken := []byte(c.GetHeader("X-RefreshToken"))

	if cap(originToken) == 0 {
		j.Unauthorized(c)
		return
	}

	shouldReGen := false
	uid, err := helper.CheckToken(originToken)
	if err != nil {
		// 不是过期错误就不需要检查refreshToken
		if err.Error() != helper.ErrorExpired || cap(originRefreshToken) == 0 {
			j.FailAuth(c)
			return
		}
		// 过期token检查一下是否可以重发
		shouldReGen, uid, err = helper.CheckRefreshToken(originRefreshToken)
		if err != nil {
			j.FailAuth(c)
			return
		}
	}
	c.Request.Header.Set("uid", uid)

	c.Next()

	/** token重签间隔：每次刷新 */
	/** refreshToken重签间隔：常规token的过期时间 * 2 */

	newToken, err := helper.GenToken(uid)
	if err != nil {
		// 这里重签失败不能抛出到响应体中，因为已经abort了
		// j.ServerError(c, err)
		c.Writer.Header().Add("X-Error-Report", err.Error())
		return
	}
	c.Writer.Header().Set("token", newToken)

	var newRefreshToken string
	if cap(originRefreshToken) > 0 {
		if shouldReGen {
			newRefreshToken, err = helper.ReGenRefreshToken(originRefreshToken)
			if err != nil {
				// 这里重签失败不能抛出到响应体中，因为已经abort了
				// j.ServerError(c, err)
				c.Writer.Header().Add("X-Error-Report", err.Error())
				return
			}
		}
	} else {
		newRefreshToken, err = helper.GenRefreshToken(uid)
		if err != nil {
			// 这里重签失败不能抛出到响应体中，因为已经abort了
			// j.ServerError(c, err)
			c.Writer.Header().Add("X-Error-Report", err.Error())
			return
		}
	}
	c.Writer.Header().Set("refreshToken", newRefreshToken)
}
