package helper

import (
	"errors"
	mRand "math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/btnguyen2k/olaf"
	"github.com/go-redis/redis/v8"
	"github.com/matthewhartstonge/argon2"
	"github.com/pascaldekloe/jwt"
	"github.com/rs/zerolog/log"

	"vcb_member/conf"
	"vcb_member/models"
)

// timeStart jwt时间偏移
// new Number( new Date('Mon Dec 01 2010 00:00:00 GMT+0800') )
const timeStart = 1291132800000
const jwtIssuer = "vcb-member"

// 暂时有效期为一年，因为token全部内存指定了
const jwtExpires = 365 * 30 * 24 * time.Hour

var idGenerator *olaf.Olaf
var tokenSignKey = []byte(conf.Main.Jwt.Mac)

// ErrorExpired jwt过期
const ErrorExpired = "Expired"

// ErrorInvalid jwt无效
const ErrorInvalid = "Invalid"

var defaultArgon2 = argon2.DefaultConfig()

func init() {
	idGenerator = olaf.NewOlafWithEpoch(1, timeStart)
}

// GenID 获取一个雪花ID
func GenID() string {
	return idGenerator.Id64Ascii()
}

// GenCode 获取一个四位随机数
func GenCode() string {
	return strconv.Itoa(mRand.Intn(9999-999) + 999)
}

// 可用于生成密码的字符
var passbase = strings.Split("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM0987654321", "")
var passbaseLen = len(passbase)

// GenPass 获取一个六位随机大小写数字密码
func GenPass() string {
	pass := [6]string{}
	for idx := range pass {
		pass[idx] = passbase[mRand.Intn(passbaseLen)]
	}

	return strings.Join(pass[:], "")
}

// GenToken 获取一个jwt，不负责校验用户的合法性
func GenToken(uid string) (string, error) {
	var claims jwt.Claims
	claims.Issuer = jwtIssuer
	claims.ID = uid
	claims.KeyID = GenID()

	rdb, ctx := models.GetAuthCodeRedisHelper()
	if err := rdb.Set(ctx, claims.KeyID, claims.ID, 0).Err(); err != nil {
		log.Error().Err(err).Str("用户UID", claims.ID).Msg("token签发期间无法写入redis")
		return "", err
	}

	token, err := claims.HMACSign(jwt.HS256, tokenSignKey)
	if err != nil {
		log.Error().Err(err).Str("用户UID", claims.ID).Msg("token签发错误")
		return "", err
	}

	return string(token), nil
}

// CheckToken 检查jwt
func CheckToken(token []byte) (string, error) {
	claims, err := jwt.HMACCheck(token, tokenSignKey)
	if err != nil {
		return "", errors.New("token无效")
	}

	// 校验一次keyID
	rdb, ctx := models.GetAuthCodeRedisHelper()
	UIDInRedis, err := rdb.Get(ctx, claims.KeyID).Result()
	if err != nil {
		if err != redis.Nil {
			return "", errors.New("token无效")
		}

		log.Error().Err(err).Msg("redis执行错误")
		return "", errors.New("redis错误")
	}
	if claims.ID != UIDInRedis {
		return "", errors.New("token无效")
	}

	return claims.ID, nil
}

// CalcPassHash 获取一个安全的密码Hash
func CalcPassHash(pass string) (string, error) {
	result, err := defaultArgon2.HashRaw([]byte(pass))
	if err == nil {
		return string(result.Encode()), nil
	}

	return "", err
}

// CheckPassHash 校验密码
func CheckPassHash(pass string, hash string) bool {
	raw, err := argon2.Decode([]byte(hash))
	if err != nil {
		return false
	}

	ok, err := raw.Verify([]byte(pass))
	if err != nil {
		return false
	}

	return ok
}
