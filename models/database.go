package models

import (
	"strings"
)

// User 用户表
type User struct {
	ID         string `json:"id" form:"id" gorm:"primaryKey;column:id"`
	Password   string `json:"-" form:"-" gorm:"column:pass"`
	Admin      string `json:"admin" form:"admin" gorm:"column:admin"`
	Ban        int8   `json:"ban" form:"ban" gorm:"column:ban"`
	SuperAdmin int8   `json:"-" form:"-" gorm:"column:super_admin"`
	Avast      string `json:"avast" form:"avast" gorm:"column:avast"`
	Nickname   string `json:"nickname" form:"nickname" gorm:"column:nickname"`
	Group      string `json:"group" form:"group" gorm:"column:group"`
	SoftDeletedModel
}

// TableName 指示 User 表名
func (m User) TableName() string {
	return "user"
}

// IsAdmin 是否是管理层
func (m User) IsAdmin() bool {
	return len(m.Admin) > 0
}

// IsAdmin 是否是超管
func (m User) IsSuperAdmin() bool {
	return m.SuperAdmin == 1
}

// IsBan 是否是被封禁用户
func (m User) IsBan() bool {
	return m.Ban == 1
}

// CanAccessPerson 是否可以查看对应用户
func (m User) CanAccessPerson(user User) bool {
	if m.IsBan() {
		return false
	}

	if m.ID == user.ID {
		return true
	}

	// 如果自己是还没被禁用的超管,放行
	if m.IsSuperAdmin() {
		return true
	}

	canManage := false
	for _, groupID := range strings.Split(user.Group, ",") {
		if strings.Contains(m.Admin, groupID) {
			canManage = true
		}
	}

	return canManage
}

// CanManagePerson 是否可以管理对应用户
func (m User) CanManagePerson(user User) bool {
	if m.IsBan() {
		return false
	}

	if m.ID == user.ID {
		return true
	}

	// 如果自己是还没被禁用的超管,放行
	if m.IsSuperAdmin() {
		return true
	}

	// 除非是超管,否则不能操作管理员
	// 因为 目前 前端 采用纯uid方式来进行操作
	if user.IsAdmin() {
		return false
	}

	canManage := false
	for _, groupID := range strings.Split(user.Group, ",") {
		if strings.Contains(m.Admin, groupID) {
			canManage = true
		}
	}

	return canManage
}

// IsContainAllGroup 判断传入用户的组别是否是自己的子集
func (m User) IsContainAllGroup(user User) bool {
	if m.IsAdmin() {
		return false
	}
	if m.IsBan() {
		return false
	}
	if m.ID == user.ID {
		return true
	}

	allContainer := true

	for _, groupID := range strings.Split(user.Group, ",") {
		if !strings.Contains(m.Admin, groupID) {
			allContainer = false
		}
	}

	return allContainer
}

// UserCard 卡片表
type UserCard struct {
	ID       string `json:"id" form:"id" gorm:"primaryKey;column:id"`
	UID      string `json:"uid" form:"uid" gorm:"column:uid"`
	Group    string `json:"group" form:"group" gorm:"column:group"`
	Order    int    `json:"order" form:"order" gorm:"column:order"`
	Avast    string `json:"avast" form:"avast" gorm:"column:avast"`
	Bio      string `json:"bio" form:"bio" gorm:"column:bio"`
	Nickname string `json:"nickname" form:"nickname" gorm:"column:nickname"`
	Job      string `json:"job" form:"job" gorm:"column:job"`
	Hide     int    `json:"hide" form:"hide" gorm:"column:hide"`
	Retired  int    `json:"retired" form:"retired" gorm:"column:retired"`
	SoftDeletedModel
}

// TableName 指示 User 表名
func (m UserCard) TableName() string {
	return "user_crad"
}

// IsRetired 指示 卡片 是否已退休
func (m UserCard) IsRetired() bool {
	return m.Retired == 1
}

// IsHide 指示 卡片 是否被隐藏
func (m UserCard) IsHide() bool {
	return m.Hide == 1
}

// UserCardGroup 组别表
type UserCardGroup struct {
	ID   int    `json:"id" form:"id" gorm:"primaryKey;column:id"`
	Name string `json:"name" form:"name" gorm:"column:name"`
}

// TableName 指示 UserGroup 表名
func (m UserCardGroup) TableName() string {
	return "user_card_group"
}

// UserAssociationType 账号绑定类型枚举
type UserAssociationType int8

const (
	// UserAssociationTypeWP 绑定类型 - 主站
	UserAssociationTypeWP UserAssociationType = 1
)

// UserAssociation 绑定授权表
type UserAssociation struct {
	ID       string              `json:"id" form:"id" gorm:"primaryKey;column:id"`
	UID      string              `json:"uid" form:"uid" gorm:"column:uid"`
	AuthCode string              `json:"association" form:"association" gorm:"column:auth"`
	Type     UserAssociationType `json:"type" form:"type" gorm:"column:type"`
	SoftDeletedModel
}

// TableName 指示 UserAssociation 表名
func (m UserAssociation) TableName() string {
	return "login_association"
}
