package controllers

import (
	"fmt"
	"io/ioutil"

	beego "github.com/beego/beego/v2/server/web"
)

type AccessFileController struct {
	beego.Controller
}

/* 访问缩略图 */
func (afc *AccessFileController) AccessBlogThumbnailPicture() {
	uid := afc.GetString(":uid")
	pic_name := afc.GetString(":picname")

	thumbnail_img_path := "store/blog/" + uid + "/thumbnail/" + pic_name + ".png"

	afc.Ctx.Output.Header("Content-Type", "image/png")
	afc.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", thumbnail_img_path))
	file, err := ioutil.ReadFile(thumbnail_img_path)
	if err != nil {
		fmt.Printf("[%v] 文件不存在！\n", thumbnail_img_path)
		afc.Ctx.Output.Header("Content-Type", "image/jpg")
		file, err = ioutil.ReadFile("store/nofound/404_pic.jpeg")
	}
	afc.Ctx.WriteString(string(file))

}

/* 访问原图 */
func (afc *AccessFileController) AccessBlogOriginalPicture() {
	uid := afc.GetString(":uid")
	pic_name := afc.GetString(":picname")

	large_img_path := "store/blog/" + uid + "/large/" + pic_name + ".jpg"

	afc.Ctx.Output.Header("Content-Type", "image/png")
	afc.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", large_img_path))
	file, err := ioutil.ReadFile(large_img_path)
	if err != nil {
		fmt.Printf("[%v] 文件不存在！\n", large_img_path)
		afc.Ctx.Output.Header("Content-Type", "image/jpg")
		file, err = ioutil.ReadFile("store/nofound/404_pic.jpeg")
	}
	afc.Ctx.WriteString(string(file))

}

/* 访问默认歌单封面图片 */
func (afc *AccessFileController) AccessDefaultCover() {
	cover_path := afc.GetString(":coverpath")

	large_img_path := "store/playlist/cover/" + cover_path

	afc.Ctx.Output.Header("Content-Type", "image/png")
	afc.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", large_img_path))
	file, err := ioutil.ReadFile(large_img_path)
	if err != nil {
		fmt.Printf("[%v] 文件不存在！\n", large_img_path)
		afc.Ctx.Output.Header("Content-Type", "image/jpg")
		file, err = ioutil.ReadFile("store/playlist/cover/default.jpg")
	}
	afc.Ctx.WriteString(string(file))
}

/* 访问用户头像 */
func (afc *AccessFileController) AccessUserAvatar() {
	cover_path := afc.GetString(":u_path")

	large_img_path := "store/user/avatar/" + cover_path

	afc.Ctx.Output.Header("Content-Type", "image/png")
	afc.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", large_img_path))
	file, err := ioutil.ReadFile(large_img_path)
	if err != nil {
		fmt.Printf("[%v] 文件不存在！\n", large_img_path)
		afc.Ctx.Output.Header("Content-Type", "image/png")
		file, err = ioutil.ReadFile("store/store/user/avatar/default.png")
	}
	afc.Ctx.WriteString(string(file))
}
