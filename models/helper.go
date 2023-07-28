package models

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"vcb_member/conf"
)

var dbInstance *gorm.DB
var dbOnce sync.Once

var authcodeRedisOnce sync.Once
var authcodeRedisInstance *redis.Client
var authCodeRedisContext = context.Background()

// GetDBHelper 获取数据库实例
func GetDBHelper() *gorm.DB {
	dbOnce.Do(func() {
		initDBHelper()
	})
	return dbInstance
}

type customDBLogWriter struct{}

func (c customDBLogWriter) Printf(msg string, data ...interface{}) {
	log.Printf(msg, data...)
}

func initDBHelper() {
	dsn := fmt.Sprintf(
		"%v:%v@tcp([%v]:%v)/%v?charset=utf8mb4&parseTime=true&loc=Local",
		conf.Main.Database.User,
		conf.Main.Database.Pass,
		conf.Main.Database.Host,
		conf.Main.Database.Port,
		conf.Main.Database.Dbname,
	)

	customDBLogger := logger.New(
		customDBLogWriter{},
		logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Warn,
			Colorful:      true,
		},
	)

	if conf.Main.Debug {
		customDBLogger = customDBLogger.LogMode(logger.Info)
	}

	dbConfig := gorm.Config{
		Logger: customDBLogger,
	}

	engine, err := gorm.Open(mysql.Open(dsn), &dbConfig)
	if err != nil {
		log.Panic().Err(err).Msg("gorm open error")
	}

	sqlDB, err := engine.DB()
	if err != nil {
		log.Panic().Err(err).Msg("gorm get DB flag error")
	}

	//test DB if connection
	err = sqlDB.Ping()
	if err != nil {
		log.Panic().Err(err).Msg("gorm ping error")
	}

	//设置连接池
	sqlDB.SetMaxIdleConns(5)            //空闲数大小
	sqlDB.SetMaxOpenConns(50)           //最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) //重用超时

	dbInstance = engine

	log.Info().Msg("main db started")
}

// GetAuthCodeRedisHelper 获取redis实例
func GetAuthCodeRedisHelper() (*redis.Client, context.Context) {
	authcodeRedisOnce.Do(func() {
		newAuthCodeRedisHelper()
	})
	return authcodeRedisInstance, authCodeRedisContext
}
func newAuthCodeRedisHelper() {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%v:%v",
			conf.Main.Redis.Host,
			conf.Main.Redis.Port,
		),
		Password: conf.Main.Redis.Pass,
		DB:       0,
	})

	_, err := rdb.Ping(authCodeRedisContext).Result()
	if err != nil {
		log.Panic().Err(err).Msg("redis ping error")
	}
	log.Info().Msg("redis started")
	authcodeRedisInstance = rdb
}
