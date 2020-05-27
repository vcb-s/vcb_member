package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"vcb_member/helper"
	"vcb_member/models"
)

type resetPassReq struct {
	UID         string `json:"uid" form:"uid"`
	Current     string `json:"current" form:"current"`
	NewPassword string `json:"new" form:"new"`
}

// ResetPass 重设密码
func ResetPass(c *gin.Context) {
	var (
		j           JSONData
		req         resetPassReq
		userInAuth  models.User
		userToReset models.User
	)
	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	uidInAuth := c.Request.Header.Get("uid")
	if len(req.UID) == 0 {
		userToReset.UID = uidInAuth
	} else {
		userToReset.UID = req.UID
	}

	if err := models.GetDBHelper().First(&userToReset, "`id` = ?", userToReset.UID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			j.BadRequest(c)
			return
		}
		j.ServerError(c, err)
		return
	}

	if uidInAuth == userToReset.UID {
		userInAuth = userToReset
	} else {
		if err := models.GetDBHelper().First(&userInAuth, "`id` = ?", uidInAuth).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				j.BadRequest(c)
				return
			}
			j.ServerError(c, err)
			return
		}
		if !userInAuth.CanManagePerson(userToReset) {
			if !helper.CheckPassHash(req.Current, userInAuth.Password) {
				j.Message = "密码错误"
				j.FailAuth(c)
				return
			}
		}
	}

	if req.NewPassword == "" {
		req.NewPassword = helper.GenCode()
	}

	newPass, err := helper.CalcPassHash(req.NewPassword)
	if err != nil {
		j.ServerError(c, err)
		return
	}

	userToReset.Password = newPass

	result := models.GetDBHelper().Model(&userToReset).Updates(userToReset)
	if result.Error != nil {
		j.ServerError(c, err)
		return
	}
	if result.RowsAffected == 0 {
		j.ServerError(c, errors.New("用户不存在"))
		return
	}

	j.Data = map[string]interface{}{
		"newPass": req.NewPassword,
	}

	j.ResponseOK(c)
	return
}
