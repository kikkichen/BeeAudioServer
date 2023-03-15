package remote

import (
	"BeeAudioServer/models/netmodel"
	"BeeAudioServer/utils"
	"fmt"

	"github.com/go-resty/resty/v2"
)

/*	获取来自 net_cloud 音乐Url链接
 *	@param	request_client	resty网络i气逆供求客户端
 *	@param	id		歌曲ID (可以是多个)
 *	@param	level	音质
 */
func RequestAudioUrl(
	request_client *resty.Client,
	songIds string,
	level string,
) []netmodel.SongData {
	respMsg := new(netmodel.ResponseSongData)
	/* 请求目标歌曲的获取链接 */
	_, err := request_client.R().
		SetQueryParams(map[string]string{
			"id":    songIds,
			"level": level,
		}).
		SetResult(&respMsg).
		SetAuthToken(utils.NET_TOKEN).
		ForceContentType("application/json").
		SetJSONEscapeHTML(false).
		Get(utils.NET_ACCESS_ADDRESS + "/song/url/v1")

	/* 抛出异常 */
	if err != nil {
		panic(err)
	}

	/* 异常处理 */
	defer func(respMsg netmodel.ResponseSongData) {
		if err := recover(); err != nil {
			/* 发现错误不抛出， 返回信号， 待调用处理 */
			fmt.Printf("🔴 song_id:%v Url链接请求出错\n: %v\n", songIds, err)
		}
	}(*respMsg)

	return respMsg.Data
}
