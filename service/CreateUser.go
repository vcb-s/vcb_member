package service

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"vcb_member/helper"
	"vcb_member/models"
)

type createUserReq struct {
	Group []string `json:"group" form:"group" gorm:"column:group"`
}

// CreateUser 创建新的用户
func CreateUser(c *gin.Context) {
	var (
		j                JSONData
		req              createUserReq
		userInAuth       models.User
		userToCreate     models.User
		userCardToCreate models.UserCard
	)

	userInAuth.ID = c.Request.Header.Get("uid")
	userToCreate.ID = helper.GenID()
	// 两次GenID之间必须存在一定延时，避免序号连续
	// userCardToCreate.ID = helper.GenID()

	password := helper.GenCode()
	passwordHash, err := helper.CalcPassHash(password)
	if err != nil {
		j.ServerError(c, err)
		return
	}

	userToCreate.Password = passwordHash

	userCardToCreate.ID = helper.GenID()
	userCardToCreate.UID = userToCreate.ID
	userCardToCreate.Hide = 1

	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	// 查询权限
	if err := models.GetDBHelper().First(&userInAuth, "`id` = ?", userInAuth.ID).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	// 不是管理员
	if !userInAuth.IsAdmin() {
		j.Message = "只有管理员允许创建用户"
		j.FailAuth(c)
		return
	}

	err = models.GetDBHelper().Transaction(func(db *gorm.DB) error {
		// 更新该用户的组别信息
		userToCreate.Group = strings.Join(req.Group, ",")
		userToCreate.Nickname = "新用户"
		userCardToCreate.Nickname = userToCreate.Nickname
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

	j.Data = map[string]interface{}{
		"cardID": userCardToCreate.ID,
		"UID":    userToCreate.ID,
		"pass":   password,
	}

	j.ResponseOK(c)
	return
}