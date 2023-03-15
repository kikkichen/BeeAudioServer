package remote

import (
	"BeeAudioServer/models/netmodel"
	"BeeAudioServer/utils"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

/*	通过艺人ID获取艺人详情 （包含艺人名称、WIKI描述、曲目数目以及专辑数目）
 *	@param	request_client	resty网络i气逆供求客户端
 *	@param	ArtistId	艺人ID
 *
 */
func RequestArtistDetail(
	request_client *resty.Client,
	ArtistId string,
) netmodel.ArtistInfo {
	var respMsg netmodel.RequestArtistOutSideBody
	/* 请求艺人信息 */
	responseBody, err := request_client.R().
		SetQueryParams(map[string]string{
			"id": ArtistId,
		}).
		SetResult(&respMsg).
		SetAuthToken(utils.NET_TOKEN).
		ForceContentType("application/json").
		SetJSONEscapeHTML(false).
		Get(utils.LOCAL_NETCLOUD_ADDRESS + "/artist/detail")

	/* 异常处理 */
	defer func(resp *resty.Response, err error) {
		if err := recover(); err != nil {
			/* 发现错误不抛出， 返回信号， 待调用处理 */
			fmt.Printf("🔴 艺人信息请求出错\n[请求体]:%v\n\n[报错]%v", resp, err)
		}
	}(responseBody, err)

	return respMsg.Data.Artist
}

/*	通过艺人ID获取艺人曲目作品列表 （分页）
 *	@param	request_client	resty网络i气逆供求客户端
 *	@param	ArtistId	艺人ID
 *	@param	page	页码
 *	@param	size	单页容量大小
 *
 */
func RequestArtistSongs(
	request_client *resty.Client,
	ArtistId string,
	page int,
	size int,
) []netmodel.SongInfo {
	var respMsg netmodel.ArtistArtsBody
	/* 请求艺人曲目信息 */
	responseBody, err := request_client.R().
		SetQueryParams(map[string]string{
			"id":     ArtistId,
			"limit":  strconv.Itoa(size),
			"offset": strconv.Itoa((page - 1) * size),
		}).
		SetResult(&respMsg).
		SetAuthToken(utils.NET_TOKEN).
		ForceContentType("application/json").
		SetJSONEscapeHTML(false).
		Get(utils.LOCAL_NETCLOUD_ADDRESS + "/artist/songs")

	/* 异常处理 */
	defer func(resp *resty.Response, err error) {
		if err := recover(); err != nil {
			/* 发现错误不抛出， 返回信号， 待调用处理 */
			fmt.Printf("🔴 艺人创作曲目列表请求出错\n[请求体]:%v\n\n[报错]%v", resp, err)
		}
	}(responseBody, err)

	return respMsg.Songs
}

/*	通过艺人ID获取艺人专辑信息 （分页）
 *	@param	request_client	resty网络i气逆供求客户端
 *	@param	ArtistId	艺人ID
 *	@param	page	页码
 *	@param	size	单页容量大小
 *
 */
func RequestArtistAlbums(
	request_client *resty.Client,
	ArtistId string,
	page int,
	size int,
) []netmodel.AlbumInfo {
	var respMsg netmodel.ArtistAlbumBody
	/* 请求艺人曲目信息 */
	responseBody, err := request_client.R().
		SetQueryParams(map[string]string{
			"id":     ArtistId,
			"limit":  strconv.Itoa(size),
			"offset": strconv.Itoa((page - 1) * size),
		}).
		SetResult(&respMsg).
		SetAuthToken(utils.NET_TOKEN).
		ForceContentType("application/json").
		SetJSONEscapeHTML(false).
		Get(utils.LOCAL_NETCLOUD_ADDRESS + "/artist/album")

	/* 异常处理 */
	defer func(resp *resty.Response, err error) {
		if err := recover(); err != nil {
			/* 发现错误不抛出， 返回信号， 待调用处理 */
			fmt.Printf("🔴 艺人专辑列表请求出错\n[请求体]:%v\n\n[报错]%v", resp, err)
		}
	}(responseBody, err)

	return respMsg.Albums
}
