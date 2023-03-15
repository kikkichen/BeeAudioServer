package controllers

import (
	model "BeeAudioServer/models"
	"BeeAudioServer/models/localmodel"
	"BeeAudioServer/repository"
	"BeeAudioServer/utils"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm/clause"
)

type UserController struct {
	MainController
}

/**	查询目标用户简易信息
 *
 */
func (uc *UserController) GetSimpleUserInfo() {

	userId, err := uc.GetUint64("user_id", 9900100001)
	resultUser := repository.GetUserInfo(utils.SqlDB, userId)

	if err != nil {
		log.Fatal(err)
	}

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if resultUser.Uid == 0 {
		client_response.OK = 0
		client_response.Message = "No Found User"
		client_response.Data = nil
	} else {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = resultUser
	}
	uc.Data["json"] = &client_response
	uc.ServeJSON()
}

/**	获取 携带用户类型的用户信息
 *
 */
func (uc *UserController) GetUserInfoV2() {

	userId, err := uc.GetUint64("user_id", 9900100001)
	resultUser := repository.GetUserInfoV2(utils.SqlDB, userId)

	if err != nil {
		log.Fatal(err)
	}

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if resultUser.Uid == 0 {
		client_response.OK = 0
		client_response.Message = "No Found User"
		client_response.Data = nil
	} else {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = resultUser
	}
	uc.Data["json"] = &client_response
	uc.ServeJSON()
}

/* 获取用户的详细信息 (密码脱敏) */
func (uc *UserController) GetUserDetail() {

	userId, err := uc.GetUint64("user_id", 9900100001)
	resultUser := repository.GetUserInfoWithoutPasword(utils.SqlDB, userId)

	if err != nil {
		log.Fatal(err)
	}

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if resultUser.Uid == 0 {
		client_response.OK = 0
		client_response.Message = "No Found User"
		client_response.Data = nil
	} else {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = resultUser
	}
	uc.Data["json"] = &client_response
	uc.ServeJSON()
}

/* 获取用户的关注数、粉丝数和互粉数目 */
func (uc *UserController) GetUserFollowAndFansCount() {

	userId, err := uc.GetUint64("user_id", 9900100001)

	follow, fans, friends := repository.GetCountOfFollowAndFans(utils.SqlDB, userId)
	resultMap := map[string]interface{}{
		"follows": follow,
		"fans":    fans,
		"friends": friends,
	}

	if err != nil {
		log.Fatal(err)
	}

	client_response := model.ResponseBody{}
	client_response.OK = 1
	client_response.Message = "success"
	client_response.Data = resultMap
	uc.Data["json"] = &client_response
	uc.ServeJSON()
}

/* 通过邮箱查找ID */
func (uc *UserController) FindUserIdByEmail() {

	email := uc.GetString("email")

	var uid uint64
	utils.SqlDB.Model(&localmodel.UserModel{}).
		Select("user_table.uid").
		Joins("LEFT JOIN mblog.user_permission_table ON user_table.uid = user_permission_table.uid").
		Where("user_permission_table.email = ?", email).
		First(&uid)
	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if uid == 0 {
		client_response.OK = 0
		client_response.Message = "No Found User"
		client_response.Data = nil
	} else {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = uid
	}
	uc.Data["json"] = &client_response
	uc.ServeJSON()

}

/* 关键词搜索用户 */
func (uc *UserController) SearchUserByKeyWords() {

	keyword := uc.GetString("keyword")
	page, err := uc.GetInt("page", 1)
	size, err := uc.GetInt("size", 20)

	result_user_list := repository.SearchUserByKeyWord(utils.SqlDB, keyword, page, size)

	if err != nil {
		log.Fatal(err)
	}

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if len(result_user_list) == 0 {
		client_response.OK = 0
		client_response.Message = "No Found Any User"
		client_response.Data = nil
	} else {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = result_user_list
	}
	uc.Data["json"] = &client_response
	uc.ServeJSON()
}

/* 通过邮箱查询有无该用户 */
func (uc *UserController) CheckExistAcount() {
	/* 从参数中获取查询ID */
	userId, err := uc.GetUint64("user_id", 9900100001)
	var findUser localmodel.UserModel
	utils.SqlDB.Model(&localmodel.UserModel{}).
		Where("uid = ?", userId).
		First(&findUser)

	if err != nil {
		log.Fatal(err)
	}
}

