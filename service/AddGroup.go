package service

import (
	"strings"

	"github.com/gin-gonic/gin"

	"vcb_member/models"
)

type addGroupReq struct {
	UID   string   `json:"uid" form:"uid" binding:"required"`
	Group []string `json:"group" form:"group" binding:"required"`
}

// TableName 指示 User 表名
func (m addGroupReq) TableName() string {
	return models.User{}.TableName()
}

// AddGroup 拉某个人进组
func AddGroup(c *gin.Context) {
	var (
		j            JSONData
		req          addGroupReq
		userToPullIn models.User
		userInAuth   models.User
	)
	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	userInAuth.ID = c.Request.Header.Get("uid")
	userToPullIn.ID = req.UID

	if err := models.GetDBHelper().Where(userToPullIn).First(&userInAuth).Error; err != nil {
		j.ServerError(c, err)
		return
	}
	if err := models.GetDBHelper().Where(userInAuth).First(&userToPullIn).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	if (userInAuth.IsContainAllGroup(models.User{Group: strings.Join(req.Group, ",")})) {
		j.Message = "你只能把用户拉进你能管理的组内"
		j.BadRequest(c)
		return
	}

	// 更新该用户的组别信息
	nextGroups := strings.Split(userToPullIn.Group, ",")
	for _, group := range req.Group {
		if !strings.Contains(userToPullIn.Group, group) {
			nextGroups = append(nextGroups, group)
		}
	}

	userToPullIn.Group = strings.Join(nextGroups, ",")
	updateBuilder := models.GetDBHelper().Model(&req)

	// 修改键值
	if err := updateBuilder.Updates(&userToPullIn).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	j.ResponseOK(c)
	return
}
