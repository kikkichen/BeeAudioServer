package remote

import (
	"BeeAudioServer/models/netmodel"
	"BeeAudioServer/utils"
	"fmt"

	"github.com/go-resty/resty/v2"
)

/*	请求热门歌单分类Tag
 *	@param	request_client	resty网络i气逆供求客户端
 *
 */
func RequestPlayListHotTag(
	request_client *resty.Client,
) []netmodel.PlayListTag {
	var respMsg netmodel.ResponseHotTagsOutSideBody
	/* 请求播放列表Tag信息 */
	responseBody, err := request_client.R().
		SetResult(&respMsg).
		SetAuthToken(utils.NET_TOKEN).
		ForceContentType("application/json").
		SetJSONEscapeHTML(false).
		Get(utils.LOCAL_NETCLOUD_ADDRESS + "/playlist/hot")

	/* 异常处理 */
	defer func(resp *resty.Response, err error) {
		if err := recover(); err != nil {
			/* 发现错误不抛出， 返回信号， 待调用处理 */
			fmt.Printf("🔴 热门歌单Tag请求出错\n[请求体]:%v\n\n[报错]%v", resp, err)
		}
	}(responseBody, err)

	return respMsg.Tags
}

/* 请求全部播放列表Tag标签 */
func RequestPlayListAllTag(
	request_client *resty.Client,
) []netmodel.PlayListTag {
	var respMsg netmodel.ResponseAllTagsOutSideBody
	/* 请求播放列表信息 */
	responseBody, err := request_client.R().
		SetResult(&respMsg).
		SetAuthToken(utils.NET_TOKEN).
		ForceContentType("application/json").
		SetJSONEscapeHTML(false).
		Get(utils.LOCAL_NETCLOUD_ADDRESS + "/playlist/catlist")

	/* 异常处理 */
	defer func(resp *resty.Response, err error) {
		if err := recover(); err != nil {
			/* 发现错误不抛出， 返回信号， 待调用处理 */
			fmt.Printf("🔴 全部歌单Tag请求出错\n[请求体]:%v\n\n[报错]%v", resp, err)
		}
	}(responseBody, err)

	return respMsg.SubTagList
}
