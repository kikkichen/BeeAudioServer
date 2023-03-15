package controllers

import (
	model "BeeAudioServer/models"
	"BeeAudioServer/models/localmodel"
	"BeeAudioServer/repository"
	"BeeAudioServer/utils"
	"fmt"
	"log"
	"os"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/fishtailstudio/imgo"
)

type UploadFileController struct {
	beego.Controller
}

/* 上传动态博文随行图片 */
func (ufc *UploadFileController) UploadBlogImage() {
	uid, err := ufc.GetUint64("uid")
	filename := ufc.GetString("name")
	file, h, err := ufc.GetFile("uploadname")
	defer file.Close()
	if err != nil {
		log.Fatal("获取文件错误:", err)
	}

	/* 过滤无效字符 */
	filename = strings.Replace(strings.Replace(strings.Replace(filename, "{", "", -1), "}", "", -1), "\"", "", -1)

	/* Test */
	fmt.Println(fmt.Sprintf("uid: %v, name: %v", uid, filename))
	/* 判断用户路径是否存在， 不存在则创建该用户路径 */
	user_large_path := fmt.Sprintf("store/blog/%v/large", uid)
	user_thumbnail_path := fmt.Sprintf("store/blog/%v/thumbnail", uid)
	/* 判断缩略图路径与大图路径是否存在 */
	thumbnail_path_exist_signal, err := isPathExists(user_thumbnail_path)
	large_path_exist_signal, err := isPathExists(user_large_path)
	if !large_path_exist_signal {
		/* 用户博文动态图片 原图 路径不存在的情况，需要创建该目录路径 */
		err := os.MkdirAll(user_large_path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	if !thumbnail_path_exist_signal {
		/* 用户博文动态图片 缩略图 路径不存在的情况，需要创建该目录路径 */
		err := os.MkdirAll(user_thumbnail_path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	extend_name := GetFileExtendName(h.Filename)

	ufc.SaveToFile("uploadname", user_large_path+"/"+filename+"."+extend_name)
	/* 生成缩略图 */
	img := imgo.Load(user_large_path + "/" + filename + "." + extend_name)
	img_heigh := img.Height()
	img_width := img.Width()

	img.Thumbnail(img_width*1/5, img_heigh*1/5).Save(user_thumbnail_path + "/" + filename + ".png")

	client_response := model.ResponseBody{}
	client_response.OK = 1
	client_response.Message = "success"
	client_response.Data = filename
	ufc.Data["json"] = &client_response
	ufc.ServeJSON()
}

/* 上传自建歌单图片 */
func (ufc *UploadFileController) UploadPlayListCover() {
	uid, err := ufc.GetUint64("uid")
	pid, err := ufc.GetUint64("pid")
	file, h, err := ufc.GetFile("uploadfile")

	defer file.Close()
	if err != nil {
		log.Fatal("获取文件错误:", err)
	}
	/* 判断当前执行用户是否为歌单创建者 */
	var creator_id uint64
	err = utils.SqlDB.Model(&localmodel.UserPlayListModel{}).
		Select("uid").
		Where("pid = ?", pid).
		First(&creator_id).Error

	client_response := model.ResponseBody{}
	if creator_id == uid {
		/* 图片文件保存逻辑 */
		extend_name := GetFileExtendName(h.Filename)

		cover_path := fmt.Sprintf("store/playlist/cover/%v.%v", pid, extend_name)
		/* 保存/ 覆盖 歌单封面图片文件 */
		ufc.SaveToFile("uploadfile", cover_path)
		utils.SqlDB.Model(&localmodel.UserPlayListModel{}).
			Where("pid = ?", pid).
			Update("coverImgUrl", fmt.Sprintf("/playlist/cover/%v.%v", pid, extend_name))

		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = true
	} else {
		client_response.OK = 0
		client_response.Message = "error"
		client_response.Data = false
	}
	ufc.Data["json"] = &client_response
	ufc.ServeJSON()
}

/* 上传用户头像 */
func (ufc *UploadFileController) UploadUserAvatar() {
	uid, err := ufc.GetUint64("uid")
	file, h, err := ufc.GetFile("uploadfile")

	defer file.Close()
	if err != nil {
		log.Fatal("获取文件错误:", err)
	}

	client_response := model.ResponseBody{}
	var finish_signal bool = false

	extend_name := GetFileExtendName(h.Filename)
	random_code := utils.RangeRand(10, 99)
	avatar_path := fmt.Sprintf("store/user/avatar/%v%v.%v", uid, random_code, extend_name)
	/* 保存/ 覆盖 新的用户头像图片文件 */
	ufc.SaveToFile("uploadfile", avatar_path)
	/* 修改数据库中关于用户头像的路径信息 */
	finish_signal = repository.UpdateUserAvatarPath(utils.SqlDB, fmt.Sprintf("/user/avatar/%v%v.%v", uid, random_code, extend_name), uid)
	if finish_signal && err == nil {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = true
	} else {
		client_response.OK = 0
		client_response.Message = "error"
		client_response.Data = false
	}
	ufc.Data["json"] = &client_response
	ufc.ServeJSON()
}

/*	判断文件夹是否存在
 *	@param	path	目标查询文件夹目录
 */
func isPathExists(
	path string,
) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/** 依据文件名 获得文件扩展名
 *	@param	filename	文件名
 */
func GetFileExtendName(filename string) string {
	name_array := strings.Split(filename, ".")
	extend_name := name_array[len(name_array)-1]
	return extend_name
}
