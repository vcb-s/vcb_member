package inital

import (
	"os"

	"github.com/rs/zerolog/log"
)

// getLogFile 获取log文件句柄
func getLogFile(currentDate string) *os.File {
	logFile, logFileError := os.OpenFile("log/"+currentDate+".txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if logFileError != nil {
		log.Panic().Err(err).Str("currentDate", currentDate).Msg("Failed to open error log file")

	}
	return logFile
}
