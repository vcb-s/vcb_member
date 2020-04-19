package helper

import (
	"errors"
	"time"
	"vcb_member/conf"
	"vcb_member/models"

	"github.com/btnguyen2k/olaf"
	argon "github.com/dwin/goArgonPass"
	"github.com/pascaldekloe/jwt"
)

// timeStart jwt时间偏移
// new Number( new Date('Mon Dec 01 2010 00:00:00 GMT+0800') )
const timeStart = 1291132800000
const jwtIssuer = "vcb-member"

// 暂时有效期为一年，因为token全部内存指定了
const jwtExpires = 365 * 30 * 24 * time.Hour

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

// getUserTokenID 获取用户对应的tokenID，以此保持登录token稳定
func getUserTokenID(uid string) string {
	tokenID := Session.SearchByValue(AuthTokenNamespace, uid)

	if len(tokenID) == 0 {
		// 查询一次数据库，不考虑用户不存在的情况
		var user models.User
		models.GetDBHelper().First(&user, "id = ?", uid)

		if len(user.LastTokenID) == 0 {
			tokenID = GenID()
			user.LastTokenID = tokenID
			// 不处理错误，暂时没看到抛出的价值
			// 即使出现错误也要王Session中存一个tokenID，防止击穿导致疯狂查数据库
			models.GetDBHelper().Model(&user).Updates(&user)
		} else {
			tokenID = user.LastTokenID
		}
	}

	Session.Set(AuthTokenNamespace, tokenID, uid)

	return tokenID
}

// GenToken 获取一个jwt，不负责校验用户的合法性
func GenToken(uid string) (string, error) {
	var claims jwt.Claims
	now := time.Now().Round(time.Second)
	claims.Issuer = jwtIssuer
	claims.ID = uid
	claims.KeyID = getUserTokenID(uid)
	claims.Issued = jwt.NewNumericTime(now)
	claims.Expires = jwt.NewNumericTime(now.Add(jwtExpires))

	token, err := claims.HMACSign(jwt.HS256, tokenSignKey)
	if err != nil {
		return "", err
	}

	return string(token), nil
}

// CheckToken 检查jwt
func CheckToken(token []byte) (string, error) {
	claims, err := jwt.HMACCheck(token, tokenSignKey)
	if err != nil {
		return "", err
	}

	if !claims.Valid((time.Now())) {
		return "", errors.New(ErrorExpired)
	}

	sessionTokenID := getUserTokenID(claims.ID)

	// 校验一次keyID
	if sessionTokenID != claims.KeyID {
		return "", errors.New(ErrorInvalid)
	}

	return claims.ID, nil
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
