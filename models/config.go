package models

import (
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/go-xorm/xorm"
	"github.com/pelletier/go-toml"
)

type conf struct {
	Database struct {
		Host   string
		Port   int
		User   string
		Pass   string
		Dbname string
	}
	Jwt struct {
		Mac        string
		Encryption string
	}
}

// Conf 配置内容
var Conf conf
var instance *xorm.Engine
var once sync.Once

func init() {
	tomlFile, err := ioutil.ReadFile("./config.toml")

	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(tomlFile, &Conf)

	if err != nil {
		panic(err)
	}

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
		Conf.Database.User, Conf.Database.Pass, Conf.Database.Host, Conf.Database.Port, Conf.Database.Dbname,
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
