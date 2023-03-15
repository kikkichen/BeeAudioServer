package netmodel

/* 搜索类型：分为：1：单曲，10：专辑，100：艺人，1000：歌单 */

/* 搜索结果类型 */
type SearchResult interface {
	SingleSongsResult | AblumResult | ArtitstResult | PlayListResult
}

/* 搜索结果外围 - 针对客户端 */
type ResponseSearchBody[T SearchResult] struct {
	Code   int `json:"200"`
	Result T   `json:"result"`
}

/* 单曲搜索结果 */
type SingleSongsResult struct {
	Songs     []SongInfo `json:"songs"`
	SongCount int        `json:"songCount"`
}

/* 专辑搜索结果 */
type AblumResult struct {
	Albums     []AlbumInfo `json:"albums"`
	AlbumCount int         `json:"albumCount"`
}

/* 艺人搜索结果 */
type ArtitstResult struct {
	Artists     []ArtistInfo `json:"artists"`
	ArtistCount int          `json:"artistCount"`
}

/* 歌单搜索结果 */
type PlayListResult struct {
	PlayLists     []PlayListInfo `json:"playlists"`
	PlaylistCount int            `json:"playlistCount"`
}
