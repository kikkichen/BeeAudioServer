package netmodel

/* 关于首页热门的请求集合 */

/*	热门歌单请求 外围1请求体
 *
 */
type PlayListCollection struct {
	List []PlayListInfo `json:"playlists"`
}
