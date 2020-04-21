package service

import (
	"vcb_member/models"
)

// yes/no --- 12

type userSearchReq struct {
	Keyword string `json:"keyword" form:"keyword" binding:"required"`
}
type userListReq struct {
	CardID      string `json:"id" form:"id"`
	KeyWord     string `json:"keyword" form:"keyword"`
	IncludeHide int    `json:"includeHide" form:"includeHide"`
	Group       int    `json:"group" form:"group"`
	Retired     int    `json:"retired" form:"retired"`
	Sticky      int    `json:"sticky" form:"sticky"`
	Current     int    `json:"page" form:"page"`
	PageSize    int    `json:"pageSize" form:"pageSize"`
	Tiny        int    `json:"tiny" form:"tiny"`
}
type loginReq struct {
	UID      string `json:"uid" form:"uid" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}
type resetPassReq struct {
	UID         string `json:"uid" form:"uid"`
	Current     string `json:"current" form:"current"`
	NewPassword string `json:"new" form:"new" binding:"required"`
}
type loginWithWPCodeReq struct {
	Code string `json:"code" form:"code" binding:"required"`
}
type createBindForWPReq = loginWithWPCodeReq
type updateUserReq struct {
	models.UserCard
	ID  string `json:"id" form:"id" gorm:"PRIMARY_KEY;column:id" binding:"required"`
	UID string `json:"-" form:"-" gorm:"column:uid"`
}
type personInfoReq struct {
	UID string `json:"uid,omitempty" form:"uid,omitempty"`
}
