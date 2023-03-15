package remote

import (
	"BeeAudioServer/models/netmodel"
	"BeeAudioServer/utils"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

/*	请求歌单详情
 *	@param	request_client	resty网络i气逆供求客户端
 *	@param	playlist_id	目标歌单请求ID
 */
func RequestPlayListDetail(
	request_client *resty.Client,
	playlist_id string,
) netmodel.PlayListInfo {
	var respMsg netmodel.ResponsePlayListOutSideBody
	/* 请求播放歌单详细信息 */
	responseBody, err := request_client.R().
		SetQueryParams(map[string]string{
			"id": playlist_id,
		}).
		SetResult(&respMsg).
		SetAuthToken(utils.NET_TOKEN).
		ForceContentType("application/json").
		SetJSONEscapeHTML(false).
		Get(utils.LOCAL_NETCLOUD_ADDRESS + "/playlist/detail")

	/* 异常处理 */
	defer func(resp *resty.Response, err error) {
		if err := recover(); err != nil {
			/* 发现错误不抛出， 返回信号， 待调用处理 */
			fmt.Printf("🔴 歌单详细信息请求出错\n[请求体]:%v\n\n[报错]%v", resp, err)
		}
	}(responseBody, err)

	return respMsg.PlayListInfo
}

/*	请求目标歌单所有曲目列表
 *	@param	request_client	resty网络i气逆供求客户端
 *	@param	playlist_id	目标歌单请求ID
 *	@param	page	页码
 *	@param	size	单页大小
 */
func RequestPlayListAllSong(
	request_client *resty.Client,
	playlist_id string,
	page int,
	size int,
) []netmodel.SongInfo {
	var respMsg netmodel.ResponsePlayListAllSongBody
	/* 请求播放歌单详细信息 */
	responseBody, err := request_client.R().
		SetQueryParams(map[string]string{
			"id":     playlist_id,
			"limit":  strconv.Itoa(size),
			"offset": strconv.Itoa((page - 1) * size),
		}).
		SetResult(&respMsg).
		SetAuthToken(utils.NET_TOKEN).
		ForceContentType("application/json").
		SetJSONEscapeHTML(false).
		Get(utils.LOCAL_NETCLOUD_ADDRESS + "/playlist/track/all")

	/* 异常处理 */
	defer func(resp *resty.Response, err error) {
		if err := recover(); err != nil {
			/* 发现错误不抛出， 返回信号， 待调用处理 */
			fmt.Printf("🔴 歌单全部歌曲请求出错\n[请求体]:%v\n\n[报错]%v", resp, err)
		}
	}(responseBody, err)

	return respMsg.Songs
}

/*	请求热门歌单列表
 *	@param	request_client	resty网络请求客户端对象
 */
func RequestTopPlayList(
	request_client *resty.Client,
) []netmodel.PlayListInfo {
	var respMsg netmodel.PlayListCollection
	/* 请求播放歌单详细信息 */
	responseBody, err := request_client.R().
		SetResult(&respMsg).
		SetAuthToken(utils.NET_TOKEN).
		ForceContentType("application/json").
		SetJSONEscapeHTML(false).
		Get(utils.LOCAL_NETCLOUD_ADDRESS + "/top/playlist/highquality")

	/* 异常处理 */
	defer func(resp *resty.Response, err error) {
		if err := recover(); err != nil {
			/* 发现错误不抛出， 返回信号， 待调用处理 */
			fmt.Printf("🔴 热门歌单列表信息请求出错\n[请求体]:%v\n\n[报错]%v", resp, err)
		}
	}(responseBody, err)

	return respMsg.List
}

/*	请求指定Tag/Cat歌单列表
 *	@param	request_client	resty网络请求客户端对象
 *	@param	cat 与之歌单匹配的Tag/Cat
 *	@param	page	页码
 *	@param	size	单页大小 - 默认容量为50
 */
func RequestTargetTagsPlayLists(
	request_client *resty.Client,
	cat string,
	page int,
	size int,
) []netmodel.PlayListInfo {
	var respMsg netmodel.PlayListCollection
	/* 请求播放歌单详细信息 */
	responseBody, err := request_client.R().
		SetQueryParams(map[string]string{
			"cat":    cat,
			"limit":  strconv.Itoa(size),
			"offset": strconv.Itoa((page - 1) * size),
		}).
		SetResult(&respMsg).
		SetAuthToken(utils.NET_TOKEN).
		ForceContentType("application/json").
		SetJSONEscapeHTML(false).
		Get(utils.LOCAL_NETCLOUD_ADDRESS + "/top/playlist/highquality")

	/* 异常处理 */
	defer func(resp *resty.Response, err error) {
		if err := recover(); err != nil {
			/* 发现错误不抛出， 返回信号， 待调用处理 */
			fmt.Printf("🔴 请求指定Tag歌单列表信息请求出错\n[请求体]:%v\n\n[报错]%v", resp, err)
		}
	}(responseBody, err)

	return respMsg.List
}
