package controllers

import (
	model "BeeAudioServer/models"
	"BeeAudioServer/repository"
	"BeeAudioServer/utils"
	"log"
)

type AttitudeConttoller struct {
	MainController
}

/* 查询目标用户博文的点赞列表 (分页) */
func (ac *AttitudeConttoller) GettargetBlogAttitudes() {

	blogId, err := ac.GetUint64("blog_id", 9900100001)
	page, err := ac.GetInt("page", 1)

	size, err := ac.GetInt("size", 20)

	result := repository.GetBlogAttitudes(utils.SqlDB, blogId, page, size)

	if err != nil {
		log.Fatal(err)
	}
	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if len(result) > 0 {
		client_response.OK = 1
		client_response.Message = "success"

		client_response.Data = result
	} else {
		client_response.OK = 1
		client_response.Message = "have not any comment"
		client_response.Data = result
	}
	ac.Data["json"] = &client_response

	ac.ServeJSON()
}

/* 查询用户是否对目标博文动态有过点赞记录 */
func (ac *AttitudeConttoller) IsAttitudedRecordExist() {

	uid, err := ac.GetUint64("uid")
	bid, err := ac.GetUint64("bid")

	result := repository.IsPullAttitude(utils.SqlDB, bid, uid)

	if err != nil {
		log.Fatal(err)
	}

	client_response := model.ResponseBody{}
	if err == nil {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = result
	} else {
		client_response.OK = 0
		client_response.Message = "Error"
		client_response.Data = result
	}
	ac.Data["json"] = &client_response

	ac.ServeJSON()
}

/* 为目标博文点赞 */
func (ac *AttitudeConttoller) AttitudeTargetBlog() {

	uid, err := ac.GetUint64("uid", 9900619251)
	source := ac.GetString("source", "BeeAudio客户端")
	bid, err := ac.GetUint64("bid", 0)

	result, err := repository.ActionAttitude(utils.SqlDB, bid, uid, source)

	if err != nil {
		log.Fatal(err)
	}
	client_response := model.ResponseBody{}
	if err == nil {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = result
	} else {
		client_response.OK = 0
		client_response.Message = "Error"
		client_response.Data = result
	}
	ac.Data["json"] = &client_response

	ac.ServeJSON()
}
