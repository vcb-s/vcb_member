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
