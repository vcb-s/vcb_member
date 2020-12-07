package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"vcb_member/helper"
	"vcb_member/models"
)

type loginReq struct {
	UID      string `json:"uid" form:"uid" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// Login 登录
func Login(c *gin.Context) {
	var (
		j    JSONData
		req  loginReq
		user models.User
	)
	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	if err := models.GetDBHelper().First(&user, "`id` = ?", req.UID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			j.Message = "用户不存在"
			j.NotAcceptable(c)
			return
		}
		j.ServerError(c, err)
		return
	}

	// 如果用户组别为空或者被Ban了
	if user.Group == "" || user.IsBan() {
		j.Message = "该账号不允许登录"
		j.NotAcceptable(c)
		return
	}

	// 如果用户密码为空，返回提示找网站组设置初始密码
	if user.Password == "" {
		j.Message = "请先联系网站组设置登录数据"
		j.NotAcceptable(c)
		return
	}

	if !helper.CheckPassHash(req.Password, user.Password) {
		j.Message = "密码不正确"
		j.FailAuth(c)
		return
	}

	// 签发密钥
	token, err := helper.GenToken(user.ID)
	if err != nil {
		j.ServerError(c, err)
		return
	}

	c.Writer.Header().Set("X-Token", token)

	j.ResponseOK(c)
	return
}
