package service

import (
	"github.com/gin-gonic/gin"

	"vcb_member/models"
)

// TinyUserCardList 用户列表（登录页面用）
func TinyUserCardList(c *gin.Context) {
	var (
		j            JSONData
		userCardList = make([]models.UserCard, 0)
	)

	var sqlBuilder = models.GetDBHelper().Select("`id`, `uid`, `avast`, `nickname`")

	total := 0

	err := sqlBuilder.Find(&userCardList).Count(&total).Error
	if err != nil {
		j.ServerError(c, err)
		return
	}

	j.Data = map[string]interface{}{"res": userCardList, "total": total}
	j.ResponseOK(c)
	return
}
