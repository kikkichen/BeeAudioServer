package controllers

import (
	model "BeeAudioServer/models"
	"BeeAudioServer/repository"
	"BeeAudioServer/utils"
)

type SubscribeDataController struct {
	MainController
}

/* 获取我的音频项目订阅信息 */
func (sdc *SubscribeDataController) GetMySubscribeData() {
	var err error
	uid, err := sdc.GetUint64("uid")
	/* 查询我的订阅信息数据 */
	result, err := repository.GetSubscribeData(utils.SqlDB, uid)

	client_response := model.ResponseBody{}

	if err != nil {
		client_response.OK = 0
		client_response.Code = 200
		client_response.Message = "error"
		client_response.Data = nil
	} else {
		client_response.OK = 1
		client_response.Code = 200
		client_response.Message = "success"
		client_response.Data = result
	}
	sdc.Data["json"] = &client_response
	sdc.ServeJSON()
}

/* 同步我的音频订阅项目数据 */
func (sdc *SubscribeDataController) SyncMySubscribeData() {
	var err error
	uid, err := sdc.GetUint64("uid")
	subscribe_data := sdc.GetString("data")

	client_response := model.ResponseBody{}

	if len(subscribe_data) == 0 {
		/* 无效的同步请求 */
		client_response.OK = 0
		client_response.Code = 200
		client_response.Message = "Unvalid Request"
		client_response.Data = false
	} else {
		/* 同步数据 */
		err = repository.SyncSubscribeData(utils.SqlDB, uid, subscribe_data)
		if err != nil {
			client_response.OK = 0
			client_response.Code = 200
			client_response.Message = "error"
			client_response.Data = false
		} else {
			client_response.OK = 1
			client_response.Code = 200
			client_response.Message = "success"
			client_response.Data = true
		}
	}
	sdc.Data["json"] = &client_response
	sdc.ServeJSON()
}
