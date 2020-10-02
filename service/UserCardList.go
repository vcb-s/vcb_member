package service

import (
	"fmt"
	"math/rand"

	"github.com/gin-gonic/gin"

	"vcb_member/models"
)

type userCardListReq struct {
	Current  int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	CardID   string `json:"id" form:"id"`
	KeyWord  string `json:"keyword" form:"keyword"`
	Group    int    `json:"group" form:"group"`
	Retired  int    `json:"retired" form:"retired"`
	// 只返回置顶相关卡片
	Sticky int `json:"sticky" form:"sticky"`
	// 只返回 ID/UID/头像/昵称
	Tiny int `json:"tiny" form:"tiny"`
	// 把隐藏设置为1的也范围
	IncludeHide int `json:"includeHide" form:"includeHide"`
	// 不乱序
	InOrder int `json:"inOrder" form:"inOrder"`
}

// UserCardList 用户列表
func UserCardList(c *gin.Context) {
	var (
		j            JSONData
		req          userCardListReq
		userCardList = make([]models.UserCard, 0)
	)
	if err := c.ShouldBind(&req); err != nil {
		j.Message = err.Error()
		j.BadRequest(c)
		return
	}

	UserCardTableName := models.UserCard{}.TableName()
	UserTableName := models.User{}.TableName()

	var sqlBuilder = models.GetDBHelper().Table(UserCardTableName)
	sqlBuilder = sqlBuilder.Order("`order` DESC").Order("`id` ASC")

	sqlBuilder = sqlBuilder.Joins(fmt.Sprintf(
		"left join %s on %s.id = %s.uid",
		UserTableName,
		UserTableName,
		UserCardTableName,
	))

	sqlBuilder = sqlBuilder.Where(
		fmt.Sprintf("`%s`.`ban` <> ?", UserTableName),
		1,
	)

	// 如果没有指定包括隐藏
	if req.IncludeHide != 1 {
		// 那就指定hide不包含1（是）
		sqlBuilder = sqlBuilder.Not(`hide`, 1)
	}
	if req.Tiny == 1 {
		sqlBuilder = sqlBuilder.Select(
			fmt.Sprintf(
				"`%s`.`id`, `%s`.`uid`, `%s`.`avast`, `%s`.`nickname`",
				UserCardTableName,
				UserCardTableName,
				UserCardTableName,
				UserCardTableName,
			),
		)
	}
	if req.Group > 0 {
		sqlBuilder = sqlBuilder.Where(
			fmt.Sprintf("`%s`.`group` like ?", UserCardTableName),
			fmt.Sprintf("%%%d%%", req.Group),
		)
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
	if req.CardID != "" {
		sqlBuilder = sqlBuilder.Where(fmt.Sprintf(
			"`%s`.`id` = ?",
			UserCardTableName,
		), req.CardID)
	} else if req.KeyWord != "" {
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
	if req.InOrder != 1 && originUserCardListLen > 1 {
		// 没有指定按序就是乱序
		stickyUserList := make([]models.UserCard, 0)
		lastStickyUserIndex := -1

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
