package controllers

import (
	model "BeeAudioServer/models"
	"BeeAudioServer/models/localmodel"
	"BeeAudioServer/repository"
	"BeeAudioServer/utils"
	"fmt"
	"log"
)

type BlogController struct {
	MainController
}

/* 请求关注用户的发布博文 - 动态页 */
func (bc *BlogController) GetSubscribeBlog() {
	userId, err := bc.GetUint64("user_id", 9900100001)
	page, err := bc.GetInt("page", 1)
	size, err := bc.GetInt("size", 20)

	if err != nil {
		log.Fatal(err)
	}

	result := repository.GetMyFocusBlogs(utils.SqlDB, userId, page, size)
	var clientBlogs []localmodel.ClientBlog
	for _, item := range result {
		clientBlogs = append(clientBlogs, item.ToJsonClient(utils.SqlDB))
	}

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if len(clientBlogs) > 0 {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = clientBlogs
	} else {
		client_response.OK = 0
		client_response.Message = "You have not more Info"
		client_response.Data = clientBlogs
	}
	bc.Data["json"] = &client_response
	bc.ServeJSON()
}

/* 请求互粉用户关注的博文 - 动态页 */
func (bc *BlogController) GetFriendsBlog() {
	userId, err := bc.GetUint64("user_id", 9900100001)
	page, err := bc.GetInt("page", 1)
	size, err := bc.GetInt("size", 20)

	if err != nil {
		log.Fatal(err)
	}

	result := repository.GetMyFriendsBlogs(utils.SqlDB, userId, page, size)
	var clientBlogs []localmodel.ClientBlog
	for _, item := range result {
		clientBlogs = append(clientBlogs, item.ToJsonClient(utils.SqlDB))
	}

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if len(clientBlogs) > 0 {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = clientBlogs
	} else {
		client_response.OK = 2
		client_response.Message = "You maybe have not any fans"
		client_response.Data = clientBlogs
	}
	bc.Data["json"] = &client_response
	bc.ServeJSON()
}

/* 请求指定用户的博文 - 动态页 */
func (bc *BlogController) GetUserBlog() {
	userId, err := bc.GetUint64("user_id", 9900100001)
	page, err := bc.GetInt("page", 1)
	size, err := bc.GetInt("size", 20)

	if err != nil {
		log.Fatal(err)
	}

	/* 是否获取原创博文 */
	isOriginal, err := bc.GetBool("isOri", false)

	result := repository.GetTargetUserBlog(utils.SqlDB, userId, isOriginal, page, size)
	var clientBlogs []localmodel.ClientBlog
	for _, item := range result {
		clientBlogs = append(clientBlogs, item.ToJsonClient(utils.SqlDB))
	}

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if len(clientBlogs) > 0 {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = clientBlogs
	} else {
		client_response.OK = 2
		client_response.Message = "The User havent send Blog"
		client_response.Data = clientBlogs
	}
	bc.Data["json"] = &client_response
	bc.ServeJSON()
}

/* 根据博文ID, 查询博文内容 */
func (bc *BlogController) GetTargetBlogDetail() {
	blogId, err := bc.GetUint64("blog_id", 9900100001)

	if err != nil {
		log.Fatal(err)
	}

	result := repository.GetBlogDetail(utils.SqlDB, blogId)

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if result.Bid != 0 {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = result.ToJsonClient(utils.SqlDB)
	} else {
		client_response.OK = 0
		client_response.Message = "Not Found Blog"
		client_response.Data = nil
	}
	bc.Data["json"] = &client_response
	bc.ServeJSON()
}

/* 通过关键字搜索博文 */
func (bc *BlogController) SearchBlogByKeyWords() {
	keyword := bc.GetString("keyword")
	page, err := bc.GetInt("page", 1)
	size, err := bc.GetInt("size", 20)

	if err != nil {
		log.Fatal(err)
	}

	result := repository.SearchBlogByKeyWord(utils.SqlDB, keyword, page, size)
	var clientBlogs []localmodel.ClientBlog
	for _, item := range result {
		clientBlogs = append(clientBlogs, item.ToJsonClient(utils.SqlDB))
	}

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if len(clientBlogs) > 0 {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = clientBlogs
	} else {
		client_response.OK = 0
		client_response.Message = "You have not more Info"
		client_response.Data = clientBlogs
	}
	bc.Data["json"] = &client_response
	bc.ServeJSON()
}

