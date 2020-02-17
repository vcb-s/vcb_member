package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

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
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	userList := make([]userListResponseRes, 0, req.PageSize)
	var sqlBuilder = models.GetDBHelper().Table("user")

	sqlBuilder.Where("`delete_at` is NULL")

	if req.Group > 0 {
		sqlBuilder.Where("`group` like ?", fmt.Sprintf("%%%d%%", req.Group))
	}
	if req.Retired == 1 {
		sqlBuilder.Where("`retired` = ?", 1)
	}
	if req.Sticky == 1 {
		sqlBuilder.Where("`order` > ?", 0)
	}
	sqlBuilder.Limit(req.PageSize, req.PageSize*(req.Current-1))

	sqlBuilder.OrderBy("`order` desc, `id` asc")

	total, err := sqlBuilder.FindAndCount(&userList)
	if err != nil {
		j.ServerError(c, err)
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
		j.ServerError(c, err)
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
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	hasUser, err := models.GetDBHelper().Where("id = ?", req.UID).Get(&user)
	if err != nil {
		j.ServerError(c, err)
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

	if !helper.CheckPassHash(req.Password, user.Password) {
		j.FailAuth(c)
		return
	}

	// 签发密钥
	token, err := helper.GenToken(user.ID)
	if err != nil {
		j.ServerError(c, err)
		return
	}
	refreshToken, err := helper.GenRefreshToken(user.ID)
	if err != nil {
		j.ServerError(c, err)
		return
	}

	c.Writer.Header().Set("x-token", token)
	c.Writer.Header().Set("x-refreshToken", refreshToken)

	j.ResponseOK(c)
	return
}

// ResetPass 重设自己的密码
func ResetPass(c *gin.Context) {
	var (
		j    JSONData
		req  resetPassReq
		user models.User
	)
	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	uid := c.Request.Header.Get("uid")

	hasValue, err := models.GetDBHelper().Where("id = ?", uid).Get(&user)
	if err != nil {
		j.ServerError(c, err)
		return
	}

	if !hasValue {
		j.BadRequest(c)
		return
	}

	if !helper.CheckPassHash(req.Current, user.Password) {
		j.Message = "密码错误"
		j.FailAuth(c)
		return
	}

	newHash, err := helper.CalcPassHash(req.NewPassword)
	if err != nil {
		j.ServerError(c, err)
		return
	}

	user.Password = newHash
	user.JwtID = ""
	_, err = models.GetDBHelper().Where("id = ?", uid).Cols("password, jwt_id").Update(&user)
	if err != nil {
		j.ServerError(c, err)
		return
	}

	j.ResponseOK(c)
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
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	loginUID := c.Request.Header.Get("uid")

	_, err := models.GetDBHelper().Where("id = ?", loginUID).Get(&adminUser)
	hasUser, err := models.GetDBHelper().Where("id = ?", req.UID).Get(&user)
	if err != nil {
		j.ServerError(c, err)
		return
	}

	if adminUser.SuperAdmin != 1 {
		j.Message = "您不是管理员"
		j.FailAuth(c)
		return
	}

	if !hasUser {
		j.ServerError(c, errors.New("用户不存在"))
		return
	}

	// 产生一个随机密码
	var newPass string = helper.GenRandPass()

	newHash, err := helper.CalcPassHash(newPass)
	if err != nil {
		j.ServerError(c, err)
		return
	}

	user.Password = newHash
	user.JwtID = ""

	_, err = models.GetDBHelper().Where("id = ?", req.UID).Cols("password, jwt_id").Update(&user)
	if err != nil {
		j.ServerError(c, err)
		return
	}

	j.Data = map[string]interface{}{"password": newPass}
	j.ResponseOK(c)
	return
}

// LoginFromWP 绑定登录
func LoginFromWP(c *gin.Context) {
	var (
		j           JSONData
		req         loginWithWPCodeReq
		association models.UserAssociation
	)
	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	// 根据 Authorization code 换取 AccessToken
	accessToken, err := helper.GetAccessTokenFromCode(req.Code)
	if err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	// 根据accessToken换取主站ID
	userInWP, err := helper.GetUserInfoFromAccesstoken(accessToken)
	if err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	// 根据主站ID在第三方绑定表查找
	hasValue, err := models.GetDBHelper().Where("type = ? AND association = ?", models.UserAssociationTypeWP, strconv.Itoa(userInWP.ID)).Get(&association)
	if err != nil {
		j.ServerError(c, err)
		return
	}
	// 没找到就返回没授权
	if !hasValue {
		j.Message = "没有找到用户"
		j.FailAuth(c)
		return
	}
	// 找到了就按照UID签发
	token, err := helper.GenToken(association.UID)
	if err != nil {
		j.ServerError(c, err)
		return
	}
	c.Writer.Header().Set("token", token)

	refreshToken, err := helper.GenRefreshToken(association.UID)
	if err != nil {
		j.ServerError(c, err)
		return
	}
	c.Writer.Header().Set("refreshToken", refreshToken)

	j.ResponseOK(c)
	return
}

// CreateBindForWP 添加主站绑定
func CreateBindForWP(c *gin.Context) {
	var (
		j           JSONData
		req         createBindForWPReq
		association models.UserAssociation
	)
	UID := c.Request.Header.Get("uid")
	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	// 查询用户是否有同类型绑定，不允许重复
	hasValue, err := models.GetDBHelper().Where("type = ? AND uid = ?", models.UserAssociationTypeWP, UID).Get(&association)
	if err != nil {
		j.ServerError(c, err)
		return
	}
	if hasValue {
		j.ServerError(c, errors.New("你已绑定其他账号"))
		return
	}

	// 根据 Authorization code 换取 AccessToken
	accessToken, err := helper.GetAccessTokenFromCode(req.Code)
	if err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	// 根据accessToken换取主站ID
	userInWP, err := helper.GetUserInfoFromAccesstoken(accessToken)
	if err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	// 检查主站ID是否已经跟其他账号绑定过
	hasValue, err = models.GetDBHelper().Where("type = ? AND association = ?", models.UserAssociationTypeWP, strconv.Itoa(userInWP.ID)).Get(&association)
	if err != nil {
		j.ServerError(c, err)
		return
	}
	// 不允许重复绑定
	if hasValue {
		j.Message = "该主站账号已被绑定"
		j.FailAuth(c)
		return
	}

	association.ID = helper.GenID()
	association.UID = UID
	association.Association = strconv.Itoa(userInWP.ID)
	association.Type = models.UserAssociationTypeWP

	// 没绑定过的就添加一条绑定
	_, err = models.GetDBHelper().InsertOne(&association)
	if err != nil {
		j.ServerError(c, err)
		return
	}

	j.ResponseOK(c)
	return
}

// DeleteWPBind 移除主站绑定
func DeleteWPBind(c *gin.Context) {
	var (
		j           JSONData
		association models.UserAssociation
	)
	UID := c.Request.Header.Get("uid")

	association.UID = UID
	association.Type = models.UserAssociationTypeWP

	colsNumber, err := models.GetDBHelper().Delete(&association)
	if err != nil {
		j.ServerError(c, err)
		return
	}
	if colsNumber == 0 {
		j.ServerError(c, errors.New("你未绑定主站账号"))
		return
	}

	j.ResponseOK(c)
	return
}

// UpdateUser 修改用户信息
func UpdateUser(c *gin.Context) {
	var (
		j         JSONData
		req       updateUserReq
		adminUser models.User
	)
	UIDFromToken := c.Request.Header.Get("uid")
	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	if req.ID == "" {
		req.ID = UIDFromToken
	}

	updateBuilder := models.GetDBHelper().Where("id = ?", req.ID).Omit("password,jwt_id")

	// 检查是不是管理员
	_, err := models.GetDBHelper().Where("id = ?", UIDFromToken).Get(&adminUser)
	if err != nil {
		j.ServerError(c, err)
		return
	}

	// 不是管理员的话这些键值不允许修改
	if adminUser.SuperAdmin == 2 {
		if req.ID != UIDFromToken {
			j.Message = "不允许修改他人信息"
			j.FailAuth(c)
			return
		}

		updateBuilder.Omit("super_admin")
	}

	// 修改键值
	_, err = updateBuilder.Update(&req)
	if err != nil {
		j.ServerError(c, err)
		return
	}

	j.ResponseOK(c)
	return
}
