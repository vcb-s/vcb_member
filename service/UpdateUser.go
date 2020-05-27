package service

import (
	"github.com/gin-gonic/gin"

	"vcb_member/models"
)

type updateUserReq struct {
	UID      string `json:"id" form:"id" gorm:"PRIMARY_KEY;column:id" binding:"required"`
	Admin    string `json:"admin" form:"admin" gorm:"column:admin"`
	Ban      int8   `json:"ban" form:"ban" gorm:"column:ban"`
	Avast    string `json:"avast" form:"avast" gorm:"column:avast"`
	Nickname string `json:"nickname" form:"nickname" gorm:"column:nickname"`
	Group    string `json:"group" form:"group" gorm:"column:group"`
}

// TableName 指示 User 表名
func (m updateUserReq) TableName() string {
	return models.User{}.TableName()
}

// UpdateUser 修改用户信息
func UpdateUser(c *gin.Context) {
	var (
		j            JSONData
		req          updateUserReq
		userToUpdate models.User
		userInAuth   models.User = models.User{}
	)

	userInAuth.UID = c.Request.Header.Get("uid")

	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	// 查询所属UID
	if err := models.GetDBHelper().Where("`id` = ?", req.UID).First(&userToUpdate).Error; err != nil {
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

	updateBuilder := models.GetDBHelper().Model(&req)

	// 修改键值
	if err := updateBuilder.Updates(&req).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	j.ResponseOK(c)
	return
}
