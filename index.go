package main

import (
	"fmt"
	"net/http"
	"vcb_member/models"
	"vcb_member/router"
)

func main() {
	addr := fmt.Sprintf(":%d", models.Conf.Server.Port)
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
