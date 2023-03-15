package netmodel

/* 专辑信息 */

type AlbumInfo struct {
	AlbumId     uint64     `json:"id"`
	AlbumName   string     `json:"name"`
	PicUrl      string     `json:"picUrl"`      // 专辑封面信息
	Description string     `json:"description"` // 专辑描述
	Size        int        `json:"size"`        // 专辑曲目数量
	Artist      ArtistInfo `json:"artist"`
}

type AlbumInfoWithSongList struct {
	ResourceState bool       `json:"resourceState"`
	Songs         []SongInfo `json:"songs"`
	Code          int        `json:"code"`
	Album         AlbumInfo  `json:"album"`
}
