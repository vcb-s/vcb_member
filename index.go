package main

import (
	"vcb_member/helper"
	"vcb_member/models"

	"fmt"
)

func init() {
	// 解析toml得到env
}

func main() {
	// 解析本地toml
	// var config Songs
	// _, err := toml.DecodeFile("config.toml", &config)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(config)

	helper.GenID()

	fmt.Print(models.Conf)

	helper.GenID()
	// r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	// r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
