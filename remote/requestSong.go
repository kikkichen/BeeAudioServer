package remote

import (
	"BeeAudioServer/models/netmodel"
	"BeeAudioServer/utils"
	"fmt"

	"github.com/go-resty/resty/v2"
)

/*	通过歌曲ID查询歌曲的详情信息（可查询多个）
 *	@param	request_client	resty网络i气逆供求客户端
 *	@param	SongIds	歌曲ID数组
 *
 */
func RequestSongDetail(
	request_client *resty.Client,
	SongIds string,
) []netmodel.SongInfo {
	var respMsg netmodel.ResponseBodyWrappedBody
	/* 请求曲目详情信息 */
	responseBody, err := request_client.R().
		SetQueryParams(map[string]string{
			"ids": SongIds,
		}).
		SetResult(&respMsg).
		SetAuthToken(utils.NET_TOKEN).
		ForceContentType("application/json").
		SetJSONEscapeHTML(false).
		Get(utils.LOCAL_NETCLOUD_ADDRESS + "/song/detail")

	/* 异常处理 */
	defer func(resp *resty.Response, err error) {
		if err := recover(); err != nil {
			/* 发现错误不抛出， 返回信号， 待调用处理 */
			fmt.Printf("🔴 曲目详情请求出错\n[请求体]:%v\n\n[报错]%v", resp, err)
		}
	}(responseBody, err)

	return respMsg.Songs
}
