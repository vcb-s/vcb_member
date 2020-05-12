package inital

import (
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
)

// setupLog 获取log文件句柄
func setupLog(file *os.File) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	fileWritter := diode.NewWriter(file, 1000, 10*time.Millisecond, func(missed int) {
		log.Error().Msg("missed log: " + string(missed))
	})

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
		Timestamp().
		Caller().
		Logger()

	log.Info().Msg("log test success")

	gin.DisableConsoleColor()
	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
}
