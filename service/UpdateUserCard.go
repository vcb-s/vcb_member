package service

import (
	"github.com/gin-gonic/gin"

	"vcb_member/models"
)

type updateUserReq struct {
	models.UserCard
	ID  string `json:"id" form:"id" gorm:"PRIMARY_KEY;column:id" binding:"required"`
	UID string `json:"-" form:"-" gorm:"column:uid"`
}

// UpdateUserCard 修改用户信息
func UpdateUserCard(c *gin.Context) {
	var (
		j            JSONData
		req          updateUserReq
		userToUpdate models.UserCard
		userInAuth   models.User = models.User{}
	)

	userInAuth.UID = c.Request.Header.Get("uid")

	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	// 查询所属UID
	if err := models.GetDBHelper().Where("`id` = ?", req.ID).First(&userToUpdate).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	// 查询权限
	if err := models.GetDBHelper().First(&userInAuth, "`id` = ?", userInAuth.UID).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	// 不是管理员且uid不匹配的话
	if userToUpdate.UID != userInAuth.UID && !userInAuth.IsAdmin() {
		j.Message = "不允许修改他人信息"
		j.FailAuth(c)
		return
	}

	updateBuilder := models.GetDBHelper().Where("id = ?", req.ID)

	// 修改键值
	if err := updateBuilder.Model(&req).Updates(&req).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	j.ResponseOK(c)
	return
}