/* 发送一条文本博文动态 */
func (bc *BlogController) SendTextBlog() {
	uid, err := bc.GetUint64("uid")
	text := bc.GetString("text")
	source := bc.GetString("source", "BeeAudio客户端")
	pic_urls := bc.GetStrings("pic_urls")
	media_data := bc.GetString("media_data")

	var pics []localmodel.PicUrl
	if len(pic_urls) != 0 {
		for _, item := range pic_urls {
			pics = append(pics, localmodel.PicUrl{
				Url: item,
			})
		}
		/* 有携带图片的博文动态数据 */
		err = repository.SendBlog(utils.SqlDB, uid, text, source, pics, media_data)
	} else {
		/* 文本博文数据 */
		err = repository.SendBlog(utils.SqlDB, uid, text, source, []localmodel.PicUrl{}, media_data)
	}

	client_response := model.ResponseBody{}
	if err == nil {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = "博文动态发送成功"
	} else {
		client_response.OK = 0
		client_response.Message = "Error"
		client_response.Data = "动态发送出了一些小问题"
	}
	bc.Data["json"] = &client_response
	bc.ServeJSON()
}

/* 发送一条带图片的博文动态 */
func (bc *BlogController) SendImageBlog() {
	uid, err := bc.GetUint64("uid")
	text := bc.GetString("text")
	source := bc.GetString("source", "BeeAudio客户端")
	pic_urls := bc.GetStrings("pic_urls")

	var pics []localmodel.PicUrl
	if len(pic_urls) != 0 {
		for _, item := range pic_urls {
			pics = append(pics, localmodel.PicUrl{
				Url: item,
			})
		}
	}

	fmt.Printf("uid : %v, text: %v, source: %v, pic_url: %v\n", uid, text, source, pics)
	// err = repository.SendBlog(db, uid, text, source, []localmodel.PicUrl{}, "")
	client_response := model.ResponseBody{}
	if err == nil {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = "博文动态发送成功"
	} else {
		client_response.OK = 0
		client_response.Message = "Error"
		client_response.Data = "动态发送出了一些小问题"
	}
	bc.Data["json"] = &client_response
	bc.ServeJSON()
}

/* 转发一条博文动态 */
func (bc *BlogController) RetweetedTextBlog() {
	uid, err := bc.GetUint64("uid", 9900619251)
	text := bc.GetString("text")
	source := bc.GetString("source", "BeeAudio客户端")
	retweetedId, err := bc.GetUint64("retweeted_id", 0)

	if err != nil {
		log.Fatal(err)
	}

	err = repository.RepostBlog(utils.SqlDB, uid, text, source, retweetedId)
	client_response := model.ResponseBody{}
	if err == nil {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = "转发发送成功"
	} else {
		client_response.OK = 0
		client_response.Message = "Error"
		client_response.Data = "转发有误"
	}
	bc.Data["json"] = &client_response
	bc.ServeJSON()
}

/**	管理员 浏览全部博文 - 分页
 *
 */
func (bc *BlogController) BrowserAllBlogByAdminPage() {
	page, err := bc.GetInt("page", 1)
	size, err := bc.GetInt("size", 20)

	list := repository.GetAllBlogsByAdminPage(utils.SqlDB, page, size)

	var clientBlogs []localmodel.ClientBlog
	for _, item := range list {
		clientBlogs = append(clientBlogs, item.ToJsonClient(utils.SqlDB))
	}

	client_response := model.ResponseBody{}
	if err == nil {
		if len(list) != 0 {
			client_response.OK = 1
			client_response.Message = "success"
			client_response.Data = clientBlogs
		} else {
			client_response.OK = 0
			client_response.Message = "empty"
			client_response.Data = []repository.MyFocusBlog{}
		}
	} else {
		bc.Ctx.ResponseWriter.WriteHeader(500)
		client_response.OK = 0
		client_response.Message = "server error"
		client_response.Data = []repository.MyFocusBlog{}
	}
	bc.Data["json"] = &client_response
	bc.ServeJSON()
}

/**	管理员 覆盖删除博内容
 *	将原本的博文动态正文替换为 [该博文由于违反社区规定已被删除] ，并且移除 图片链接 和 音频分享功能
 */
func (bc *BlogController) DeleteTargetBlog() {
	bid, err := bc.GetUint64("bid")
	err = repository.DeleteTargetBlog(utils.SqlDB, bid)
	client_response := model.ResponseBody{}
	if err == nil {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = true
	} else {
		bc.Ctx.ResponseWriter.WriteHeader(500)
		client_response.OK = 0
		client_response.Message = "server error"
		client_response.Data = false
	}
	bc.Data["json"] = &client_response
	bc.ServeJSON()
}

/**
 *
 */
