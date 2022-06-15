package service

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"vcb_member/helper"
	"vcb_member/models"
)

type loginWithWPCodeReq struct {
	Code string `json:"code" form:"code" binding:"required"`
}

// LoginFromWP 绑定登录
func LoginFromWP(c *gin.Context) {
	var (
		j           JSONData
		req         loginWithWPCodeReq
		association models.UserAssociation
	)
	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	// 根据 Authorization code 换取 AccessToken
	accessToken, err := helper.GetAccessTokenFromCode(req.Code)
	if err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	// 根据accessToken换取主站ID
	userInWP, err := helper.GetUserInfoFromAccesstoken(accessToken)
	if err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	// 根据主站ID在第三方绑定表查找
	err = models.GetDBHelper().Where("type = ? AND association = ?", models.UserAssociationTypeWP, strconv.Itoa(userInWP.ID)).First(&association).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			j.Message = "没有找到用户"
			j.FailAuth(c)
			return
		}
		j.ServerError(c, err)
		return
	}

	// 找到了就按照UID签发
	token, err := helper.GenToken(association.UID)
	if err != nil {
		j.ServerError(c, err)
		return
	}
	c.Writer.Header().Set("token", token)

	j.ResponseOK(c)
}
