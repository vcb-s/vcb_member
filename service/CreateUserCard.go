package service

import (
	"github.com/gin-gonic/gin"

	"vcb_member/helper"
	"vcb_member/models"
)

// createUserCard 创建新的用户卡片
func createUserCard(c *gin.Context) {
	var (
		j                JSONData
		userCardToCreate models.UserCard
	)

	UID := c.Request.Header.Get("uid")

	userCardToCreate.ID = helper.GenID()
	userCardToCreate.UID = UID
	userCardToCreate.Hide = 1

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
