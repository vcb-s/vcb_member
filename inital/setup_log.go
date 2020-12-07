package inital

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"

	"vcb_member/conf"
)

// setupLog 获取log文件句柄
func setupLog(file *os.File) {
	log.Debug().Bool("debug mode", conf.Main.Debug).Msg("current debug mode")

	if conf.Main.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}

	fileWritter := diode.NewWriter(file, 1000, 10*time.Millisecond, func(missed int) {
		log.Error().Msg("missed log: " + fmt.Sprint(missed))
	})

	if conf.Main.Debug {
		log.Logger = log.
			Output(
				zerolog.MultiLevelWriter(
					zerolog.ConsoleWriter{
						Out:        os.Stderr,
						TimeFormat: time.RFC3339,
					},
					fileWritter,
				),
			).
			With().
			Logger()
	} else {
		log.Logger = log.
			Output(
				zerolog.MultiLevelWriter(
					fileWritter,
				),
			).
			With().
			Logger()
	}

	if conf.Main.Debug {
		gin.DefaultWriter = io.MultiWriter(os.Stdout)
	} else {
		gin.DisableConsoleColor()
		gin.DefaultWriter = io.MultiWriter(file)
	}

	log.Debug().Msg("log setup success")
}
