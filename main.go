package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"

	_ "vcb_member/inital"

	"vcb_member/conf"
	"vcb_member/router"
)

func main() {
	addr := fmt.Sprintf(":%d", conf.Main.Server.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: router.Router,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println(err.Error())
		}
	}()

	// 退出监听
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
