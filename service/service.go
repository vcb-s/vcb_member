package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"

	"vcb_member/helper"
	"vcb_member/models"
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
	var (
		j JSONData
	)

	userGroupList := make([]userGroupListResponseRes, 0)
	var sqlBuilder = models.GetDBHelper().Table("user_group")

	total, err := sqlBuilder.FindAndCount(&userGroupList)
	if err != nil {
		fmt.Println(err)
		j.ServerError(c)
		return
	}

	j.Data = map[string]interface{}{"res": userGroupList, "total": total}
	j.ResponseOK(c)
	return
}

// Login 登录
func Login(c *gin.Context) {
	var (
		j    JSONData
		req  loginReq
		user models.User
	)
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println(err)
		j.BadRequest(c)
		return
	}

	hasUser, err := models.GetDBHelper().Table("user").Where("id = ?", req.UID).Get(&user)
	if err != nil {
		fmt.Println(err)
		j.BadRequest(c)
		return
	}

	if !hasUser {
		j.Message = "用户不存在"
		j.NotAcceptable(c)
		return
	}

	// 如果用户密码为空，返回提示找网站组设置初始密码
	if user.Password == "" {
		j.Message = "请先联系网站组设置登录数据"
		j.NotAcceptable(c)
		return
	}

	if !helper.CheckPass(req.Password, user.Password) {
		j.Message = user.Password
		j.NotAcceptable(c)
		return
	}

	// 签发密钥
	token, err := helper.GenToken(user.ID)
	if err != nil {
		j.ServerError(c)
		return
	}
	refreshToken, err := helper.GenRefreshToken(user.ID)
	if err != nil {
		j.ServerError(c)
		return
	}

	c.Writer.Header().Set("x-token", token)
	c.Writer.Header().Set("x-refreshToken", refreshToken)

	j.ResponseOK(c)

	// 密码字段不为空则进行验证
	return
}

// ResetPassForSuperAdmin 重设密码
func ResetPassForSuperAdmin(c *gin.Context) {
	var (
		j         JSONData
		req       loginReq
		user      models.User
		adminUser models.User
	)
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println(err)
		j.BadRequest(c)
		return
	}

	loginUID := c.Request.Header.Get(sercertUIDHeaderKey)

	_, err := models.GetDBHelper().Table("user").Where("id = ?", loginUID).Get(&adminUser)
	hasUser, err := models.GetDBHelper().Table("user").Where("id = ?", req.UID).Get(&user)
	if err != nil || adminUser.SuperAdmin != 1 {
		j.FailAuth(c)
		return
	}

	if !hasUser {
		j.Message = "用户不存在"
		j.NotAcceptable(c)
		return
	}

	// 产生一个明文密钥
	var newPass string
	for i := 0; i < 8; i++ {
		newPass += strconv.Itoa(rand.Intn(9))
	}

	newPassword, err := helper.GenPass(newPass)
	if err != nil {
		j.Message = err.Error()
		j.ServerError(c)
		return
	}

	user.Password = newPassword

	_, err = models.GetDBHelper().Table("user").Where("id = ?", req.UID).Cols("password").Update(&user)
	if err != nil {
		j.Message = err.Error()
		j.ServerError(c)
		return
	}

	j.Data = map[string]interface{}{"password": newPass}
	j.ResponseOK(c)
	return
}
