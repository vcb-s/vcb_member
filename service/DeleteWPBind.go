package service

import (
	"errors"

	"github.com/gin-gonic/gin"

	"vcb_member/models"
)

// DeleteWPBind 移除主站绑定
func DeleteWPBind(c *gin.Context) {
	var (
		j           JSONData
		association models.UserAssociation
	)
	UID := c.Request.Header.Get("uid")

	association.UID = UID
	association.Type = models.UserAssociationTypeWP

	result := models.GetDBHelper().Delete(&association)
	if result.Error != nil {
		j.ServerError(c, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		j.ServerError(c, errors.New("你未绑定主站账号"))
		return
	}

	j.ResponseOK(c)
}
