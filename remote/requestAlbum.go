package remote

import (
	"BeeAudioServer/models/netmodel"
	"BeeAudioServer/utils"
	"fmt"

	"github.com/go-resty/resty/v2"
)

/*	通过专辑ID获取专辑详情
 *	@param	request_client	resty网络i气逆供求客户端
 *	@param	AlbumId	专辑ID
 *
 */
func RequestAlbumDetail(
	request_client *resty.Client,
	AlbumId string,
) netmodel.AlbumInfoWithSongList {
	var respMsg netmodel.AlbumInfoWithSongList
	/* 请求专辑信息 */
	responseBody, err := request_client.R().
		SetQueryParams(map[string]string{
			"id": AlbumId,
		}).
		SetResult(&respMsg).
		SetAuthToken(utils.NET_TOKEN).
		ForceContentType("application/json").
		SetJSONEscapeHTML(false).
		Get(utils.LOCAL_NETCLOUD_ADDRESS + "/album")

	/* 异常处理 */
	defer func(resp *resty.Response, err error) {
		if err := recover(); err != nil {
			/* 发现错误不抛出， 返回信号， 待调用处理 */
			fmt.Printf("🔴 专辑信息请求出错\n[请求体]:%v\n\n[报错]%v", resp, err)
		}
	}(responseBody, err)

	return respMsg
}
