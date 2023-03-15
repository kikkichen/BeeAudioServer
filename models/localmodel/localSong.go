package localmodel

import (
	"BeeAudioServer/models/netmodel"
	"BeeAudioServer/utils"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

/* 本地歌曲数据信息 模型 */
type LocalSongModel struct {
	Id        uint64 `gorm:"column:id"`
	Name      string `gorm:"column:name"`
	ArtistId  uint64 `gorm:"column:ar_id"`
	AlbumId   uint64 `gorm:"column:al_id"`
	LocalPath string `gorm:"column:local_path"`
	Privilege int    `gorm:"column:privilege"`
	Quality   string `gorm:"column:quality"`
	Useful    bool   `gorm:"column:useful"`
	Source    string `gorm:"column:source"`
}

func (*LocalSongModel) TableName() string {
	return "song_table"
}

/*	将数据库中gorm模型类型的 Song， 映射转换请求体类型的SongInfo类型
 *	@param	client	resty网络连接客户端对象
 */
func (this *LocalSongModel) MapToSongInfoType(client *resty.Client) netmodel.SongInfo {
	var respArtist netmodel.RequestArtistOutSideBody
	var respAlbum netmodel.AlbumInfoWithSongList
	/* 请求艺人信息 */
	responseBody, err := client.R().
		SetQueryParams(map[string]string{
			"id": strconv.FormatUint(this.ArtistId, 10),
		}).
		SetResult(&respArtist).
		SetAuthToken(utils.NET_TOKEN).
		ForceContentType("application/json").
		SetJSONEscapeHTML(false).
		Get(utils.LOCAL_NETCLOUD_ADDRESS + "/artist/detail")

	/* 请求专辑信息 */
	responseBody, err = client.R().
		SetQueryParams(map[string]string{
			"id": strconv.FormatUint(this.AlbumId, 10),
		}).
		SetResult(&respAlbum).
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

	return netmodel.SongInfo{
		SongId: this.Id,
		Name:   this.Name,
		Ar: []netmodel.InnerSongInfoAr{{
			ArId: respArtist.Data.Artist.ArtistId,
			Name: respArtist.Data.Artist.Name,
		}},
		Al: netmodel.InnerSongInfoAl{
			AlId:   respAlbum.Album.AlbumId,
			Name:   respAlbum.Album.AlbumName,
			PicUrl: respAlbum.Album.PicUrl,
		},
		Source:    this.Source,
		Privilege: this.Privilege,
		// Usable: this.Useful,
	}
}
