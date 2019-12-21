package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"time"
	"vcb_member/models"

	"github.com/btnguyen2k/olaf"
	argon "github.com/dwin/goArgonPass"
	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/blake2b"
)

// timeStart jwt时间偏移
// new Number( new Date('Mon Dec 01 2010 00:00:00 GMT+0800') )
const timeStart = 1291132800000
const jwtIssuer = "vcb-member"
const jwtExpires = 5 * time.Minute
const jwtRefreshExpires = 7 * 24 * time.Hour

var signKey = []byte(models.Conf.Jwt.Mac)
var encryptKey = []byte(models.Conf.Jwt.Encryption)

// GenID 获取一个雪花ID
func GenID() string {
	o := olaf.NewOlafWithEpoch(1, timeStart)

	return o.Id64Ascii()
}

// GenToken 获取一个jwt
func GenToken(uid string) (string, error) {
	var claims jwt.Claims
	now := time.Now().Round(time.Second)
	claims.Issuer = jwtIssuer
	claims.Issued = jwt.NewNumericTime(now)
	claims.Expires = jwt.NewNumericTime(now.Add(jwtExpires))
	claims.Audiences = []string{uid}

	token, err := claims.HMACSign(jwt.HS256, signKey)
	if err != nil {
		return "", err
	}

	return string(token), nil
}

// GenRefreshToken 获取一个 refreshjwt
func GenRefreshToken(uid string) (string, error) {
	var claims jwt.Claims
	now := time.Now().Round(time.Second)
	tokenID := GenID()
	// 记录到用户的token字段，refresh的时候读取验证

	claims.ID = tokenID
	claims.Issuer = jwtIssuer
	claims.Issued = jwt.NewNumericTime(now)
	claims.Expires = jwt.NewNumericTime(now.Add(jwtRefreshExpires))
	claims.Audiences = []string{uid}

	token, err := claims.HMACSign(jwt.HS256, signKey)
	if err != nil {
		return "", err
	}

	return string(token), nil
}

// GenPass 获取一个安全的密码Hash
func GenPass(pass string) (string, error) {
	return argon.Hash(pass)
}

// CheckPass 校验密码
func CheckPass(pass string, hash string) bool {
	if argon.Verify(pass, hash) != nil {
		return false
	}

	return true
}

// GenIVByte 产生随机IV
func GenIVByte() ([]byte, error) {
	nonce := make([]byte, 12)
	io.ReadFull(rand.Reader, nonce)
	_, err := io.ReadFull(rand.Reader, nonce)

	return nonce, err
}

// GenHashtext 加密给定字符串
func GenHashtext(plaintext string) string {
	plainByte := []byte(plaintext)
	hash := blake2b.Sum256(plainByte)
	return base64.URLEncoding.EncodeToString(hash[:])
}

// Ciphertext 加密结果
type Ciphertext struct {
	Iv         string
	Ciphertext string
}

// GenCiphertext 加密给定字符串
func GenCiphertext(plaintext string) (Ciphertext, error) {
	plainByte := []byte(plaintext)
	var result Ciphertext

	iv, err := GenIVByte()
	if err != nil {
		panic(err)
	}

	result.Iv = base64.URLEncoding.EncodeToString(iv)

	block, err := aes.NewCipher(encryptKey)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	result.Ciphertext = base64.URLEncoding.EncodeToString(aesgcm.Seal(nil, iv, plainByte, nil))

	return result, err
}
