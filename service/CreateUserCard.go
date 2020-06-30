package service

import (
	"github.com/gin-gonic/gin"

	"vcb_member/helper"
	"vcb_member/models"
)

type createUserCardReq struct {
	models.UserCard
	ID  string `json:"-" form:"-" gorm:"PRIMARY_KEY;column:id"`
	UID string `json:"-" form:"-" gorm:"column:uid"`
}

// CreateUserCard 创建新的用户卡片
func CreateUserCard(c *gin.Context) {
	var (
		j                JSONData
		userToBind       models.User
		userCardToCreate models.UserCard
	)

	UID := c.Request.Header.Get("uid")

	userToBind.ID = UID

	userCardToCreate.ID = helper.GenID()
	userCardToCreate.Hide = 1

	if err := models.GetDBHelper().Model(&userToBind).First(&userToBind, "id = ?", UID).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	userCardToCreate.UID = UID
	userCardToCreate.Avast = userToBind.Avast
	userCardToCreate.Nickname = userToBind.Nickname
	userCardToCreate.Group = userToBind.Group

	if err := models.GetDBHelper().Model(&userCardToCreate).Create(&userCardToCreate).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	j.Data = map[string]interface{}{
		"ID": userCardToCreate.ID,
	}

	j.ResponseOK(c)
	return
}
