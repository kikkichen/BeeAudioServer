package controllers

import (
	model "BeeAudioServer/models"
	"BeeAudioServer/remote"

	"github.com/go-resty/resty/v2"
)

type SearchController struct {
	MainController
}

/* 搜索结果请求 */
func (sc *SearchController) GetSearchResult() {
	/* 请求参数获取 */
	var keywords string = sc.GetString("keywords")
	search_type, _ := sc.GetInt("type", 1)
	page, _ := sc.GetInt("page", 1)
	size, _ := sc.GetInt("size", 30)

	/* 需要获取的歌单ID */
	client := resty.New()
	search_result := remote.RequestSearchResult(client, search_type, keywords, page, size)

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	client_response.OK = 1
	client_response.Message = "success"
	client_response.Data = search_result
	sc.Data["json"] = &client_response
	sc.ServeJSON()
}
