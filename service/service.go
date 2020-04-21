package service

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"vcb_member/helper"
	"vcb_member/models"
)

// UserCardList 用户列表
func UserCardList(c *gin.Context) {
	var (
		j            JSONData
		req          userListReq
		userCardList = make([]models.UserCard, 0)
	)
	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	var sqlBuilder = models.GetDBHelper().Order("`order` DESC").Order("`id` ASC")

	if req.IncludeHide != 1 {
		sqlBuilder = sqlBuilder.Not("hide", 1)
	}
	if req.Tiny == 1 {
		sqlBuilder = sqlBuilder.Select("`id`, `uid`, `avast`, `nickname`")
	}
	if req.Group > 0 {
		sqlBuilder = sqlBuilder.Where("`group` like ?", fmt.Sprintf("%%%d%%", req.Group))
	}
	if req.Retired == 1 {
		sqlBuilder = sqlBuilder.Where("`retired` = ?", 1)
	}
	if req.Sticky == 1 {
		sqlBuilder = sqlBuilder.Where("`order` > ?", 0)
	}
	if req.PageSize > 0 && req.Current > 0 {
		sqlBuilder = sqlBuilder.Limit(req.PageSize).Offset(req.PageSize * (req.Current - 1))
	}

	// 有指定 CardID 时就忽略 KeyWord 参数
	if len(req.CardID) > 0 {
		sqlBuilder = sqlBuilder.Where("`id` = ?", req.CardID)
	} else if len(req.KeyWord) > 0 {
		keyword := fmt.Sprintf("%%%s%%", req.KeyWord)
		sqlBuilder = sqlBuilder.Where("`bio` like ? OR `nickname` like ? OR `id` = ?", keyword, keyword, req.KeyWord)
	}

	total := 0

	err := sqlBuilder.Find(&userCardList).Count(&total).Error
	if err != nil {
		j.ServerError(c, err)
		return
	}

	// 乱序
	originUserCardListLen := len(userCardList)
	if req.Sticky != 1 && req.Tiny != 1 && originUserCardListLen > 0 {
		// 没有筛选置顶，也就是数组需要乱序
		// 如果筛选了置顶整个数组就是有顺序的
		stickyUserList := make([]models.UserCard, 0)
		lastStickyUserIndex := 0

		// 找到置顶部分
		for i := 0; i < originUserCardListLen; i++ {
			if userCardList[i].Order > 0 {
				lastStickyUserIndex = i
			} else {
				break
			}
		}

		// 置顶用户列表
		stickyUserList = userCardList[:lastStickyUserIndex+1]
		userCardList = userCardList[lastStickyUserIndex+1:]

		rand.Shuffle(len(userCardList), func(i, j int) {
			userCardList[i], userCardList[j] = userCardList[j], userCardList[i]
		})

		userCardList = append(stickyUserList, userCardList...)
	}

	j.Data = map[string]interface{}{"res": userCardList, "total": total}
	j.ResponseOK(c)
	return
}

// TinyUserCardList 用户列表（登录页面用）
func TinyUserCardList(c *gin.Context) {
	var (
		j            JSONData
		userCardList = make([]models.UserCard, 0)
	)

	var sqlBuilder = models.GetDBHelper().Select("`id`, `uid`, `avast`, `nickname`")

	total := 0

	err := sqlBuilder.Find(&userCardList).Count(&total).Error
	if err != nil {
		j.ServerError(c, err)
		return
	}

	j.Data = map[string]interface{}{"res": userCardList, "total": total}
	j.ResponseOK(c)
	return
}

// GroupList 组别列表
func GroupList(c *gin.Context) {
	var (
		j JSONData
	)

	userGroupList := make([]models.UserCardGroup, 0)

	total := 0

	if err := models.GetDBHelper().Find(&userGroupList).Count(&total).Error; err != nil {
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

	if err := models.GetDBHelper().First(&user, "`id` = ?", req.UID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			j.Message = "用户不存在"
			j.NotAcceptable(c)
			return
		}
		j.ServerError(c, err)
		return
	}

	// 如果用户密码为空，返回提示找网站组设置初始密码
	if user.Password == "" {
		j.Message = "请先联系网站组设置登录数据"
		j.NotAcceptable(c)
		return
	}

	if !helper.CheckPassHash(req.Password, user.Password) {
		j.Message = "密码不正确"
		j.FailAuth(c)
		return
	}

	// 签发密钥
	token, err := helper.GenToken(user.UID)
	if err != nil {
		j.ServerError(c, err)
		return
	}

	c.Writer.Header().Set("X-Token", token)

	j.ResponseOK(c)
	return
}

