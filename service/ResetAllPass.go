package service

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"vcb_member/helper"
	"vcb_member/models"
)

type resetAllPassReq struct {
	NewPassword string `json:"new" form:"new"`
}

type miniUser struct {
	ID            string `json:"id" form:"id" gorm:"PRIMARY_KEY;column:id"`
	Password      string `json:"-" form:"-" gorm:"column:pass"`
	PasswordInStr string `json:"pass" form:"pass" gorm:"-"`
}

// TableName 指示 User 表名
func (m miniUser) TableName() string {
	return models.User{}.TableName()
}

// ResetAllPass 重设密码
func ResetAllPass(c *gin.Context) {
	var (
		j          JSONData
		req        resetAllPassReq
		userInAuth models.User
		// userToReset models.User
	)
	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	uidInAuth := c.Request.Header.Get("uid")

	if err := models.GetDBHelper().First(&userInAuth, "`id` = ?", uidInAuth).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			j.BadRequest(c)
			return
		}
		j.ServerError(c, err)
		return
	}

	// 特殊标记，一个不存在的组别，作为超级权限判断
	if !strings.Contains(userInAuth.Admin, "-1") {
		j.NotAcceptable(c)
		return
	}

	allUsers := []miniUser{}

	// 获取所有用户
	models.GetDBHelper().Select("id").Find(&allUsers)
	// tableName := allUsers[0].TableName()

	// 更新所有用户
	err := models.GetDBHelper().Transaction(func(db *gorm.DB) error {
		for idx := range allUsers {
			allUsers[idx].PasswordInStr = helper.GenPass()

			pass, errForPassHash := helper.CalcPassHash(allUsers[idx].PasswordInStr)
			if errForPassHash != nil {
				return errForPassHash
			}

			allUsers[idx].Password = pass

			if errForSQL := db.Model(&allUsers[idx]).Update(&allUsers[idx]).Error; errForSQL != nil {
				return errForSQL
			}
		}

		return nil
	})

	if err != nil {
		j.ServerError(c, err)
		return
	}

	j.Data = map[string]interface{}{
		"allUsers": allUsers,
	}

	j.ResponseOK(c)
	return
}
