package helper

import (
	"errors"
	"time"
	"vcb_member/conf"

	"github.com/btnguyen2k/olaf"
	argon "github.com/dwin/goArgonPass"
	"github.com/pascaldekloe/jwt"
)

// timeStart jwt时间偏移
// new Number( new Date('Mon Dec 01 2010 00:00:00 GMT+0800') )
const timeStart = 1291132800000
const jwtIssuer = "vcb-member"
const jwtExpires = 30 * time.Minute

var idGenerator *olaf.Olaf
var tokenSignKey = []byte(conf.Main.Jwt.Mac)
var refreshTokenSignKey = []byte(conf.Main.Jwt.Encryption)

// ErrorExpired jwt过期
const ErrorExpired = "Expired"

// ErrorInvalid jwt无效
const ErrorInvalid = "Invalid"

func init() {
	idGenerator = olaf.NewOlafWithEpoch(1, timeStart)
}

// GenID 获取一个雪花ID
func GenID() string {
	return idGenerator.Id64Ascii()
}

// GenToken 获取一个jwt
func GenToken(uid string) (string, error) {
	var claims jwt.Claims
	now := time.Now().Round(time.Second)
	claims.Issuer = jwtIssuer
	claims.Issued = jwt.NewNumericTime(now)
	claims.Expires = jwt.NewNumericTime(now.Add(jwtExpires))
	claims.Subject = uid

	token, err := claims.HMACSign(jwt.HS256, tokenSignKey)
	if err != nil {
		return "", err
	}

	Session.Set(AuthToken, uid, "1")

	return string(token), nil
}

// CheckToken 检查jwt
func CheckToken(tokenString []byte) (string, error) {
	claims, err := jwt.HMACCheck(tokenString, tokenSignKey)
	if err != nil {
		return "", err
	}

	if !claims.Valid((time.Now())) {
		return "", errors.New(ErrorExpired)
	}

	if !Session.Has(AuthToken, claims.Subject) {
		return "", errors.New(ErrorInvalid)
	}

	return claims.Subject, nil
}

// CalcPassHash 获取一个安全的密码Hash
func CalcPassHash(pass string) (string, error) {
	return argon.Hash(pass)
}

// CheckPassHash 校验密码
func CheckPassHash(pass string, hash string) bool {
	err := argon.Verify(pass, hash)
	if err != nil {
		return false
	}

	return true
}
