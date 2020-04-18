package helper

import (
	"vcb_member/models"
)

// UserCanGo 校验token用户是否可以修改指定uid的信息
func UserCanGo(uidInAuth string, uidInParam string) (CanGo bool) {
	var userInAuth models.User

	if uidInAuth == "" || uidInParam == "" {
		return false
	}

	if uidInAuth == uidInParam {
		return true
	}

	if err := models.GetDBHelper().First(&userInAuth, uidInAuth).Error; err != nil {
		return false
	}

	if userInAuth.IsAdmin != 1 {
		return false
	}

	return true
}
