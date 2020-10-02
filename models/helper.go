package models

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"

	"vcb_member/conf"
)

var dbInstance *gorm.DB
var dbOnce sync.Once
var redisOnce sync.Once

var authcodeRedisInstance *redis.Client
var redisContext = context.Background()

func init() {
	GetDBHelper()
	GetAuthCodeRedisHelper()
}

// GetDBHelper 获取数据库实例
func GetDBHelper() *gorm.DB {
	dbOnce.Do(func() {
		dbInstance = newDBHelper()
	})
	return dbInstance
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
		log.Panic().Err(err).Msg("gorm auth error")
	}

	//test DB if connection
	err = engine.DB().Ping()
	if err != nil {
		log.Panic().Err(err).Msg("gorm ping error")
	}

	//设置连接池
	engine.DB().SetMaxIdleConns(10)           //空闲数大小
	engine.DB().SetMaxOpenConns(100)          //最大打开连接数
	engine.DB().SetConnMaxLifetime(time.Hour) //重用超时

	if conf.Main.Debug {
		engine.SetLogger(&log.Logger)
		engine.LogMode(true)
	} else {
		engine.LogMode(false)
	}
	return engine
}

// GetAuthCodeRedisHelper 获取redis实例
func GetAuthCodeRedisHelper() (*redis.Client, context.Context) {
	dbOnce.Do(func() {
		newAuthCodeRedisHelper()
	})
	return authcodeRedisInstance, redisContext
}
func newAuthCodeRedisHelper() {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%v:%v",
			conf.Main.Redis.Host,
			conf.Main.Redis.Port,
		),
		Password: conf.Main.Redis.Pass, // no password set
		DB:       0,                    // use default DB
	})

	_, err := rdb.Ping(redisContext).Result()
	if err != nil {
		log.Fatal().Err(err).Msg("redis ping error")
		panic(err)
	}
	authcodeRedisInstance = rdb
}
