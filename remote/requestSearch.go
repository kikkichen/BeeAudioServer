package remote

import (
	"BeeAudioServer/models/localmodel"
	"BeeAudioServer/models/netmodel"
	"BeeAudioServer/utils"
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*	请求 搜索结果
 *	@param	request_client	resty网络请求客户端对象
 *	@param	search_type	搜索类型
 *	@param	keyword	搜索关键字
 *	@param	page	页码
 *	@param	size	单页大小
 *
 */
func RequestSearchResult(
	request_client *resty.Client,
	search_type int,
	keywords string,
	page int,
	size int,
) any {
	var net_tail string = "/cloudsearch"
	/* 依据参数类型判断请求搜索类型 */
	if search_type == 10 {
		/* 请求专辑搜索结果 */
		var respMsg netmodel.ResponseSearchBody[netmodel.AblumResult]
		responseBody, err := request_client.R().
			SetQueryParams(map[string]string{
				"keywords": keywords,
				"type":     strconv.Itoa(search_type),
				"limit":    strconv.Itoa(size),
				"offset":   strconv.Itoa((page - 1) * size),
			}).
			SetResult(&respMsg).
			SetAuthToken(utils.NET_TOKEN).
			ForceContentType("application/json").
			SetJSONEscapeHTML(false).
			Get(utils.LOCAL_NETCLOUD_ADDRESS + net_tail)

		/* 异常处理 */
		defer func(resp *resty.Response, err error) {
			if err := recover(); err != nil {
				/* 发现错误不抛出， 返回信号， 待调用处理 */
				fmt.Printf("🔴 专辑搜索列表信息请求出错\n[请求体]:%v\n\n[报错]%v", resp, err)
			}
		}(responseBody, err)

		return respMsg.Result
	} else if search_type == 100 {
		/* 请求艺人搜索结果 */
		var respMsg netmodel.ResponseSearchBody[netmodel.ArtitstResult]
		responseBody, err := request_client.R().
			SetQueryParams(map[string]string{
				"keywords": keywords,
				"type":     strconv.Itoa(search_type),
				"limit":    strconv.Itoa(size),
				"offset":   strconv.Itoa((page - 1) * size),
			}).
			SetResult(&respMsg).
			SetAuthToken(utils.NET_TOKEN).
			ForceContentType("application/json").
			SetJSONEscapeHTML(false).
			Get(utils.LOCAL_NETCLOUD_ADDRESS + net_tail)

		/* 异常处理 */
		defer func(resp *resty.Response, err error) {
			if err := recover(); err != nil {
				/* 发现错误不抛出， 返回信号， 待调用处理 */
				fmt.Printf("🔴 艺人搜索列表信息请求出错\n[请求体]:%v\n\n[报错]%v", resp, err)
			}
		}(responseBody, err)

		return respMsg.Result
	} else if search_type == 1000 {
		/* 请求歌单搜索结果 */
		var respMsg netmodel.ResponseSearchBody[netmodel.PlayListResult]
		responseBody, err := request_client.R().
			SetQueryParams(map[string]string{
				"keywords": keywords,
				"type":     strconv.Itoa(search_type),
				"limit":    strconv.Itoa(size),
				"offset":   strconv.Itoa((page - 1) * size),
			}).
			SetResult(&respMsg).
			SetAuthToken(utils.NET_TOKEN).
			ForceContentType("application/json").
			SetJSONEscapeHTML(false).
			Get(utils.LOCAL_NETCLOUD_ADDRESS + net_tail)

		/* 异常处理 */
		defer func(resp *resty.Response, err error) {
			if err := recover(); err != nil {
				/* 发现错误不抛出， 返回信号， 待调用处理 */
				fmt.Printf("🔴 歌单搜索列表信息请求出错\n[请求体]:%v\n\n[报错]%v", resp, err)
			}
		}(responseBody, err)

		return respMsg.Result
	} else {
		/* 请求单曲搜索结果 */
		var respMsg netmodel.ResponseSearchBody[netmodel.SingleSongsResult]
		responseBody, err := request_client.R().
			SetQueryParams(map[string]string{
				"keywords": keywords,
				"type":     strconv.Itoa(search_type),
				"limit":    strconv.Itoa(size),
				"offset":   strconv.Itoa((page - 1) * size),
			}).
			SetResult(&respMsg).
			SetAuthToken(utils.NET_TOKEN).
			ForceContentType("application/json").
			SetJSONEscapeHTML(false).
			Get(utils.LOCAL_NETCLOUD_ADDRESS + net_tail)

		/* 异常处理 */
		defer func(resp *resty.Response, err error) {
			if err := recover(); err != nil {
				/* 发现错误不抛出， 返回信号， 待调用处理 */
				fmt.Printf("🔴 单曲搜索列表信息请求出错\n[请求体]:%v\n\n[报错]%v", resp, err)
			}
		}(responseBody, err)

		/* 匹配数据库过滤单曲信息 */
		if len(respMsg.Result.Songs) != 0 {
			db, err := gorm.Open(mysql.New(mysql.Config{
				DSN:                       utils.MYSQL_DSN,
				DefaultStringSize:         256,
				DisableDatetimePrecision:  true,
				DontSupportRenameIndex:    true,
				DontSupportRenameColumn:   true,
				SkipInitializeWithVersion: false,
			}), &gorm.Config{})

			if err != nil {
				log.Fatal(err)
			}
			respMsg.Result.Songs = FilterSongs(db, respMsg.Result.Songs)
		}

		return respMsg.Result
	}
}

/*	过滤获得有效歌单
 *	@param	db	gorm连接对象
 *	@param	songs	进行过滤的请求音频列表
 */
func FilterSongs(db *gorm.DB, songs []netmodel.SongInfo) []netmodel.SongInfo {
	var song_info_ids []uint64
	for _, song := range songs {
		song_info_ids = append(song_info_ids, song.SongId)
	}

	/* 数据库查询有音乐ID无记录 */
	var result_songs []localmodel.LocalSongModel
	db.Model(&localmodel.LocalSongModel{}).
		Where("id IN ?", song_info_ids).
		Find(&result_songs)

	/* 若存在可覆盖的记录 */
	if len(songs) > 0 {
		for index, item := range songs {
			songs[index].Source = "126.net"
			songs[index].Privilege = 0
			songs[index].Usable = true
			for _, song := range result_songs {
				if item.SongId == song.Id {
					songs[index].Source = song.Source
					songs[index].Privilege = song.Privilege
					songs[index].Usable = song.Useful
				}
			}
		}
	}
	return songs
}
