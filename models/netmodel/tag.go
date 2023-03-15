package netmodel

/*	分类Tag
 *
 */
type PlayListTag struct {
	Name     string `json:"name"`
	Category int    `json:"category"`
}

/* 请求外围Tag响应体 */
type ResponseHotTagsOutSideBody struct {
	Tags []PlayListTag `json:"tags"`
}

/* 请求全部Tag外围 */
type ResponseAllTagsOutSideBody struct {
	Code       int           `json:"code"`
	SubTagList []PlayListTag `json:"sub"`
}

/*	热门歌单列表Tag请求体示例

{
	"tags": [
		{
			"playlistTag": {
				"id": 5001,
				"name": "华语",
				"category": 0,
				"usedCount": 8539255,
				"type": 0,
				"position": 1,
				"createTime": 1378707544870,
				"highQuality": 1,
				"highQualityPos": 1,
				"officialPos": 1
			},
			"activity": false,
			"hot": true,
			"usedCount": 8539255,
			"position": 1,
			"category": 0,
			"createTime": 1378707544870,
			"name": "华语",
			"id": 5001,
			"type": 1
		},
		....
	}

*/
