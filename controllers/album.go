package controllers

import (
	model "BeeAudioServer/models"
	"BeeAudioServer/remote"
	"BeeAudioServer/utils"

	"github.com/go-resty/resty/v2"
)

type AlbumController struct {
	MainController
}

/* 获取专辑详细信息 */
func (ac *AlbumController) GetAlbumDetail() {

	/* 需要获取的艺人ID */
	var album_id string = ac.GetString("album_id")
	client := resty.New()
	album_detail := remote.RequestAlbumDetail(client, album_id)

	if len(album_detail.Songs) != 0 {
		album_detail.Songs = remote.FilterSongs(utils.SqlDB, album_detail.Songs)
	}

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if !album_detail.ResourceState {
		/* 若没有获取到有效的专辑信息， 返回错误处理 */
		client_response.OK = 0
		client_response.Message = "No Exist Album"
		client_response.Data = nil
	} else {
		/* 若请求数据有效，返回有效的专辑详细信息数据 */
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = album_detail
	}
	ac.Data["json"] = &client_response
	ac.ServeJSON()
}
