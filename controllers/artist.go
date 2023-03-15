package controllers

import (
	model "BeeAudioServer/models"
	"BeeAudioServer/remote"
	"BeeAudioServer/utils"

	"github.com/go-resty/resty/v2"
)

type ArtistController struct {
	MainController
}

/* 获取艺人详细信息 */
func (ac *ArtistController) GetArtistDetail() {
	/* 需要获取的艺人ID */
	var artist_id string = ac.GetString("artist_id")
	client := resty.New()
	artist_detail := remote.RequestArtistDetail(client, artist_id)
	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if artist_detail.ArtistId == 0 {
		/* 若没有获取到有效的艺人信息， 返回错误处理 */
		client_response.OK = 0
		client_response.Message = "No Exist Artist"
		client_response.Data = nil
	} else {
		/* 若请求数据有效，返回有效的艺人详细信息数据 */
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = artist_detail
	}
	ac.Data["json"] = &client_response
	ac.ServeJSON()
}

/* 获取艺人的曲目信息列表 */
func (ac *ArtistController) GetArtistSongs() {
	/* 需要获取的艺人ID */
	var artist_id string = ac.GetString("artist_id")
	playlist_page, _ := ac.GetInt("page", 1)
	playlist_size, _ := ac.GetInt("size", 50)
	client := resty.New()
	artist_songs := remote.RequestArtistSongs(client, artist_id, playlist_page, playlist_size)
	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if len(artist_songs) == 0 {
		/* 若没有获取到有效的艺人信息， 返回错误处理 */
		client_response.OK = 0
		client_response.Message = "No Exist Songs"
		client_response.Data = nil
	} else {
		/* 若请求数据有效，返回有效的艺人详细信息数据 */
		/* 歌曲信息过滤 */
		artist_songs = remote.FilterSongs(utils.SqlDB, artist_songs)

		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = artist_songs
	}
	ac.Data["json"] = &client_response
	ac.ServeJSON()
}

/* 获取艺人的专辑信息列表 */
func (ac *ArtistController) GetArtistAlbums() {
	/* 需要获取的艺人ID */
	var artist_id string = ac.GetString("artist_id")
	playlist_page, _ := ac.GetInt("page", 1)
	playlist_size, _ := ac.GetInt("size", 30)
	client := resty.New()
	artist_albums := remote.RequestArtistAlbums(client, artist_id, playlist_page, playlist_size)
	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if len(artist_albums) == 0 {
		/* 若没有获取到有效的艺人信息， 返回错误处理 */
		client_response.OK = 0
		client_response.Message = "No Exist Songs"
		client_response.Data = nil
	} else {
		/* 若请求数据有效，返回有效的艺人详细信息数据 */
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = artist_albums
	}
	ac.Data["json"] = &client_response
	ac.ServeJSON()
}
