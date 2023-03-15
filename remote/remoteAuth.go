package remote

import (
	"BeeAudioServer/utils"
	"log"

	"github.com/go-resty/resty/v2"
)

/* 服务端验证Token返回体 */
type OauthBody struct {
	ClientId  string `json:"client_id"`
	Domain    string `json:"domain"`
	ExpiresIn int    `json:"expires_in"`
	Scope     string `json:"scope"`
	UserId    string `json:"user_id"`
}

func AuthToken(client *resty.Client, token string) uint64 {
	var auth_info OauthBody
	response, err := client.R().
		SetAuthToken(token).
		SetResult(&auth_info).
		ForceContentType("application/json").
		SetJSONEscapeHTML(false).
		Get(utils.LOCAL_OAUTH2_SERVER + "/verify")

	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode() == 200 {
		return utils.StringParseToUint64(auth_info.UserId)
	} else {
		return 0
	}
}