/* 通过手机号码 新增用户 */
func (uc *UserController) CreateNewAcountByPhone() {
	/* 从POST请求体中解析请求数据 */
	u := RegisterUserForm{}
	if err := uc.ParseForm(&u); err != nil {
		log.Fatal(err)
	}

	/* 数据库中查询有无邮箱重复结果 */
	var count_result int64
	utils.SqlDB.Model(&localmodel.UserPermissionModel{}).
		Where("phone = ?", u.Phone).
		Count(&count_result)

	client_response := model.ResponseBody{}

	/* 若存在邮箱记录， 反馈给客户端，表示该邮箱已注册 */
	if count_result > 0 {
		client_response.OK = 0
		client_response.Message = "Email have been register"
		client_response.Data = nil
		uc.Data["json"] = &client_response
		uc.ServeJSON()
	} else {
		new_password := utils.GenerateStringByMD5(u.Email)

		new_uid := repository.AddUser(utils.SqlDB, new_password, u.Email, "")
		/* 新建playlist */
		temp_pid := utils.StringParseToUint64("6900" + utils.RandomNumberString(8, utils.DefaultNumber))
		temp_user_playlist := localmodel.UserPlayListModel{
			Uid:         new_uid,
			Pid:         temp_pid,
			Cover:       fmt.Sprintf("/playlist/cover/default.jpg"),
			Name:        fmt.Sprintf("用户%v喜爱的歌曲", new_uid),
			CreateTime:  time.Now(),
			Description: "ta的歌单，什么都没有写。",
			Tags:        "[]",
			Public:      1,
		}
		utils.SqlDB.Model(&localmodel.UserPlayListModel{}).
			Clauses(clause.OnConflict{UpdateAll: true}).
			Create(&temp_user_playlist)

		/* 将该歌单绑定用户uid */
		utils.SqlDB.Model(&localmodel.UserFavoritePlaylist{}).
			Clauses(clause.OnConflict{UpdateAll: true}).
			Create(&localmodel.UserFavoritePlaylist{
				Uid: new_uid,
				Pid: temp_pid,
			})

		registered_user := repository.GetUserInfoV2(utils.SqlDB, new_uid)
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = registered_user
		uc.Data["json"] = &client_response
		uc.ServeJSON()
	}

}

/**	修改用户密码
 *
 */
func (uc *UserController) ModifierUserPassword() {
	var err error
	/* 用户ID */
	uid, err := uc.GetUint64("uid")
	/* 原密码 */
	origialPassword := uc.GetString("original_password")
	/* 新密码 */
	newPassword := uc.GetString("new_password")

	client_response := model.ResponseBody{}

	var origialEncryptionPassword string
	err = utils.SqlDB.Model(&localmodel.UserPermissionModel{}).
		Select("password").
		Where("uid = ?", uid).
		First(&origialEncryptionPassword).Error
	if origialEncryptionPassword == utils.GenerateStringByMD5(origialPassword) {
		/* 原密码核对正确， 执行密码修改逻辑 */
		result := repository.UpdateUserPassword(utils.SqlDB, newPassword, uid)
		if result && err == nil {
			client_response.OK = 1
			client_response.Code = 200
			client_response.Message = "success"
			client_response.Data = true
		} else {
			client_response.OK = 1
			client_response.Code = 500
			client_response.Message = "error"
			client_response.Data = false
		}
	} else {
		/* 原密码核对错误， 无更改权限  */
		client_response.OK = 0
		client_response.Code = 403
		client_response.Message = "forbid control"
		client_response.Data = false
	}
	uc.Data["json"] = &client_response
	uc.ServeJSON()
}

/**	通过邮箱账户 修改用户密码
 *
 */
func (uc *UserController) ModifierUserPasswordByEmailAccount() {
	var err error
	/* 用户Email */
	email := uc.GetString("email")
	/* 原密码 */
	origialPassword := uc.GetString("original_password")
	/* 新密码 */
	newPassword := uc.GetString("new_password")

	client_response := model.ResponseBody{}

	/* 通过邮箱查找到账户ID */
	var userId uint64
	err = utils.SqlDB.Model(&localmodel.UserPermissionModel{}).
		Select("uid").
		Where("email = ?", email).
		First(&userId).Error

	var origialEncryptionPassword string
	err = utils.SqlDB.Model(&localmodel.UserPermissionModel{}).
		Select("password").
		Where("uid = ?", userId).
		First(&origialEncryptionPassword).Error
	if origialEncryptionPassword == utils.GenerateStringByMD5(origialPassword) {
		/* 原密码核对正确， 执行密码修改逻辑 */
		result := repository.UpdateUserPassword(utils.SqlDB, newPassword, userId)
		if result && err == nil {
			client_response.OK = 1
			client_response.Code = 200
			client_response.Message = "success"
			client_response.Data = true
		} else {
			client_response.OK = 1
			client_response.Code = 500
			client_response.Message = "error"
			client_response.Data = false
		}
	} else {
		/* 原密码核对错误， 无更改权限  */
		client_response.OK = 0
		client_response.Code = 403
		client_response.Message = "forbid control"
		client_response.Data = false
	}
	uc.Data["json"] = &client_response
	uc.ServeJSON()
}

/**	修改用户自身非敏感信息
 *
 */
