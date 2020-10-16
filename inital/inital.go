package inital

import (
	"os"

	"github.com/rs/zerolog/log"
)

var file *os.File
var err error

func init() {
	setupMathSeed()
	file, err = getLogFile()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to open error log file")
	}
	log.Debug().Msg("log file opened")
	setupLog(file)
}

// Clean 获取初始化过程中出现的需要clean的东西
func Clean() {
	file.Close()
}
