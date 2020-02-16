package models

import (
	"fmt"
	"log"
	"sync"

	"github.com/go-xorm/xorm"

	"vcb_member/conf"
)

// User 用户表
type User struct {
	ID         string `json:"id" form:"id" xorm:"id"`
	Retired    int    `json:"retired" form:"retired" xorm:"retired"`
	Avast      string `json:"avast" form:"avast" xorm:"avast"`
	Bio        string `json:"bio" form:"bio" xorm:"bio"`
	Nickname   string `json:"nickname" form:"nickname" xorm:"nickname"`
	Job        string `json:"job" form:"job" xorm:"job"`
	Order      int    `json:"order" form:"order" xorm:"order"`
	Password   string `json:"password" form:"password" xorm:"password"`
	Group      string `json:"group" form:"group" xorm:"group"`
	JwtID      string `json:"jwt_id" form:"jwt_id" xorm:"jwt_id"`
	SuperAdmin int    `json:"super_admin" form:"super_admin" xorm:"super_admin"`
}

// TableName 指示 User 表名
func (m User) TableName() string {
	return "user"
}

// UserGroup 组别表
type UserGroup struct {
	ID   int    `json:"id" form:"id" xorm:"id"`
	Name string `json:"name" form:"name" xorm:"name"`
}

// TableName 指示 UserGroup 表名
func (m UserGroup) TableName() string {
	return "user_group"
}

// UserAssociationType 账号绑定类型枚举
type UserAssociationType int8

const (
	// UserAssociationTypeWP 绑定类型 - 主站
	UserAssociationTypeWP UserAssociationType = 1
)

// UserAssociation 绑定授权表
type UserAssociation struct {
	ID          string              `json:"id" form:"id" xorm:"id"`
	UID         string              `json:"uid" form:"uid" xorm:"uid"`
	Association string              `json:"association" form:"association" xorm:"association"`
	Type        UserAssociationType `json:"type" form:"type" xorm:"type"`
}

// TableName 指示 UserAssociation 表名
func (m UserAssociation) TableName() string {
	return "user_association"
}

var instance *xorm.Engine
var once sync.Once

func init() {
	GetDBHelper()
}

// GetDBHelper 获取数据库实例
func GetDBHelper() *xorm.Engine {
	once.Do(func() {
		instance = newDBHelper()
	})
	return instance
}

func newDBHelper() *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", fmt.Sprintf(
		"%v:%v@tcp([%v]:%v)/%v?charset=utf8mb4&parseTime=true&loc=Local",
		conf.Main.Database.User,
		conf.Main.Database.Pass,
		conf.Main.Database.Host,
		conf.Main.Database.Port,
		conf.Main.Database.Dbname,
	))
	if err != nil {
		log.Fatalln("xorm err", err)
	}
	//test DB if connection
	err = engine.Ping()
	if err != nil {
		log.Fatalln("xorm Ping err", err)
	}

	//设置连接池
	engine.SetMaxIdleConns(2)     //空闲数大小
	engine.SetMaxOpenConns(10)    //最大打开连接数
	engine.SetConnMaxLifetime(-1) //重用超时

	//start sql log print
	engine.ShowSQL(true)
	engine.ShowExecTime(true)
	return engine
}
