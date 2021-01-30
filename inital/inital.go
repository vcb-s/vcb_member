package inital

import (
	"os"
	"time"
)

var file *os.File
var ticker *time.Ticker
var err error
var lastDate string

const dateFormat = "2006-01-02"

func init() {
	setupMathSeed()
	setupLog(getLogFile(time.Now().Format(dateFormat)))
	go func() {
		setupRotatingLog()
	}()
}

// 日志轮转
func setupRotatingLog() {
	lastDate = time.Now().Format(dateFormat)
	ticker = time.NewTicker(time.Minute)

	for range ticker.C {
		currentDate := time.Now().Format(dateFormat)
		if currentDate != lastDate {
			lastDate = currentDate
			setupLog(getLogFile(currentDate))
		}
	}
}

// Clean 获取初始化过程中出现的需要clean的东西
func Clean() {
	file.Close()
	ticker.Stop()
}
