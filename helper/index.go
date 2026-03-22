package helper

import (
	"crypto/aes"
	"crypto/cipher"
	cRand "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	mRand "math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/btnguyen2k/consu/olaf"
	badger "github.com/dgraph-io/badger/v4"
	"github.com/matthewhartstonge/argon2"

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

var encKey []byte

func getEncKey() ([]byte, error) {
	if len(encKey) == 0 {
		key, err := hex.DecodeString(conf.Main.Token.Key)
		if err != nil {
			return nil, err
		}

		if len(key) != (256 / 8) {
			return nil, errors.New("key should be 32 bytes len")
		}

		encKey = key
	}

	return encKey, nil
}

// enc aes-256 encrypt
func enc(plaintext []byte) ([]byte, error) {
	key, err := getEncKey()
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce, err := genNonce(aead.NonceSize())
	if err != nil {
		return nil, err
	}

	return aead.Seal( /** 让 密文 append 在 nonce 后边，简化书写 */ nonce, nonce, plaintext, nil), nil
}

// dec aes-256 decrypt
func dec(ciphertext []byte) ([]byte, error) {
	key, err := getEncKey()
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aead.NonceSize()

	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	return aead.Open(nil, ciphertext[:nonceSize], ciphertext[nonceSize:], nil)
}

// genNonce 产生一个AES nonce
func genNonce(nonceSize int) ([]byte, error) {
	nonce := make([]byte, nonceSize)
	_, err := cRand.Read(nonce)

	if err != nil {
		return nil, err
	}

	return nonce, nil
}

// GenToken 签署一个token
func GenToken(uidStr string) (string, error) {
	uid := []byte(uidStr)

	token, err := enc(uid)
	if err != nil {
		return "", err
	}

	tokenStore := models.GetAuthTokenStore()
	err = tokenStore.Update(func(txn *badger.Txn) error {
		return txn.Set(uid, token)
	})

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.Strict().EncodeToString(token), nil
}

// CheckToken 校验一个token
func CheckToken(tokenStr string) (string, error) {
	ciphertext, err := base64.URLEncoding.Strict().DecodeString(tokenStr)
	if err != nil {
		return "", errors.New("token解码失败")
	}

	uid, err := dec(ciphertext)
	if err != nil {
		return "", errors.New("token无效")
	}

	tokenStore := models.GetAuthTokenStore()
	err = tokenStore.View(func(txn *badger.Txn) error {
		item, err := txn.Get(uid)
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			if reflect.DeepEqual(val, ciphertext) {
				return nil
			}

			return errors.New("token已失效")
		})

		return err
	})

	if err != nil {
		return "", err
	}

	return string(uid), nil
}

// CalcPassHash 获取一个安全的密码Hash
func CalcPassHash(pass string) (string, error) {
	result, err := defaultArgon2.HashRaw([]byte(pass))
	if err != nil {
		return "", err
	}

	return string(result.Encode()), nil
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

// /** GenCRandKey 生成指定长度的密码学随机key */
// func GenCRandKey(byteLen int) string {
// 	r := make([]byte, byteLen)
// 	_, err := cRand.Read(r)

// 	if err != nil {
// 		log.Panic().Err(err)
// 	}

// 	return hex.EncodeToString(r)
// }
