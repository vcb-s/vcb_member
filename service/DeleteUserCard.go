package service

import (
	"github.com/gin-gonic/gin"

	"vcb_member/models"
)

type deleteUserReq struct {
	ID string `json:"id" form:"id"`
}

// DeleteUserCard （软）删除用户
func DeleteUserCard(c *gin.Context) {
	var (
		j                      JSONData
		req                    deleteUserReq
		userInAuth             models.User
		userCardToDeleteBelong models.User
		userCardToDelete       models.UserCard
	)

	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	userInAuth.ID = c.Request.Header.Get("uid")
	userCardToDelete.ID = req.ID

	// 查询权限
	if err := models.GetDBHelper().First(&userInAuth, "`id` = ?", userInAuth.ID).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	// 查询卡片
	if err := models.GetDBHelper().First(&userCardToDelete, "`id` = ?", userCardToDelete.ID).Error; err != nil {
		j.ServerError(c, err)
		return
	}
	// 查询卡片所属人
	if err := models.GetDBHelper().First(&userCardToDeleteBelong, "`id` = ?", userCardToDelete.UID).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	// 不是管理员
	if !userInAuth.CanManagePerson(userCardToDeleteBelong) {
		j.Message = "只有管理员允许创建用户"
		j.FailAuth(c)
		return
	}

	// 删除卡片
	if err := models.GetDBHelper().Delete(&userCardToDelete, "`id` = ?", userCardToDelete.ID).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	j.Data = map[string]interface{}{
		"id": userCardToDelete.ID,
	}

	j.ResponseOK(c)
}
