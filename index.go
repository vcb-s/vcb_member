package main

import (
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"vcb_member/helper"
	"vcb_member/models"
)

func main() {
	fmt.Println("数据库连接 ---------")
	var total int
	_, err := models.GetDBHelper().SQL(`
		select count(id) from user_group;
	`).Get(&total)

	if err != nil {
		panic(err)
	}

	fmt.Println(total)

	fmt.Println("签发 refreshToken --------------")
	tokenString, err := helper.GenRefreshToken("hhvagrhhxd")
	token := []byte(tokenString)
	tokenString, err = helper.ReGenRefreshToken(token)
	fmt.Println(tokenString)
	uid, err := helper.CheckRefreshToken(token)
	if err != nil {
		panic(err)
	}
	if uid == "" {
		fmt.Println("检验失败")
		return
	}

	fmt.Println("校验成功")
	fmt.Println("新旧token：--------")
	fmt.Println(tokenString)
	fmt.Println("---------------")
	tokenString, err = helper.ReGenRefreshToken(token)
	if err != nil {
		panic(err)
	}
	fmt.Println(tokenString)

}
