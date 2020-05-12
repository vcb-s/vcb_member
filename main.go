package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	_ "vcb_member/inital"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"

	"vcb_member/conf"
	"vcb_member/helper"
	"vcb_member/models"
	"vcb_member/router"
)

func main() {
	file, err := os.OpenFile("log/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer file.Close()

	if err != nil {
		log.Error().Err(err).Msg("Failed to open error log file")
		panic("can not open log file")
	}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	log.Logger = log.
		Output(
			zerolog.MultiLevelWriter(
				zerolog.ConsoleWriter{
					Out:        os.Stderr,
					TimeFormat: time.RFC3339,
				},
				diode.NewWriter(file, 1000, 10*time.Millisecond, func(missed int) {
					log.Error().Msg("missed log: " + string(missed))
				}),
			),
		).
		With().
		Timestamp().
		Caller().
		Logger()

	log.Info().Msg("log test success")

	model := models.GetDBHelper()
	defer model.Close()

	if p, err := helper.CalcPassHash("0000"); err == nil {
		log.Info().Msg("example pass encode for 0000")
		log.Info().Msg(p)
	} else {
		log.Error().Err(err).Msg("calc pass test failed")
		return
	}

	addr := fmt.Sprintf(":%d", conf.Main.Server.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: router.Router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("start server fail")
		}
	}()

	// 退出监听
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shuting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server Shutdown With Error")
	}

	log.Info().Msg("Server Shutdown Success")
}
