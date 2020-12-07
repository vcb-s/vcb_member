package service

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"vcb_member/models"
)

type retireReq struct {
	UID   string `json:"uid" form:"uid" binding:"required"`
	Group string `json:"group" form:"group" binding:"required"`
}

type userCardIOnlyID struct {
	ID string `json:"id" form:"id" gorm:"primaryKey;column:id"`
}

// TableName 指示 User 表名
func (m userCardIOnlyID) TableName() string {
	return models.UserCard{}.TableName()
}

// Retire 退休
func Retire(c *gin.Context) {
	var (
		j            JSONData
		req          kickOffReq
		userToRetire models.User
		userInAuth   models.User
	)
	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	userInAuth.ID = c.Request.Header.Get("uid")
	userToRetire.ID = req.UID

	groupsToRemove := map[string]bool{}
	for _, group := range strings.Split(req.Group, ",") {
		groupsToRemove[group] = true
	}

	if err := models.GetDBHelper().Where(userInAuth).First(&userInAuth).Error; err != nil {
		j.ServerError(c, err)
		return
	}
	if err := models.GetDBHelper().Where(userToRetire).First(&userToRetire).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	if !userInAuth.CanManagePerson(userToRetire) {
		j.Message = "你没有权限管理该用户"
		j.BadRequest(c)
		return
	}

	err := models.GetDBHelper().Transaction(func(db *gorm.DB) error {
		// 更新该用户的所有卡片
		var userCards []userCardIOnlyID

		if err := db.Where(models.UserCard{UID: userToRetire.ID}).Find(&userCards).Error; err != nil {
			return err
		}

		var cardIDs []string

		for _, card := range userCards {
			cardIDs = append(cardIDs, card.ID)
		}

		if err := db.Model(&models.UserCard{}).Where(`id IN ?`, cardIDs).Updates(models.UserCard{Retired: 1}).Error; err != nil {
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
