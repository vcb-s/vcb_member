package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"vcb_member/helper"
	"vcb_member/models"
)

func main() {
	uid := helper.GenID()

	fmt.Println("models.Conf -------------------")
	fmt.Println(models.Conf)
	fmt.Println("uid -------------------")
	fmt.Println(uid)

	fmt.Println("GenToken -------------------")
	token, err := helper.GenToken(uid)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Println(token)

	fmt.Println("GenRefreshToken -------------------")
	token, err = helper.GenRefreshToken(uid)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Println(token)

	fmt.Println("GenCiphertext -------------------")
	result, err := helper.GenCiphertext(uid)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Println(result.Iv, result.Ciphertext)

	fmt.Println("GenHashtext -------------------")
	fmt.Println(helper.GenHashtext(uid))

	fmt.Println("GenIVByte -------------------")
	nonce, err := helper.GenIVByte()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	base64Nonce := base64.URLEncoding.EncodeToString(nonce)
	hexNonce := hex.EncodeToString(nonce)
	fmt.Println(base64Nonce, len(base64Nonce))
	fmt.Println(hexNonce, len(hexNonce))

	fmt.Println("GenPass -------------------")
	passUID, err := helper.GenPass(uid)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Println(passUID)
	fmt.Println("CheckPass -------------------")
	fmt.Println(helper.CheckPass(uid, passUID))
}
