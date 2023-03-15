package auth

import (
	"BeeAudioServer/utils"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type OathTokenBody struct {
	ClientId  string `json:"client_id"`
	Domain    string `json:"domain"`
	ExpiresIn int    `json:"expires_in"`
	Scope     string `json:"scope"`
	UserId    string `json:"user_id"`
}

/* 验证Token */
func VertifyOauth2Token(request_client *resty.Client, token string) (uint64, error) {
	/* 验证Token 请求响应体 */
	var oath_result OathTokenBody

	/* 待获取请求用户ID */
	var login_user_id uint64
	responseBody, err := request_client.R().
		SetAuthToken(token).
		SetResult(&oath_result).
		ForceContentType("application/json").
		SetJSONEscapeHTML(false).
		Get(utils.LOCAL_OAUTH2_SERVER + "/verify")

	/* 异常处理 */
	defer func(resp *resty.Response, err error) {
		if err := recover(); err != nil {
			/* 发现错误不抛出， 返回信号， 待调用处理 */
			fmt.Printf("🔴 验证用户登陆Token请求出错\n[请求体]:%v\n\n[报错]%v", resp, err)
		}
	}(responseBody, err)

	login_user_id = utils.StringParseToUint64(oath_result.UserId)
	return login_user_id, err
}
