package netmodel

/* 歌单详情 */
type PlayListInfo struct {
	Id          uint64     `json:"id"`
	Name        string     `json:"name"`
	Cover       string     `json:"coverImgUrl"`
	UserId      uint64     `json:"userId"`
	CreateTime  uint64     `json:"createTime"`
	Description string     `json:"description"`
	Tags        []string   `json:"tags"`
	CreatorUser Creator    `json:"creator"`
	PlayList    []SongInfo `json:"tracks"`
	PlayListIds []TrackId  `json:"trackIds"`
}

/* 播放列表创建用户 */
type Creator struct {
	Id            uint64 `json:"userId"`
	NickName      string `json:"nickname"`
	Signatrue     string `json:"signatrue"`
	Description   string `json:"description"`
	AvatarUrl     string `json:"avatarUrl"`
	BackgroundUrl string `json:"backgroundUrl"`
}

/* 歌单详情信息 响应体外围 */
type ResponsePlayListOutSideBody struct {
	Code         string       `json:"code"`
	PlayListInfo PlayListInfo `json:"playlist"`
}

/* 获取歌单所有歌曲 响应体外围 */
type ResponsePlayListAllSongBody struct {
	Songs []SongInfo `json:"songs"`
}

/* 歌单中群不歌曲的ID */
type TrackId struct {
	SongId uint64 `json:"id"`
}
