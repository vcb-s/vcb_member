package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vcb_member/conf"
	"vcb_member/inital"

	"github.com/rs/zerolog/log"

	"vcb_member/models"
	"vcb_member/router"
)

func main() {
	defer inital.Clean()

	// 配置数据库连接注销
	model := models.GetDBHelper()
	sqlDB, err := model.DB()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to open get sqlDB in defer close")
	}
	defer sqlDB.Close()

	// 配置redis连接注销
	rbd, _ := models.GetAuthCodeRedisHelper()
	defer rbd.Close()

	addr := fmt.Sprintf(":%d", conf.Main.Server.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: router.Router,
	}

	ginStartErrorDetect, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go func() {
		// server.ListenAndServe 会一直占用该线程
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("start server fail")
			cancel()
		}
	}()

	<-ginStartErrorDetect.Done()
	if ginStartErrorDetect.Err() == context.DeadlineExceeded {
		log.Debug().Msg("gin started")
	}

	// 退出监听
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Debug().Msg("Shuting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Server Shutdown With Error")
		cancel()
	}

	<-shutdownCtx.Done()
	if shutdownCtx.Err() == context.DeadlineExceeded {
		log.Debug().Msg("Server Shutdown Success")
	}
}
