package service

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"vcb_member/helper"
	"vcb_member/models"
)

type createBindForWPReq = loginWithWPCodeReq

// CreateBindForWP 添加主站绑定
func CreateBindForWP(c *gin.Context) {
	var (
		j          JSONData
		req        createBindForWPReq
		userToBind models.UserAssociation
	)
	uidToBind := c.Request.Header.Get("uid")
	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	// 查询用户是否有同类型绑定，不允许重复
	err := models.GetDBHelper().Where("type = ? AND uid = ?", models.UserAssociationTypeWP, uidToBind).First(&userToBind).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			j.ServerError(c, errors.New("你已绑定其他账号"))
			return
		}
		j.ServerError(c, err)
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

	// 检查主站ID是否已经跟其他账号绑定过
	err = models.GetDBHelper().Where("type = ? AND association = ?", models.UserAssociationTypeWP, strconv.Itoa(userInWP.ID)).First(&userToBind).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	// 不允许重复绑定
	if err == nil {
		j.Message = "该主站账号已被绑定"
		j.FailAuth(c)
		return
	}

	userToBind.ID = helper.GenID()
	userToBind.UID = uidToBind
	userToBind.AuthCode = strconv.Itoa(userInWP.ID)
	userToBind.Type = models.UserAssociationTypeWP

	// 没绑定过的就添加一条绑定
	err = models.GetDBHelper().Create(&userToBind).Error
	if err != nil {
		j.ServerError(c, err)
		return
	}

	j.ResponseOK(c)
	return
}
