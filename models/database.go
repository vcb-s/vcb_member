package models

import (
	"fmt"
	"sync"
	"vcb_member/conf"

	"github.com/go-xorm/xorm"
)

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
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=true&loc=Local",
		conf.Main.Database.User, conf.Main.Database.Pass, conf.Main.Database.Host, conf.Main.Database.Port, conf.Main.Database.Dbname,
	))
	if err != nil {
		panic(err.Error())
	}
	//test DB if connection
	err = engine.Ping()
	if err != nil {
		panic(err.Error())
	}

	//设置连接池
	engine.SetMaxIdleConns(250) //空闲数大小
	engine.SetMaxOpenConns(300) //最大打开连接数

	//start sql log print
	engine.ShowSQL(true)
	engine.ShowExecTime(true)
	return engine
}
