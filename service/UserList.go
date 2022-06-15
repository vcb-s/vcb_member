package service

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"vcb_member/models"
)

type userListReq struct {
	Current    int    `json:"page" form:"page"`
	PageSize   int    `json:"pageSize" form:"pageSize"`
	ID         string `json:"id" form:"id"`
	KeyWord    string `json:"keyword" form:"keyword"`
	Group      int    `json:"group" form:"group"`
	Retired    int    `json:"retired" form:"retired"`
	IncludeBan int    `json:"includeBan" form:"includeBan"`
}

// UserList 用户列表
func UserList(c *gin.Context) {
	var (
		j            JSONData
		req          userListReq
		userCardList = make([]models.User, 0)
	)
	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	UserTableName := models.User{}.TableName()

	var sqlBuilder = models.GetDBHelper().Table(UserTableName)
	sqlBuilder = sqlBuilder.Order("`id` ASC")

	if req.IncludeBan != 1 {
		sqlBuilder = sqlBuilder.Where(
			fmt.Sprintf("`%s`.`ban` <> ?", UserTableName),
			1,
		)
	}

	if req.Group > 0 {
		sqlBuilder = sqlBuilder.Where(
			fmt.Sprintf("`%s`.`group` like ?", UserTableName),
			fmt.Sprintf("%%%d%%", req.Group),
		)
	}
	if req.Retired == 1 {
		sqlBuilder = sqlBuilder.Where("`retired` = ?", 1)
	}
	if req.PageSize > 0 && req.Current > 0 {
		sqlBuilder = sqlBuilder.Limit(req.PageSize).Offset(req.PageSize * (req.Current - 1))
	}

	// 有指定 ID 时就忽略 KeyWord 参数
	if req.ID != "" {
		sqlBuilder = sqlBuilder.Where(fmt.Sprintf(
			"`%s`.`id` = ?",
			UserTableName,
		), req.ID)
	} else if req.KeyWord != "" {
		keyword := fmt.Sprintf("%%%s%%", req.KeyWord)
		sqlBuilder = sqlBuilder.Where("`bio` like ? OR `nickname` like ? OR `id` = ?", keyword, keyword, req.KeyWord)
	}

	total := int64(0)

	err := sqlBuilder.Find(&userCardList).Count(&total).Error
	if err != nil {
		j.ServerError(c, err)
		return
	}

	j.Data = map[string]interface{}{"res": userCardList, "total": total}
	j.ResponseOK(c)
}
