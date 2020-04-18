package models

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jinzhu/gorm"

	"vcb_member/conf"
)

// User 用户表
type User struct {
	UID      string `json:"id,omitempty" form:"id" gorm:"PRIMARY_KEY;column:id"`
	Password string `json:"password,omitempty" form:"password" gorm:"column:pass"`
	IsAdmin  int8   `json:"isAdmin,omitempty" form:"isAdmin" gorm:"column:is_admin"`
	Ban      int8   `json:"ban,omitempty" form:"ban" gorm:"column:ban"`
	SoftDeletedModel
}

// TableName 指示 User 表名
func (m User) TableName() string {
	return "user"
}

// IsBan 是否可以管理对应uid用户
func (m User) IsBan() bool {
	return m.Ban == 1
}

// CanManagePerson 是否可以管理对应uid用户
func (m User) CanManagePerson(uidInRequest string) bool {
	return !m.IsBan() && (m.IsAdmin == 1 || m.UID == uidInRequest)
}

// UserCard 用户表
type UserCard struct {
	ID       string `json:"id,omitempty" form:"id" gorm:"PRIMARY_KEY;column:id"`
	UID      string `json:"uid,omitempty" form:"uid" gorm:"column:uid"`
	Group    string `json:"group,omitempty" form:"group" gorm:"column:group"`
	Order    int    `json:"order,omitempty" form:"order" gorm:"column:order"`
	Avast    string `json:"avast,omitempty" form:"avast" gorm:"column:avast"`
	Bio      string `json:"bio,omitempty" form:"bio" gorm:"column:bio"`
	Nickname string `json:"nickname,omitempty" form:"nickname" gorm:"column:nickname"`
	Job      string `json:"job,omitempty" form:"job" gorm:"column:job"`
	Hide     int    `json:"hide,omitempty" form:"hide" gorm:"column:hide"`
	Retired  int    `json:"retired,omitempty" form:"retired" gorm:"column:retired"`
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
	ID   int    `json:"id,omitempty" form:"id" gorm:"PRIMARY_KEY;column:id"`
	Name string `json:"name,omitempty" form:"name" gorm:"column:name"`
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
	ID       string              `json:"id,omitempty" form:"id" gorm:"PRIMARY_KEY;column:id"`
	UID      string              `json:"uid,omitempty" form:"uid" gorm:"column:uid"`
	AuthCode string              `json:"association,omitempty" form:"association" gorm:"column:auth"`
	Type     UserAssociationType `json:"type,omitempty" form:"type" gorm:"column:type"`
	SoftDeletedModel
}

// TableName 指示 UserAssociation 表名
func (m UserAssociation) TableName() string {
	return "login_association"
}

var instance *gorm.DB
var once sync.Once

func init() {
	GetDBHelper()
}

// GetDBHelper 获取数据库实例
func GetDBHelper() *gorm.DB {
	once.Do(func() {
		instance = newDBHelper()
	})
	return instance
}

func newDBHelper() *gorm.DB {
	engine, err := gorm.Open("mysql", fmt.Sprintf(
		"%v:%v@tcp([%v]:%v)/%v?charset=utf8mb4&parseTime=true&loc=Local",
		conf.Main.Database.User,
		conf.Main.Database.Pass,
		conf.Main.Database.Host,
		conf.Main.Database.Port,
		conf.Main.Database.Dbname,
	))
	if err != nil {
		log.Fatalln("gorm err", err)
	}
	//test DB if connection
	err = engine.DB().Ping()
	if err != nil {
		log.Fatalln("gorm Ping err", err)
	}

	//设置连接池
	engine.DB().SetMaxIdleConns(10)           //空闲数大小
	engine.DB().SetMaxOpenConns(100)          //最大打开连接数
	engine.DB().SetConnMaxLifetime(time.Hour) //重用超时

	engine.LogMode(true)
	return engine
}
