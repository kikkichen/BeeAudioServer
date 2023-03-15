package controllers

import (
	model "BeeAudioServer/models"
	"BeeAudioServer/repository"

	"github.com/go-resty/resty/v2"
)

type HistoryController struct {
	MainController
}

/**	查看历史记录
 *
 */
func (hc *HistoryController) AccessMyHistoryPlayRecord() {
	uid, err := hc.GetUint64("uid")
	page, err := hc.GetInt("page", 1)
	size, err := hc.GetInt("size", 20)

	client := resty.New()

	song_info_list, song_history_info_list, err := repository.BrowserPlayHistory(client, uid, page, size)

	/* 集成元数据 */
	var data_list []map[string]interface{}
	for index, song_info := range song_info_list {
		data_list = append(data_list, map[string]interface{}{
			"song":    song_info,
			"play_at": song_history_info_list[index].PlayAt,
		})
	}

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if err == nil {
		client_response.OK = 1
		client_response.Code = 200
		client_response.Message = "success"
		client_response.Data = data_list
	} else {
		client_response.OK = 0
		client_response.Code = 500
		client_response.Message = "Browser Error"
		client_response.Data = nil
	}
	hc.Data["json"] = &client_response
	hc.ServeJSON()
}

/**	添加 / 更新 历史记录
 *
 */
func (hc *HistoryController) AddMyHistoryPlayItem() {
	uid, err := hc.GetUint64("uid")
	sid, err := hc.GetUint64("sid")
	client := resty.New()

	err = repository.AddPlayHistory(client, uid, sid)
	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if err == nil {
		client_response.OK = 1
		client_response.Code = 200
		client_response.Message = "add success"
		client_response.Data = true
	} else {
		client_response.OK = 0
		client_response.Code = 500
		client_response.Message = "failed"
		client_response.Data = false
	}
	hc.Data["json"] = &client_response
	hc.ServeJSON()
}

/**	清空历史记录
 *
 */
func (hc *HistoryController) ClearMyHistoryData() {
	uid, err := hc.GetUint64("uid")
	err = repository.ClearPlayHistory(uid)

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if err == nil {
		client_response.OK = 1
		client_response.Code = 200
		client_response.Message = "clear success"
		client_response.Data = true
	} else {
		client_response.OK = 0
		client_response.Code = 500
		client_response.Message = "failed"
		client_response.Data = false
	}
	hc.Data["json"] = &client_response
	hc.ServeJSON()
}
