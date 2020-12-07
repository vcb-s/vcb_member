package service

import (
	"github.com/gin-gonic/gin"

	"vcb_member/models"
)

// GroupList 组别列表
func GroupList(c *gin.Context) {
	var (
		j JSONData
	)

	userGroupList := make([]models.UserCardGroup, 0)

	total := int64(0)

	if err := models.GetDBHelper().Find(&userGroupList).Count(&total).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	j.Data = map[string]interface{}{"res": userGroupList, "total": total}
	j.ResponseOK(c)
	return
}
