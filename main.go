package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	_ "vcb_member/inital"

	"vcb_member/conf"
	"vcb_member/helper"
	"vcb_member/models"
	"vcb_member/router"
)

func main() {
	model := models.GetDBHelper()
	defer model.Close()

	if p, err := helper.CalcPassHash("12345678"); err == nil {
		fmt.Println("example pass encode for 12345678", p)
	} else {
		fmt.Println(err)
		return
	}

	addr := fmt.Sprintf(":%d", conf.Main.Server.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: router.Router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(err.Error())
		}
	}()

	// 退出监听
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shuting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown With Error: ", err)
	}

	log.Println("Server Shutdown Success")
}
