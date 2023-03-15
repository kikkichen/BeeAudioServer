package controllers

import (
	model "BeeAudioServer/models"
	"BeeAudioServer/repository"
	"BeeAudioServer/utils"
	"log"
)

/* 查询目标博文的转发列表 */
func (bc *BlogController) GetBlogReportList() {

	blogId, err := bc.GetUint64("blog_id", 9900100001)
	page, err := bc.GetInt("page", 1)
	size, err := bc.GetInt("size", 20)
	isReport, err := bc.GetBool("is_report", false)

	var result []repository.BlogReportItem
	if isReport {
		result = repository.GetRepostBlogReposts(utils.SqlDB, blogId, page, size)
	} else {
		result = repository.GetBlogReposts(utils.SqlDB, blogId, page, size)
	}

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
		client_response.Message = "No Report Info"
		client_response.Data = result
	}
	bc.Data["json"] = &client_response
	bc.ServeJSON()
}
