package inital

import (
	"os"
)

var logFile *os.File
var logFileError error

// getLogFile 获取log文件句柄
func getLogFile() (*os.File, error) {
	if logFile == nil {
		logFile, logFileError = os.OpenFile("log/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	}
	return logFile, logFileError
}
