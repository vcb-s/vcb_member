package service

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"vcb_member/models"
)

type kickOffReq struct {
	UID   string `json:"uid" form:"uid" binding:"required"`
	Group string `json:"group" form:"group" binding:"required"`
}

// KickOff 踢出
func KickOff(c *gin.Context) {
	var (
		j             JSONData
		req           kickOffReq
		userToKickOff models.User
		userInAuth    models.User
	)
	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	userInAuth.UID = c.Request.Header.Get("uid")
	userToKickOff.UID = req.UID

	groupsToRemove := map[string]bool{}
	for _, group := range strings.Split(req.Group, ",") {
		groupsToRemove[group] = true
	}

	if err := models.GetDBHelper().Where(userToKickOff).First(&userInAuth).Error; err != nil {
		j.ServerError(c, err)
		return
	}
	if err := models.GetDBHelper().Where(userInAuth).First(&userToKickOff).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	if !userInAuth.CanManagePerson(userToKickOff) {
		j.Message = "你没有权限管理该用户"
		j.BadRequest(c)
		return
	}

	err := models.GetDBHelper().Transaction(func(db *gorm.DB) error {
		// 更新该用户的组别信息
		nextGroups := []string{}
		nextAdmin := []string{}
		for _, group := range strings.Split(userToKickOff.Group, ",") {
			if !groupsToRemove[group] {
				nextGroups = append(nextGroups, group)
			}
		}
		for _, group := range strings.Split(userToKickOff.Admin, ",") {
			if !groupsToRemove[group] {
				nextAdmin = append(nextAdmin, group)
			}
		}
		if err := db.Model(userToKickOff).Updates(map[string]interface{}{
			"Group": strings.Join(nextGroups, ","),
			"Admin": strings.Join(nextAdmin, ","),
		}).Error; err != nil {
			return err
		}

		// 更新该用户的所有卡片
		var userCards []models.UserCard

		if err := db.Where(models.UserCard{UID: userToKickOff.UID}).Find(&userCards).Error; err != nil {
			return err
		}

		// 这里暂时没找到有什么办法可以在一次sql中update多行
		for idx := range userCards {
			nextGroups := []string{}
			for _, group := range strings.Split(userCards[idx].Group, ",") {
				if !groupsToRemove[group] {
					nextGroups = append(nextGroups, group)
				}
			}
			userCards[idx].Group = strings.Join(nextGroups, ",")
			if len(userCards[idx].Group) == 0 {
				userCards[idx].Hide = 1
			}

			if err := db.Model(&userCards[idx]).Where(models.UserCard{ID: userCards[idx].ID}).Updates(map[string]interface{}{
				"Group": userCards[idx].Group,
				"Hide":  userCards[idx].Hide,
			}).Error; err != nil {
				return err
			}
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
