package service

// ?? --- 12

type userListReq struct {
	Group    int `json:"group" form:"group"`
	Retired  int `json:"retired" form:"retired"`
	Sticky   int `json:"sticky" form:"sticky"`
	Current  int `json:"current" form:"current" binding:"required,min=1"`
	PageSize int `json:"pageSize" form:"pageSize" binding:"required,min=1,max=20"`
}
type loginReq struct {
	UID      string `json:"uid" form:"uid" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}
type resetPassReq struct {
	Current     string `json:"current" form:"current" binding:"required"`
	NewPassword string `json:"new" form:"new" binding:"required"`
}
