package controllers

import (
	model "BeeAudioServer/models"
	"BeeAudioServer/repository"
	"BeeAudioServer/utils"
	"log"
)

type CommentController struct {
	MainController
}

/* 查询目标博文的评论列表 (分页) */
func (cc *CommentController) GetTargetBlogComment() {
	blogId, err := cc.GetUint64("blog_id", 9900100001)
	page, err := cc.GetInt("page", 1)
	size, err := cc.GetInt("size", 20)

	result := repository.GetBlogComments(utils.SqlDB, blogId, page, size)

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
	cc.Data["json"] = &client_response
	cc.ServeJSON()
}

/* 评论一条文本博文 */
func (cc *CommentController) CommentTextBlog() {
	uid, err := cc.GetUint64("uid", 9900619251)
	text := cc.GetString("text")
	source := cc.GetString("source", "BeeAudio客户端")
	rootId, err := cc.GetUint64("root_id", 0)
	blogId, err := cc.GetUint64("bid", 0)

	err = repository.CommentBlog(utils.SqlDB, uid, blogId, rootId, text, source, false)

	if err != nil {
		log.Fatal(err)
	}

	client_response := model.ResponseBody{}
	if err == nil {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = "评论成功"
	} else {
		client_response.OK = 0
		client_response.Message = "Error"
		client_response.Data = "评论出错"
	}
	cc.Data["json"] = &client_response
	cc.ServeJSON()
}
