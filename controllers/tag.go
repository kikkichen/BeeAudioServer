package controllers

import (
	model "BeeAudioServer/models"
	"BeeAudioServer/models/netmodel"
	"BeeAudioServer/remote"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/go-resty/resty/v2"
)

type TagController struct {
	MainController
}

type TagsForClient struct {
	HotTag []netmodel.PlayListTag `json:"hot"`
	AllTag []netmodel.PlayListTag `json:"all"`
}

/* 请求Tag */
func (tc *TagController) GetTagListInfo() {
	client := resty.New()
	hot_tags := remote.RequestPlayListHotTag(client)
	all_tags := remote.RequestPlayListAllTag(client)
	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if (len(hot_tags) == 0) || (len(all_tags) == 0) {
		client_response.OK = 0
		client_response.Message = "No PlayList Tag Return"
		client_response.Data = nil
	} else {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = TagsForClient{
			HotTag: hot_tags,
			AllTag: all_tags,
		}
	}
	tc.Data["json"] = &client_response
	tc.ServeJSON()
}

/* 获取对应热门Tag/Cat图片封面 */
func (tc *TagController) GetHotTagCoverImage() {
	/*  标签封面图片在工程中的路径 */
	tag_cover_path := "./store/tag"
	cat := tc.GetString("cat")

	client := resty.New()
	hot_tags := remote.RequestPlayListHotTag(client)

	/* 判断请求的目的tag/Cat是否存在于热门Tag集合中 */
	loop_count := 0
	for _, item := range hot_tags {
		if cat == item.Name {
			break
		}
		loop_count += 1
	}

	/* 若请求Tag不存在于热门Tag中 */
	if loop_count >= 10 {
		img := path.Join(tag_cover_path, "/404_pic.jpeg")
		tc.Ctx.Output.Header("Content-Type", "image/jpeg")
		tc.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", img))
		file, err := ioutil.ReadFile(img)
		if err != nil {
			fmt.Printf("🔴 未找到 [404] Tag/Cat封面图片\n\n[报错]%v", err)
			return
		}
		tc.Ctx.WriteString(string(file))
	} else {
		img := path.Join(tag_cover_path, "/"+strings.Replace(cat, "/", "", -1)+".jpg")
		tc.Ctx.Output.Header("Content-Type", "image/jpg")
		tc.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", img))
		file, err := ioutil.ReadFile(img)
		if err != nil {
			fmt.Printf("🔴 未找到指定Tag/Cat封面图片\n\n[报错]%v", err)
			return
		}
		tc.Ctx.WriteString(string(file))
	}

}
