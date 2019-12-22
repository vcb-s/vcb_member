package main

import (
	"fmt"
	"net/http"
	"vcb_member/conf"
	"vcb_member/router"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	addr := fmt.Sprintf(":%d", conf.Main.Server.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: router.Router,
	}

	go func() {
		// service connections
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println(err.Error())
		}
	}()
}