// ResetPass 重设密码
func ResetPass(c *gin.Context) {
	var (
		j           JSONData
		req         resetPassReq
		userInAuth  models.User
		userToReset models.User
	)
	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	uidInAuth := c.Request.Header.Get("uid")
	if len(req.UID) == 0 {
		userToReset.UID = uidInAuth
	} else {
		userToReset.UID = req.UID
	}

	if err := models.GetDBHelper().First(&userInAuth, "`id` = ?", uidInAuth).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			j.BadRequest(c)
			return
		}
		j.ServerError(c, err)
		return
	}

	if !userInAuth.IsAdmin() {
		if !helper.CheckPassHash(req.Current, userInAuth.Password) {
			j.Message = "密码错误"
			j.FailAuth(c)
			return
		}
	}

	newPass, err := helper.CalcPassHash(req.NewPassword)
	if err != nil {
		j.ServerError(c, err)
		return
	}

	userToReset.Password = newPass

	resylt := models.GetDBHelper().Model(&userToReset).Update("password")
	if resylt.Error != nil {
		j.ServerError(c, err)
		return
	}
	if resylt.RowsAffected == 0 {
		j.ServerError(c, errors.New("用户不存在"))
		return
	}

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
	err = models.GetDBHelper().Where("type = ? AND association = ?", models.UserAssociationTypeWP, strconv.Itoa(userInWP.ID)).First(&association).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			j.Message = "没有找到用户"
			j.FailAuth(c)
			return
		}
		j.ServerError(c, err)
		return
	}

	// 找到了就按照UID签发
	token, err := helper.GenToken(association.UID)
	if err != nil {
		j.ServerError(c, err)
		return
	}
	c.Writer.Header().Set("token", token)

	j.ResponseOK(c)
	return
}

// CreateBindForWP 添加主站绑定
func CreateBindForWP(c *gin.Context) {
	var (
		j          JSONData
		req        createBindForWPReq
		userToBind models.UserAssociation
	)
	uidToBind := c.Request.Header.Get("uid")
	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	// 查询用户是否有同类型绑定，不允许重复
	err := models.GetDBHelper().Where("type = ? AND uid = ?", models.UserAssociationTypeWP, uidToBind).First(&userToBind).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			j.ServerError(c, errors.New("你已绑定其他账号"))
			return
		}
		j.ServerError(c, err)
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
	err = models.GetDBHelper().Where("type = ? AND association = ?", models.UserAssociationTypeWP, strconv.Itoa(userInWP.ID)).First(&userToBind).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	// 不允许重复绑定
	if err == nil {
		j.Message = "该主站账号已被绑定"
		j.FailAuth(c)
		return
	}

	userToBind.ID = helper.GenID()
	userToBind.UID = uidToBind
	userToBind.AuthCode = strconv.Itoa(userInWP.ID)
	userToBind.Type = models.UserAssociationTypeWP

	// 没绑定过的就添加一条绑定
	err = models.GetDBHelper().Create(&userToBind).Error
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
	return
}

// UpdateUser 修改用户信息
func UpdateUser(c *gin.Context) {
	var (
		j            JSONData
		req          updateUserReq
		userToUpdate models.UserCard
		userInAuth   models.User = models.User{}
	)

	userInAuth.UID = c.Request.Header.Get("uid")

	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	// 查询所属UID
	if err := models.GetDBHelper().Where("`id` = ?", req.ID).First(&userToUpdate).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	// 覆盖前端传入的UID
	req.UID = userToUpdate.UID

	// 查询权限
	if err := models.GetDBHelper().First(&userInAuth, "`id` = ?", userInAuth.UID).Error; err != nil {
		j.ServerError(c, err)
		return
	}

	// 不是管理员且uid不匹配的话
	if userToUpdate.UID != userInAuth.UID && !userInAuth.IsAdmin() {
		j.Message = "不允许修改他人信息"
		j.FailAuth(c)
		return
	}

	updateBuilder := models.GetDBHelper().Where("id = ?", req.ID)

	// 修改键值
	result := updateBuilder.Update(&req)
	if result.Error != nil {
		j.ServerError(c, result.Error)
		return
	}

	if result.RowsAffected == 0 {
		j.Message = "没有修改任何内容"
		j.ServerError(c, result.Error)
		return
	}

	j.ResponseOK(c)
	return
}

// PersonInfo 个人卡片列表
func PersonInfo(c *gin.Context) {
	var (
		j   JSONData
		req personInfoReq

		uidInAuth     = c.Request.Header.Get("uid")
		userInAuth    models.User
		userInRequest models.User

		userCardList = make([]models.UserCard, 0)
		userList     = make([]models.User, 0)
	)

	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	uidInRequest := req.UID
	if uidInRequest == "" {
		uidInRequest = uidInAuth
	}

	if err := models.GetDBHelper().First(&userInAuth, "`id` = ?", uidInAuth).Error; err != nil || !userInAuth.CanManagePerson(uidInRequest) {
		j.Message = "你无权获取该用户信息"
		j.BadRequest(c)
		return
	}

	if err := models.GetDBHelper().First(&userInRequest, "`id` = ?", uidInRequest).Error; err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	cardTotal := 0
	userTotal := 0
	{
		var sqlBuilder = models.GetDBHelper().Where("`uid` = ?", uidInRequest)

		err := sqlBuilder.Find(&userCardList).Count(&cardTotal).Error
		if err != nil {
			j.Message = err.Error()
			j.BadRequest(c)
			return
		}
	}

	if userInRequest.IsAdmin() {
		var sqlBuilder = models.GetDBHelper()
		groupsBelongUser := strings.Split(userInRequest.Admin, ",")

		for _, group := range groupsBelongUser {
			sqlBuilder = sqlBuilder.Or("`group` like ?", fmt.Sprintf("%%%s%%", group))
		}

		if err := sqlBuilder.Find(&userList).Count(&userTotal).Error; err != nil {
			j.Message = err.Error()
			j.BadRequest(c)
			return
		}
	}

	j.Data = map[string]interface{}{
		"info": userInRequest,
		"cards": map[string]interface{}{
			"total": cardTotal,
			"res":   userCardList,
		},
		"users": map[string]interface{}{
			"total": userTotal,
			"res":   userList,
		},
	}
	j.ResponseOK(c)
	return
}
