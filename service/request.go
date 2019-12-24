package service

// 是否 --- 12

type userListReq struct {
	Group    int `json:"group" form:"group"`
	Retired  int `json:"retired" form:"retired"`
	Sticky   int `json:"sticky" form:"sticky"`
	Current  int `json:"current" form:"current" binding:"required,min=1"`
	PageSize int `json:"pageSize" form:"pageSize" binding:"required,min=1,max=20"`
}
