package service

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"vcb_member/models"
)

type personInfoReq struct {
	UID string `json:"uid,omitempty" form:"uid,omitempty"`
}

// PersonInfo 个人卡片列表
func PersonInfo(c *gin.Context) {
	var (
		j   JSONData
		req personInfoReq

		uidInAuth     = c.Request.Header.Get("uid")
		userInAuth    models.User
		userInRequest models.User

		userCardList = make([]models.UserCard, 0)
		userList     = make([]models.User, 0)
	)

	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	uidInRequest := req.UID
	if uidInRequest == "" {
		uidInRequest = uidInAuth
	}

	if err := models.GetDBHelper().First(&userInAuth, "`id` = ?", uidInAuth).Error; err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	if err := models.GetDBHelper().First(&userInRequest, "`id` = ?", uidInRequest).Error; err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	if !userInAuth.CanManagePerson(userInRequest) {
		j.Message = "你无权获取该用户信息"
		j.BadRequest(c)
		return
	}

	cardTotal := int64(0)
	userTotal := int64(0)
	{
		var sqlBuilder = models.GetDBHelper().Where("`uid` = ?", uidInRequest)

		err := sqlBuilder.Find(&userCardList).Count(&cardTotal).Error
		if err != nil {
			j.Message = err.Error()
			j.BadRequest(c)
			return
		}
	}

	if userInRequest.IsAdmin() {
		var sqlBuilder = models.GetDBHelper()
		groupsBelongUser := strings.Split(userInRequest.Admin, ",")

		for _, group := range groupsBelongUser {
			sqlBuilder = sqlBuilder.Or("`group` like ?", fmt.Sprintf("%%%s%%", group))
		}

		if err := sqlBuilder.Find(&userList).Count(&userTotal).Error; err != nil {
			j.Message = err.Error()
			j.BadRequest(c)
			return
		}
	}

	j.Data = map[string]interface{}{
		"info": userInRequest,
		"cards": map[string]interface{}{
			"total": cardTotal,
			"res":   userCardList,
		},
		"users": map[string]interface{}{
			"total": userTotal,
			"res":   userList,
		},
	}
	j.ResponseOK(c)
	return
}