func (uc *UserController) ModifierMyInfo() {
	var err error
	/* 用户ID */
	uid, err := uc.GetUint64("uid")
	/* 用户昵称 */
	nickname := uc.GetString("name")
	/* 用户简介 */
	description := uc.GetString("description")
	/* email */
	email := uc.GetString("email")
	/* phone */
	phone_number := uc.GetString("phone")
	/* birthday */
	birthday, err := uc.GetUint64("birthday", 0)

	/* 用户原信息 */
	original_user_info := repository.GetUserInfoWithoutPasword(utils.SqlDB, uid)

	/* 正常完成信号 */
	var finish_signal bool = false
	if nickname != original_user_info.Name && len(nickname) != 0 {
		finish_signal = repository.UpdateUserName(utils.SqlDB, nickname, uid)
	}
	if description != original_user_info.Description && len(description) != 0 {
		finish_signal = repository.UpdateUserDescription(utils.SqlDB, description, uid)
	}
	if email != original_user_info.Email && len(email) != 0 {
		finish_signal = repository.UpdateUserEmail(utils.SqlDB, email, uid)
	}
	if phone_number != original_user_info.Phone && len(phone_number) != 0 {
		finish_signal = repository.UpdateUserPhone(utils.SqlDB, phone_number, uid)
	}
	if birthday != uint64(original_user_info.Birthday.Unix()) && birthday != 0 {
		finish_signal = repository.UpdateUserBirthday(utils.SqlDB, utils.UnixMilliToTime(int64(birthday)), uid)
	}

	client_response := model.ResponseBody{}
	if finish_signal && err == nil {
		/* 修改正常情况 */
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = true
	} else {
		/* 修改错误情况 */
		client_response.OK = 0
		client_response.Message = "error"
		client_response.Data = false
	}
	uc.Data["json"] = &client_response
	uc.ServeJSON()
}

/**	浏览用户目录
 *
 */
func (uc *UserController) BrowserAllUser() {
	page, err := uc.GetInt("page")
	size, err := uc.GetInt("size")
	sort, err := uc.GetBool("sort", false)

	/* 用户列表 - 分页 */
	user_list, err := repository.BrowserAllUser(utils.SqlDB, page, size, sort)
	/* 用户量总数 */
	total, err := repository.GetUserTotalNumber(utils.SqlDB)

	client_response := model.ResponseBody{}
	if err == nil {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = map[string]interface{}{
			"users": user_list,
			"total": total,
		}
	} else {
		/* 修改错误情况 */
		uc.Ctx.ResponseWriter.WriteHeader(500)
		client_response.OK = 0
		client_response.Message = "error"
		client_response.Data = nil
	}
	uc.Data["json"] = &client_response
	uc.ServeJSON()
}

/**	管理员通过关键字搜索相关用户
 *
 */
func (uc *UserController) SearchUserByAminKeyword() {
	keyword := uc.GetString("keyword")
	page, err := uc.GetInt("page", 1)
	size, err := uc.GetInt("size", 20)
	sort, err := uc.GetBool("sort", false)

	/* 用户列表 - 分页 */
	user_list, err := repository.SearchUserByKeyword(utils.SqlDB, keyword, page, size, sort)

	client_response := model.ResponseBody{}
	if err == nil {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = user_list
	} else {
		uc.Ctx.ResponseWriter.WriteHeader(500)
		client_response.OK = 0
		client_response.Message = "error"
		client_response.Data = nil
	}
	uc.Data["json"] = &client_response
	uc.ServeJSON()
}

/**	管理员修改用户信息
 *
 */
func (uc *UserController) ModifierUserDetailByAdmin() {
	uid, err := uc.GetUint64("uid")
	nickname := uc.GetString("name")
	description := uc.GetString("description")
	email := uc.GetString("email")
	phone := uc.GetString("phone")

	result_user, err := repository.ModifierUserDetailByAdmin(utils.SqlDB, uid, nickname, description, email, phone)

	client_response := model.ResponseBody{}
	if err == nil {
		/* 修改正常情况 */
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = result_user
	} else {
		/* 修改错误情况 */
		uc.Ctx.ResponseWriter.WriteHeader(500)
		client_response.OK = 0
		client_response.Message = "error"
		client_response.Data = nil
	}
	uc.Data["json"] = &client_response
	uc.ServeJSON()
}

/**	管理员注销用户
 *
 */
func (uc *UserController) LogoutUserByAdmin() {
	uid, err := uc.GetUint64("uid")

	err = repository.LogoutUserByAdmin(utils.SqlDB, uid)

	client_response := model.ResponseBody{}
	if err == nil {
		/* 修改正常情况 */
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = true
	} else {
		/* 修改错误情况 */
		uc.Ctx.ResponseWriter.WriteHeader(500)
		client_response.OK = 0
		client_response.Message = "error"
		client_response.Data = false
	}
	uc.Data["json"] = &client_response
	uc.ServeJSON()
}
