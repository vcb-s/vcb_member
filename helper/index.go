package helper

import (
	"crypto/ed25519"
	"fmt"

	"github.com/btnguyen2k/olaf"
	argon "github.com/dwin/goArgonPass"
)

// example from RFC 8037, appendix A.1
var testKeyEd25519Private = ed25519.PrivateKey([]byte{
	0x9d, 0x61, 0xb1, 0x9d, 0xef, 0xfd, 0x5a, 0x60,
	0xba, 0x84, 0x4a, 0xf4, 0x92, 0xec, 0x2c, 0xc4,
	0x44, 0x49, 0xc5, 0x69, 0x7b, 0x32, 0x69, 0x19,
	0x70, 0x3b, 0xac, 0x03, 0x1c, 0xae, 0x7f, 0x60,
	// public key suffix
	0xd7, 0x5a, 0x98, 0x01, 0x82, 0xb1, 0x0a, 0xb7,
	0xd5, 0x4b, 0xfe, 0xd3, 0xc9, 0x64, 0x07, 0x3a,
	0x0e, 0xe1, 0x72, 0xf3, 0xda, 0xa6, 0x25, 0x93,
	0xaf, 0x02, 0x1a, 0x68, 0xf7, 0x07, 0x51, 0x1a,
})

// example from RFC 8037, appendix A.1
var testKeyEd25519Public = ed25519.PublicKey([]byte{
	0xd7, 0x5a, 0x98, 0x01, 0x82, 0xb1, 0x0a, 0xb7,
	0xd5, 0x4b, 0xfe, 0xd3, 0xc9, 0x64, 0x07, 0x3a,
	0x0e, 0xe1, 0x72, 0xf3, 0xda, 0xa6, 0x23, 0x25,
	0xaf, 0x02, 0x1a, 0x68, 0xf7, 0x07, 0x51, 0x1a,
})

// GenID 获取一个雪花ID
func GenID() {
	// new Number( new Date('Mon Dec 01 2010 00:00:00 GMT+0800') )
	o := olaf.NewOlafWithEpoch(1, 1291132800000)

	id64 := o.Id64()
	id64Hex := o.Id64Hex()
	id64Ascii := o.Id64Ascii()
	fmt.Println("ID 64-bit (int)   : ", id64, " / Timestamp: ", o.ExtractTime64(id64))
	fmt.Println("ID 64-bit (hex)   : ", id64Hex, " / Timestamp: ", o.ExtractTime64Hex(id64Hex))
	fmt.Println("ID 64-bit (ascii) : ", id64Ascii, " / Timestamp: ", o.ExtractTime64Ascii(id64Ascii))

	id128 := o.Id128()
	id128Hex := o.Id128Hex()
	id128Ascii := o.Id128Ascii()
	fmt.Println("ID 128-bit (int)  : ", id128.String(), " / Timestamp: ", o.ExtractTime128(id128))
	fmt.Println("ID 128-bit (hex)  : ", id128Hex, " / Timestamp: ", o.ExtractTime128Hex(id128Hex))
	fmt.Println("ID 128-bit (ascii): ", id128Ascii, " / Timestamp: ", o.ExtractTime128Ascii(id128Ascii))
}

// GenToken 获取一个jwt
func GenToken() {
	// var claims jwt.Claims
	// now := time.Now().Round(time.Second)
	// claims.Issuer = "inori"
	// claims.Issued = NewNumericTime(now)
	// claims.Expires = NewNumericTime(now.Add(10 * time.Minute))

	// token, err := claims.EdDSASign(privateKey)
	// if err != nil {
	// 	return "", err
	// }
}

// GenPass 获取一个安全的密码Hash
func GenPass(pass string) (string, error) {
	return argon.Hash(pass)
}

// CheckPass 校验密码
func CheckPass(pass string, hash string) error {
	return argon.Verify(pass, hash)
}
