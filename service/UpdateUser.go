package service

import (
	"github.com/gin-gonic/gin"

	"vcb_member/models"
)

type updateUserReq struct {
	UID      string `json:"id" form:"id" gorm:"primaryKey;column:id" binding:"required"`
	Admin    string `json:"admin" form:"admin" gorm:"column:admin"`
	Ban      int8   `json:"ban" form:"ban" gorm:"column:ban"`
	Avast    string `json:"avast" form:"avast" gorm:"column:avast"`
	Nickname string `json:"nickname" form:"nickname" gorm:"column:nickname"`
	// 组别修改走另一套逻辑，
	// 因为组别有转入逻辑，有踢出逻辑，这些的权限判断都比较独立
	// Group    []string `json:"group" form:"group" gorm:"column:group"`
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

	userInAuth.ID = c.Request.Header.Get("uid")

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

	// 查询授权用户
	if err := models.GetDBHelper().First(&userInAuth, "`id` = ?", userInAuth.ID).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	if !userInAuth.CanManagePerson(userToUpdate) {
		j.Message = "无权修改该用户"
		j.BadRequest(c)
		return
	}

	updateBuilder := models.GetDBHelper().Model(&req)

	// 修改键值
	if err := updateBuilder.Updates(&req).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	j.ResponseOK(c)
}
