package main

import (
	"fmt"
	"github.com/vcb_member_be/models"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

func main() {
	// 解析本地toml
	var config models.Songs
	_, err := toml.DecodeFile("config.toml", &config)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(config)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
