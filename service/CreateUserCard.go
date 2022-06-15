package service

import (
	"github.com/gin-gonic/gin"

	"vcb_member/helper"
	"vcb_member/models"
)

type createUserCardReq struct {
	models.UserCard
	ID  string `json:"-" form:"-" gorm:"primaryKey;column:id"`
	UID string `json:"-" form:"-" gorm:"column:uid"`
}

// CreateUserCard 创建新的用户卡片
func CreateUserCard(c *gin.Context) {
	var (
		j          JSONData
		req        createUserCardReq
		userToBind models.User
	)

	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	userCardToCreate := req

	UID := c.Request.Header.Get("uid")

	userToBind.ID = UID

	if err := models.GetDBHelper().Model(&userToBind).First(&userToBind, "id = ?", UID).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	userCardToCreate.ID = helper.GenID()
	userCardToCreate.UID = UID

	if err := models.GetDBHelper().Model(&userCardToCreate).Create(&userCardToCreate).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	j.Data = map[string]interface{}{
		"ID": userCardToCreate.ID,
	}

	j.ResponseOK(c)
}
