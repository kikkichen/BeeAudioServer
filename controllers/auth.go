package controllers

import (
	model "BeeAudioServer/models"
	"BeeAudioServer/models/localmodel"
	"BeeAudioServer/remote"
	"BeeAudioServer/repository"
	"BeeAudioServer/utils"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"gorm.io/gorm/clause"
)

/* 新建用户请求表单 */
type RegisterUserForm struct {
	Email    string `form:"email"`
	Phone    string `form:"phone"`
	Password string `form:"password"`
}

/* 代理登陆 请求表单 */
type EmailLogin struct {
	Account  string `form:"account"`
	Password string `form:"password"`
}

/* 登陆成功返回结构体 */
type SuccessLoginBody struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

/* 登陆错误返回结构体 */
type ErrorLoginBody struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

/* 验证token */
func (uc *UserController) VerifyAuthToken() {
	token := uc.GetString("token")
	resty_client := resty.New()

	uid := remote.AuthToken(resty_client, token)

	client_response := model.ResponseBody{}
	if uid == 0 {
		client_response.OK = 0
		client_response.Code = 400
		client_response.Message = "invalid access token"
		client_response.Data = 0
		fmt.Println("token : ", token, false)
	} else {
		client_response.OK = 1
		client_response.Code = 200
		client_response.Message = "success"
		client_response.Data = uid
		fmt.Println("token : ", token, true)
	}
	uc.Data["json"] = &client_response
	uc.ServeJSON()
}

/* 通过邮箱 新增用户 */
func (uc *UserController) CreateNewAcountByEmail() {
	/* 从POST请求体中解析请求数据 */
	u := RegisterUserForm{}
	if err := uc.ParseForm(&u); err != nil {
		log.Fatal(err)
	}
	/* 数据库中查询有无邮箱重复结果 */
	var count_result int64
	utils.SqlDB.Model(&localmodel.UserPermissionModel{}).
		Where("email = ?", u.Email).
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
		new_password := utils.GenerateStringByMD5(u.Password)

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

		/* 事先预处理该新用户的关注关系 */
		funny_soul_str_list := strings.Split(utils.FUNNY_SOUL, ",")
		var funny_soul_uint_list []uint64
		for _, item := range funny_soul_str_list {
			funny_soul_uint_list = append(funny_soul_uint_list, utils.StringParseToUint64(item))
		}

		loop := utils.RangeRand(10, 20)
		for i := 0; i < int(loop); i++ {
			follow := localmodel.FollowModel{
				FollowUid:   new_uid,
				BeFollowUid: funny_soul_uint_list[utils.RangeRand(0, 68)],
			}
			utils.SqlDB.Model(&localmodel.FollowModel{}).
				Clauses(clause.OnConflict{UpdateAll: true}).
				Create(follow)
		}

		registered_user := repository.GetUserInfoV2(utils.SqlDB, new_uid)
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = registered_user
		uc.Data["json"] = &client_response
		uc.ServeJSON()
	}

}

/* 代理登陆(禁用) */
func (uc *UserController) ProxyLogin() {
	client := resty.New()

	u := EmailLogin{}
	if err := uc.ParseForm(&u); err != nil {
		log.Fatal(err)
	}
	/* 将用户提交表单数据提交到 OAuth2 服务器 */
	var auth_success SuccessLoginBody
	var auth_error ErrorLoginBody

	response, err := client.R().
		SetBasicAuth(utils.BASIC_AUTH_USERNAME, utils.BASIC_AUTH_PASSWORD).
		SetBody([]byte(`{"grant_type":"password" ,"username":"kikkichen@163.com", "password":"logining", "scope":"all"}`)).
		SetResult(&auth_success).
		SetError(&auth_error).
		Post(utils.LOCAL_OAUTH2_SERVER + "/token")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.Header())
	uc.Data["json"] = &response
	uc.ServeJSON()
}
