package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"vcb_member/conf"
)

var form = url.Values{}

func init() {
	form.Set("grant_type", "authorization_code")
}

// AccessTokenResponse 主站返回的 accesstoken 响应体
type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

// AccessTokenResponseError 主站返回的 accesstoken 错误 响应体
type AccessTokenResponseError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// GetAccessTokenFromCode 用code换取accessToken
func GetAccessTokenFromCode(code string) (string, error) {
	hc := http.Client{}
	result := AccessTokenResponse{}
	errorResult := AccessTokenResponseError{}

	url := fmt.Sprintf(
		"https://%s:%s@vcb-s.com/oauth/token/",
		conf.Main.Third.Wp.ClientID,
		conf.Main.Third.Wp.ClientSec,
	)

	form.Set("code", code)

	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	// 去除BOM头
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))

	if resp.StatusCode != http.StatusOK {
		err = json.Unmarshal(body, &errorResult)
		if err != nil {
			return "", errors.New(string(resp.StatusCode))
		}

		return "", errors.New(errorResult.ErrorDescription)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", nil
	}
	fmt.Println(result)

	return result.AccessToken, nil
}

// UserAvatarURLS 主站用户头像结构
type UserAvatarURLS struct {
	Small  string `json:"24"`
	Medium string `json:"48"`
	Big    string `json:"96"`
}

// UserInfoResponse 主站返回的 userInfo 响应体
type UserInfoResponse struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	URL         string         `json:"url"`
	Description string         `json:"description"`
	Link        string         `json:"link"`
	Slug        string         `json:"slug"`
	Avatars     UserAvatarURLS `json:"avatar_urls"`
}

// UserInfoResponseError 主站返回的 userInfo 错误响应体
type UserInfoResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// GetUserInfoFromAccesstoken 用 accessToken 换取用户信息
func GetUserInfoFromAccesstoken(accessToken string) (UserInfoResponse, error) {
	hc := http.Client{}
	result := UserInfoResponse{}
	errorResult := UserInfoResponseError{}

	url := fmt.Sprintf(
		"https://vcb-s.com/wp-json/wp/v2/users/me?access_token=%s",
		accessToken,
	)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, err
	}

	resp, err := hc.Do(req)
	if err != nil {
		return result, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, nil
	}

	// 去除BOM头
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))

	if resp.StatusCode != http.StatusOK {
		err = json.Unmarshal(body, &errorResult)
		if err != nil {
			return result, errors.New(string(resp.StatusCode))
		}

		return result, errors.New(errorResult.Message)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, nil
	}
	fmt.Println(result)

	return result, nil
}
