package models

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"vcb_member/conf"
)

var dbInstance *gorm.DB
var dbOnce sync.Once

var authTokenStoreOnce sync.Once
var authTokenStore *badger.DB

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

/** 获取鉴权token store */
func GetAuthTokenStore() *badger.DB {
	authTokenStoreOnce.Do(func() {
		// newAuthCodeRedisHelper()
		db, err := badger.Open(badger.DefaultOptions("./kv/token"))
		if err != nil {
			log.Panic().Err(err).Msg("badger初始化失败")
		}

		authTokenStore = db
	})

	return authTokenStore
}
