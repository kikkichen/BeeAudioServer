package netmodel

/*	艺人信息获取
 *	net服务接口 ： /artist/detail
 *
 */

/* 响应体外部 */
type RequestArtistOutSideBody struct {
	Code    int                   `json:"code"`
	Message string                `json:"message"`
	Data    ArtistInfoWrappedBody `json:"data"`
}

/* 艺人信息 */
type ArtistInfo struct {
	ArtistId  uint64 `json:"id"`
	Name      string `json:"name"`
	Cover     string `json:"cover"`     // 艺人封面url
	PicUrl    string `json:"picUrl"`    // 艺人封面url - 用于搜索页条目
	BriefDesc string `json:"briefDesc"` //	艺人描述
	AlbumSize int    `json:"albumSize"` // 专辑数量
	MusicSize int    `json:"musicSize"` // 乐曲数量
}

type ArtistInfoWrappedBody struct {
	Artist ArtistInfo `json:"artist"`
}

/* 艺人曲目信息列表 */
type ArtistArtsBody struct {
	Songs []SongInfo `json:"songs"`
	More  bool       `json:"true"`
	Total int        `json:"total"`
	Code  int        `json:"code"`
}

/* 艺人炸u年纪信息列表 */
type ArtistAlbumBody struct {
	Albums []AlbumInfo `json:"hotAlbums"`
	More   bool        `json:"true"`
	Code   int         `json:"code"`
}

/* 艺人请求响应体样式：

{
	"code": 200,
	"message": "ok",
	"data": {
		"videoCount": 27,
		"artist": {
			"id": 855508,
			"cover": "http://p2.music.126.net/HteycXYRw4rEUGPNyKJEaQ==/109951166501599564.jpg",
			"name": "Daoko",
			"transNames": [],
			"identities": [],
			"identifyTag": null,
			"briefDesc": "DAOKO（1997年3月4日－），日本女性饶舌歌手，于东京都出道。",
			"rank": {
				"rank": 15,
				"type": 4
			},
			"albumSize": 31,
			"musicSize": 224,
			"mvSize": 27
		},
		"blacklist": true,
		"preferShow": 5,
		"showPriMsg": false,
		"secondaryExpertIdentiy": [
			{
				"expertIdentiyId": 5,
				"expertIdentiyName": "演唱",
				"expertIdentiyCount": 224
			},
			{
				"expertIdentiyId": 6,
				"expertIdentiyName": "作词",
				"expertIdentiyCount": 224
			},
			{
				"expertIdentiyId": 7,
				"expertIdentiyName": "作曲",
				"expertIdentiyCount": 224
			}
		]
	}
}














*/
