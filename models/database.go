package models

// User 用户表
type User struct {
	ID       string `xorm:"id"`
	Retired  int    `xorm:"retired"`
	Avast    string `xorm:"avast"`
	Bio      string `xorm:"bio"`
	Nickname string `xorm:"nickname"`
	Job      string `xorm:"job"`
	Order    int    `xorm:"order"`
	Password string `xorm:"password"`
	Group    string `xorm:"group"`
	TokenID  string `xorm:"jwt_id"`
}

// UserGroup 组别表
type UserGroup struct {
	ID   int    `xorm:"id"`
	Name string `xorm:"name"`
}

// WpAssociation 主站关联授权表
type WpAssociation struct {
	//
}
