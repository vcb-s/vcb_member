package service

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"vcb_member/helper"
	"vcb_member/models"
)

type createUserReq struct {
	Group string `json:"group" form:"group" gorm:"column:group"`
}

// TableName 指示 User 表名
func (m createUserReq) TableName() string {
	return models.User{}.TableName()
}

// CreateUser 创建新的用户
func CreateUser(c *gin.Context) {
	var (
		j                JSONData
		req              updateUserReq
		userInAuth       models.User
		userToCreate     models.User
		userCardToCreate models.UserCard
	)

	userInAuth.UID = c.Request.Header.Get("uid")
	userToCreate.UID = helper.GenID()

	userToCreate.Password = helper.GenCode()
	userCardToCreate.ID = helper.GenID()
	userCardToCreate.UID = userToCreate.UID
	userCardToCreate.Hide = 1

	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	// 查询权限
	if err := models.GetDBHelper().First(&userCardToCreate, "`id` = ?", userCardToCreate.UID).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	// 不是管理员
	if !userInAuth.IsAdmin() {
		j.Message = "只有管理员允许创建用户"
		j.FailAuth(c)
		return
	}

	err := models.GetDBHelper().Transaction(func(db *gorm.DB) error {
		// 更新该用户的组别信息
		groups := []string{}
		for _, group := range strings.Split(req.Group, ",") {
			groups = append(groups, group)
		}

		userToCreate.Group = strings.Join(groups, ",")
		userCardToCreate.Group = userToCreate.Group

		// 写入
		if err := db.Model(&userToCreate).Create(&userToCreate).Error; err != nil {
			return err
		}
		if err := db.Model(&userCardToCreate).Create(&userCardToCreate).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		j.ServerError(c, err)
		return
	}

	j.ResponseOK(c)
	return
}
