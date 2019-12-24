package service

import (
	"fmt"
	"vcb_member/models"

	"github.com/gin-gonic/gin"
)

// UserList 用户列表
func UserList(c *gin.Context) {
	var (
		j   JSONData
		req userListReq
	)
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println(err)
		j.BadRequest(c)
		return
	}

	userList := make([]userListResponseRes, 0, req.PageSize)
	var sqlBuilder = models.GetDBHelper().Table("user")
	if req.Group > 0 {
		sqlBuilder.Where(`group = ?`, req.Group)
	}
	if req.Retired == 1 {
		sqlBuilder.Where(`retired = ?`, 1)
	}
	if req.Sticky == 1 {
		sqlBuilder.Where(`"order" > ?`, 0)
	}
	sqlBuilder.Limit(req.PageSize, req.PageSize*(req.Current-1))

	sqlBuilder.OrderBy(`"order" desc, id asc`)

	total, err := sqlBuilder.FindAndCount(&userList)
	if err != nil {
		fmt.Println(err)
		j.ServerError(c)
		return
	}

	j.Data = map[string]interface{}{"res": userList, "total": total}

	j.ResponseOK(c)
	return
}

// GroupList 组别列表
func GroupList(c *gin.Context) {
	return
}

// Login 组别列表
func Login(c *gin.Context) {
	return
}
