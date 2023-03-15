package controllers

import (
	model "BeeAudioServer/models"
	"BeeAudioServer/models/localmodel"
	"BeeAudioServer/models/netmodel"
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

type PlayListController struct {
	MainController
}

/* 请求目标歌单详情 */
func (pc *PlayListController) GetPlayListDetail() {
	/* 需要获取的歌单ID */
	playlist_id, _ := pc.GetUint64("playlist_id")
	client := resty.New()
	var playlist_detail netmodel.PlayListInfo

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}

	if playlist_id < 690000000000 {
		playlist_detail = remote.RequestPlayListDetail(client, fmt.Sprint(playlist_id))

		if len(playlist_detail.PlayList) == 0 {
			/* 若没有获取到有效的艺人信息， 返回错误处理 */
			client_response.OK = 0
			client_response.Message = "No Exist PlayList"
			client_response.Data = nil
		} else {
			/* 若请求数据有效，返回有效的艺人详细信息数据 */
			playlist_detail.PlayList = remote.FilterSongs(utils.SqlDB, playlist_detail.PlayList)

			client_response.OK = 1
			client_response.Message = "success"
			client_response.Data = playlist_detail
		}
	} else {
		/* 自建歌单信息 */
		var playlist_model localmodel.UserPlayListModel
		/* 从服务器本地歌单中进行查找 */
		utils.SqlDB.Model(&localmodel.UserPlayListModel{}).
			Where("pid = ?", playlist_id).
			First(&playlist_model)

		/* 从查询结果中获取 创建者ID 从而进一步查询用户详细信息 */
		var creator_id uint64 = playlist_model.Uid
		var user_model localmodel.UserModel
		utils.SqlDB.Model(&localmodel.UserModel{}).
			Where("uid = ?", creator_id).
			First(&user_model)

		/* 获取自建歌单曲目内信息 */
		var song_ids []uint64
		var song_ids_str string
		utils.SqlDB.Model(&localmodel.PlayListTable{}).
			Select("song_id").
			Where("id = ?", playlist_id).
			Find(&song_ids)

		var track_ids []netmodel.TrackId
		for _, songId := range song_ids {
			track_ids = append(track_ids, netmodel.TrackId{
				SongId: songId,
			})
		}
		for index, item := range song_ids {
			if index == 0 {
				song_ids_str += fmt.Sprintf("%v", item)
			} else {
				if index > 20 {
					break
				}
				song_ids_str += fmt.Sprintf(",%v", item)
			}
		}
		/* 请求曲目详细信息 */
		playlist_songs := remote.RequestSongDetail(client, song_ids_str)

		playlist_songs = remote.FilterSongs(utils.SqlDB, playlist_songs)

		/* 组合这些查询信息 */
		playlist_detail = netmodel.PlayListInfo{
			Id:          playlist_model.Pid,
			Name:        playlist_model.Name,
			Cover:       playlist_model.Cover,
			UserId:      playlist_model.Uid,
			Description: playlist_model.Description,
			Tags: func(tag_string string) []string {
				var original_text = strings.Replace(strings.Replace(strings.Replace(tag_string, "[", "", -1), "]", "", -1), " ", "", -1)
				return strings.Split(original_text, ",")
			}(playlist_model.Tags),
			CreatorUser: netmodel.Creator{
				Id:            user_model.Uid,
				NickName:      user_model.Name,
				Signatrue:     "",
				Description:   user_model.Description,
				AvatarUrl:     user_model.AvatarUrl,
				BackgroundUrl: user_model.AvatarUrl,
			},
			PlayList:    playlist_songs,
			PlayListIds: track_ids,
		}
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = playlist_detail
	}

	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/* 请求目标歌单全部曲目列表 */
func (pc *PlayListController) GetPlayListAllSong() {
	/* 需要获取的歌单ID */
	playlist_id, _ := pc.GetUint64("playlist_id")
	playlist_page, _ := pc.GetInt("page", 1)
	playlist_size, _ := pc.GetInt("size", 20)

	client := resty.New()
	// fmt.Println("请求 page : ", playlist_page)

	var playlist_songs []netmodel.SongInfo
	/* 判断是否为自建歌单 */
	if playlist_id >= 690000000000 {
		var song_ids_str string
		song_ids := repository.GetPlayListAllSongID(utils.SqlDB, playlist_id, playlist_page, playlist_size)

		for index, item := range song_ids {
			if index == 0 {
				song_ids_str += fmt.Sprintf("%v", item)
			} else {
				song_ids_str += fmt.Sprintf(",%v", item)
			}
		}

		playlist_songs = remote.RequestSongDetail(client, song_ids_str)
	} else {
		playlist_songs = remote.RequestPlayListAllSong(client, fmt.Sprint(playlist_id), playlist_page, playlist_size)
	}

	/* 过滤信息 */
	playlist_songs = remote.FilterSongs(utils.SqlDB, playlist_songs)

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if len(playlist_songs) == 0 {
		/* 若没有获取到有效的艺人信息， 返回错误处理 */
		client_response.OK = 0
		client_response.Message = "No Exist PlayList"
		client_response.Data = nil
	} else {
		/* 若请求数据有效，返回有效的艺人详细信息数据 */
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = playlist_songs
	}
	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/* 请求热门歌单列表 */
func (pc *PlayListController) GetPlayListTops() {
	/* 需要获取的歌单ID */
	client := resty.New()
	playlist_collection := remote.RequestTopPlayList(client)
	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if len(playlist_collection) == 0 {
		/* 若没有获取到有效的艺人信息， 返回错误处理 */
		client_response.OK = 0
		client_response.Message = "No PlayList Collection Data"
		client_response.Data = nil
	} else {
		/* 若请求数据有效，返回有效的艺人详细信息数据 */
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = playlist_collection
	}
	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/* 请求指定Tag/Cat下的歌单列表 */
func (pc *PlayListController) GetTargetTagPlayLists() {
	/* 接收参数 */
	var cat string = pc.GetString("cat")
	playlist_page, _ := pc.GetInt("page", 1)
	playlist_size, _ := pc.GetInt("size", 50) // 设置当前请求指定Tag/Cat的歌单列表模式下，默认的容量大小为 50

	/* 需要获取的歌单ID */
	client := resty.New()
	playlist_collection := remote.RequestTargetTagsPlayLists(client, cat, playlist_page, playlist_size)
	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if len(playlist_collection) == 0 {
		/* 若没有获取到有效的艺人信息， 返回错误处理 */
		client_response.OK = 0
		client_response.Message = "No PlayList Collection Data"
		client_response.Data = nil
	} else {
		/* 若请求数据有效，返回有效的艺人详细信息数据 */
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = playlist_collection
	}
	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/* 查询用户默认喜爱歌单 */
func (pc *PlayListController) GetUserFavoritePlayLists() {
	/* 接收参数 */
	user_id, err := pc.GetUint64("user_id", 9900619251)

	/* 需要获取的歌单ID */
	playlistId := repository.QueryUserFavoritePlayList(utils.SqlDB, user_id)
	playlist := repository.GetUserPlayListInfo(utils.SqlDB, playlistId.Pid)

	if err != nil {
		log.Fatal(err)
	}

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	client_response.OK = 1
	client_response.Message = "success"
	client_response.Data = playlist.Pid
	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/* 向歌单中添加歌曲 */
func (pc *PlayListController) AddSongIntoPlayList() {
	var err error
	uid, err := pc.GetUint64("uid", 9900619251)
	pid, err := pc.GetUint64("pid", 690083340863)
	sid, err := pc.GetUint64("sid", 1897927507)

	client_response := model.ResponseBody{}

	/* 判断是否为创建者本人操作 */
	playlist_info := repository.GetUserPlayListInfo(utils.SqlDB, pid)
	if playlist_info.Uid != uid {
		/* 非创建者本人操作 */
		client_response.OK = 0
		client_response.Code = 403
		client_response.Message = "controll forbid"
		client_response.Data = false
	} else {
		result, _ := repository.AddSongIntoPlayList(utils.SqlDB, uid, pid, sid)
		client_response.OK = 1
		client_response.Code = 200
		client_response.Message = "success"
		client_response.Data = result
	}

	if err != nil {
		log.Fatal(err)
	}

	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/* 获取自建歌单中的曲目信息 */
func (pc *PlayListController) GetPlayListSongsDetail() {

	pid, err := pc.GetUint64("pid")
	playlist_page, _ := pc.GetInt("page", 1)
	playlist_size, _ := pc.GetInt("size", 20)
	client := resty.New()

	var song_ids_str string
	song_ids := repository.GetPlayListAllSongID(utils.SqlDB, pid, playlist_page, playlist_size)

	for index, item := range song_ids {
		if index == 0 {
			song_ids_str += fmt.Sprintf("%v", item)
		} else {
			song_ids_str += fmt.Sprintf(",%v", item)
		}
	}

	// fmt.Println(song_ids_str)

	song_detail_list := remote.RequestSongDetail(client, song_ids_str)

	if err != nil {
		log.Fatal(err)
	}
	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if len(song_detail_list) == 0 {
		/* 若没有获取到有效的艺人信息， 返回错误处理 */
		client_response.OK = 0
		client_response.Message = "No Songs Info"
		client_response.Data = nil
	} else {
		/* 若请求数据有效，返回有效的艺人详细信息数据 */
		song_detail_list = remote.FilterSongs(utils.SqlDB, song_detail_list)

		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = song_detail_list
	}
	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/* 通过表单信息变更目的曲目在目的歌单列表中的收录状态 */
func (pc *PlayListController) ChangeSongFromPlayLists() {
	var err error
	uid, err := pc.GetUint64("uid")
	sid, err := pc.GetUint64("sid")
	add_pids_str := pc.GetString("add_pids")
	remove_pids_str := pc.GetString("remove_pids")

	var add_pids []string = strings.Split(strings.Replace(add_pids_str, " ", "", -1), ",")
	var remove_pids []string = strings.Split(strings.Replace(remove_pids_str, " ", "", -1), ",")

	/* pid列表转换为uint64类型 */
	var add_pids_uint64, remove_pids_uint64 []uint64
	for _, pid := range add_pids {
		add_pids_uint64 = append(add_pids_uint64, utils.StringParseToUint64(pid))
	}
	for _, pid := range remove_pids {
		remove_pids_uint64 = append(remove_pids_uint64, utils.StringParseToUint64(pid))
	}

	/* 获取用户的自建歌单列表， 判断当前执行操作的用户是否有权限更改歌单 */
	var add_vaild_pid, remove_vaild_pid []uint64
	my_playlist_ids := repository.FindMyPlayList(utils.SqlDB, uid)

	for _, playlist_id := range my_playlist_ids {
		/* 添加列表ID */
		for _, add_id := range add_pids_uint64 {
			if playlist_id == add_id {
				add_vaild_pid = append(add_vaild_pid, playlist_id)
				break
			}
		}
		/* 删除列表ID */
		for _, remove_id := range remove_pids_uint64 {
			if playlist_id == remove_id {
				remove_vaild_pid = append(remove_vaild_pid, playlist_id)
				break
			}
		}
	}

	/* 批量执行曲目添加操作 */
	for _, add_id := range add_vaild_pid {
		err = utils.SqlDB.Model(&localmodel.PlayListTable{}).
			Clauses(clause.OnConflict{UpdateAll: true}).
			Create(&localmodel.PlayListTable{
				Pid:        add_id,
				Uid:        uid,
				SongId:     sid,
				CreateTime: time.Now(),
			}).Error
	}
	/* 批量执行删除操作 */
	for _, remove_id := range remove_vaild_pid {
		err = utils.SqlDB.Where("id = ? AND song_id = ?", remove_id, sid).
			Delete(&localmodel.PlayListTable{}).
			Error
	}

	client_response := model.ResponseBody{}
	if err == nil {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = true
	} else {
		client_response.OK = 0
		client_response.Message = "error"
		client_response.Data = false
	}

	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/**	从我的目的歌单中批量删除曲目
 *
 */
func (pc *PlayListController) BatchRemoveSongFromMyPlayList() {
	var err error
	uid, err := pc.GetUint64("uid")
	pid, err := pc.GetUint64("pid")
	song_ids_str := pc.GetString("song_ids")

	var song_ids []string = strings.Split(strings.Replace(song_ids_str, " ", "", -1), ",")

	client_response := model.ResponseBody{}

	/* 判断是否为创建者本人操作 */
	playlist_info := repository.GetUserPlayListInfo(utils.SqlDB, pid)
	if playlist_info.Uid != uid {
		/* 非创建者本人操作 */
		client_response.OK = 0
		client_response.Code = 403
		client_response.Message = "controll forbid"
		client_response.Data = false
	} else {
		/* 批量删除曲目 */
		for _, song_id := range song_ids {
			err = repository.RemoveSongInfoPlayList(utils.SqlDB, pid, utils.StringParseToUint64(song_id))
		}

		if err != nil {
			client_response.OK = 0
			client_response.Code = 200
			client_response.Message = "error"
			client_response.Data = false
		} else {
			client_response.OK = 1
			client_response.Code = 200
			client_response.Message = "success"
			client_response.Data = true
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	pc.Data["json"] = &client_response
	pc.ServeJSON()

}

/* 创建一个新歌单 */
func (pc *PlayListController) CreateMyPlayList() {
	/* 用户ID */
	uid, err := pc.GetUint64("uid")
	/* 歌单封面url - 初始化不必要 */
	coverImg := pc.GetString("cover_img", "/playlist/cover/default.jpg")

	/* 获取用户昵称 */
	user := repository.GetUserInfo(utils.SqlDB, uid)
	/* 歌单标题 */
	name := pc.GetString("name", fmt.Sprintf("用户%v创建的歌单", user.Name))
	/* 歌单描述文本 */
	description := pc.GetString("description", "ta的歌单，什么都没有写。")
	/* 附加Tag */
	tags := pc.GetString("tags")
	/* 歌单可见性 */
	public, err := pc.GetBool("public", true)

	/* 新建playlist */
	var temp_pid uint64

	/* ID去重 */
	for {
		temp_pid = utils.StringParseToUint64("6900" + utils.RandomNumberString(8, utils.DefaultNumber))
		var exist_signal int = 0
		utils.SqlDB.Model(&localmodel.UserPlayListModel{}).
			Select("Count(*)").
			Where("pid = ?", temp_pid).
			First(&exist_signal)
		/* 若不存在记录 */
		if exist_signal == 0 {
			break
		}
	}

	temp_user_playlist := localmodel.UserPlayListModel{
		Uid:         uid,
		Pid:         temp_pid,
		Cover:       coverImg,
		Name:        name,
		CreateTime:  time.Now(),
		Description: description,
		Tags:        tags,
		Public: func(public bool) int {
			if public {
				return 1 /* 公开 */
			} else {
				return 0 /* 不公开 */
			}
		}(public),
	}
	err = utils.SqlDB.Model(&localmodel.UserPlayListModel{}).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&temp_user_playlist).
		Error

	client_response := model.ResponseBody{}
	if err != nil {
		client_response.OK = 0
		client_response.Code = 200
		client_response.Message = "error"
		client_response.Data = nil
	} else {
		client_response.OK = 1
		client_response.Code = 200
		client_response.Message = "success"
		client_response.Data = temp_pid
	}
	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/* 删除我的自建歌单 */
func (pc *PlayListController) DeleteMyPlayList() {
	uid, err := pc.GetUint64("uid")
	pid, err := pc.GetUint64("pid")
	/* 判断是否为创建者本人操作 */
	playlist_info := repository.GetUserPlayListInfo(utils.SqlDB, pid)

	client_response := model.ResponseBody{}
	if playlist_info.Uid != uid {
		/* 非创建者本人操作 */
		client_response.OK = 0
		client_response.Code = 403
		client_response.Message = "controll forbid"
		client_response.Data = false
	} else {
		/* 歌单删除逻辑 */
		err = repository.DeleteMyPlayList(utils.SqlDB, pid)
		if err != nil {
			client_response.OK = 0
			client_response.Code = 200
			client_response.Message = "error"
			client_response.Data = false
		} else {
			client_response.OK = 1
			client_response.Code = 200
			client_response.Message = "success"
			client_response.Data = true
		}
	}

	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/* 更新自建歌单详细信息 */
func (pc *PlayListController) UpdateMyPlayList() {
	var err error
	/* 用户ID */
	uid, err := pc.GetUint64("uid")
	/* 更改歌单ID */
	pid, err := pc.GetUint64("pid")
	// /* 歌单封面url - 初始化不必要 */
	// coverImg := pc.GetString("cover_img", "/playlist/cover/default.jpg")

	/* 获取用户昵称 */
	user := repository.GetUserInfo(utils.SqlDB, uid)
	/* 歌单标题 */
	name := pc.GetString("name", fmt.Sprintf("用户%v的歌单", user.Name))
	/* 歌单描述文本 */
	description := pc.GetString("description", "ta的歌单，什么都没有写。")
	/* 附加Tag */
	tags := strings.Replace(strings.Replace(pc.GetString("tags"), "[", "", -1), "]", "", -1)
	/* 歌单可见性 */
	public, err := pc.GetBool("public", true)

	/* resty 网络请求客户端对象 */
	client := resty.New()

	client_response := model.ResponseBody{}
	if user.Uid == uid {
		/* 获取歌单原数据 */
		var original_playlist_data localmodel.UserPlayListModel
		err := utils.SqlDB.Model(&localmodel.UserPlayListModel{}).
			Where("pid = ?", pid).
			First(&original_playlist_data).Error

		/* 确认为歌单创建者本人的操作 */
		temp_user_playlist := localmodel.UserPlayListModel{
			Uid:         uid,
			Pid:         pid,
			Cover:       original_playlist_data.Cover,
			Name:        name,
			CreateTime:  original_playlist_data.CreateTime,
			Description: description,
			Tags:        "[" + tags + "]",
			Public: func(public bool) int {
				if public {
					return 1 /* 公开 */
				} else {
					return 0 /* 不公开 */
				}
			}(public),
		}
		/* 跟新歌单详细信息 */
		err = utils.SqlDB.Model(&localmodel.UserPlayListModel{}).
			Clauses(clause.OnConflict{UpdateAll: true}).
			Where("pid = ?", pid).
			Updates(&temp_user_playlist).Error
		/* 重新获取歌单信息， 当操作无误时将该信息作为响应返回 */
		/* 自建歌单信息 */
		var playlist_model localmodel.UserPlayListModel
		/* 从服务器本地歌单中进行查找 */
		utils.SqlDB.Model(&localmodel.UserPlayListModel{}).
			Where("pid = ?", pid).
			First(&playlist_model)

		/* 从查询结果中获取 创建者ID 从而进一步查询用户详细信息 */
		var creator_id uint64 = playlist_model.Uid
		var user_model localmodel.UserModel
		utils.SqlDB.Model(&localmodel.UserModel{}).
			Where("uid = ?", creator_id).
			First(&user_model)

		/* 获取自建歌单曲目内信息 */
		var song_ids []uint64
		var song_ids_str string
		utils.SqlDB.Model(&localmodel.PlayListTable{}).
			Select("song_id").
			Where("id = ?", pid).
			Find(&song_ids)
		for index, item := range song_ids {
			if index == 0 {
				song_ids_str += fmt.Sprintf("%v", item)
			} else {
				song_ids_str += fmt.Sprintf(",%v", item)
			}
		}
		/* 请求曲目详细信息 */
		playlist_songs := remote.RequestSongDetail(client, song_ids_str)
		/* 曲目清单只预览前20首 */
		if len(playlist_songs) > 20 {
			playlist_songs = playlist_songs[:20]
		}

		playlist_songs = remote.FilterSongs(utils.SqlDB, playlist_songs)
		var track_ids []netmodel.TrackId
		for _, song := range playlist_songs {
			track_ids = append(track_ids, netmodel.TrackId{
				SongId: song.SongId,
			})
		}

		/* 组合这些查询信息 */
		playlist_detail := netmodel.PlayListInfo{
			Id:          playlist_model.Pid,
			Name:        playlist_model.Name,
			Cover:       playlist_model.Cover,
			UserId:      playlist_model.Uid,
			Description: playlist_model.Description,
			Tags: func(original_tags_str string) []string {
				return strings.Split(strings.Replace(original_tags_str, " ", "", -1), ",")
			}(playlist_model.Tags),
			CreatorUser: netmodel.Creator{
				Id:            user_model.Uid,
				NickName:      user_model.Name,
				Signatrue:     "",
				Description:   user_model.Description,
				AvatarUrl:     user_model.AvatarUrl,
				BackgroundUrl: user_model.AvatarUrl,
			},
			PlayList:    playlist_songs,
			PlayListIds: track_ids,
		}

		if err != nil {
			client_response.OK = 0
			client_response.Code = 200
			client_response.Message = "error"
			client_response.Data = nil
		} else {
			client_response.OK = 1
			client_response.Code = 200
			client_response.Message = "success"
			client_response.Data = playlist_detail
		}

	} else {
		/* 返回无权限警告响应信息 */
		/* 非创建者本人操作 */
		if err != nil {
			client_response.OK = 0
			client_response.Code = 200
			client_response.Message = "error"
			client_response.Data = nil
		} else {
			client_response.OK = 0
			client_response.Code = 403
			client_response.Message = "controll forbid"
			client_response.Data = nil
		}
	}
	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/* 请求访问用户创建的歌单 - 在公开性为开放的情况下 */
func (pc *PlayListController) AccessUserPublicPlaylistCollection() {
	var err error
	/* 用户ID */
	uid, err := pc.GetUint64("uid")
	/* resty 网络请求客户端对象 */
	client := resty.New()

	/* 查询用户公开可访问的自建歌单列表 */
	client_response := model.ResponseBody{}
	collections, err := repository.AccessUserPublicPlayList(utils.SqlDB, uid)

	/* 歌单列表集合 */
	var access_playlist_collection []netmodel.PlayListInfo

	/* 获取用户信息 */
	var user_model localmodel.UserModel
	utils.SqlDB.Model(&localmodel.UserModel{}).
		Where("uid = ?", uid).
		First(&user_model)

	for _, created_playlist := range collections {

		/* 获取自建歌单曲目内信息 */
		var song_ids []uint64
		var song_ids_str string
		utils.SqlDB.Model(&localmodel.PlayListTable{}).
			Select("song_id").
			Where("id = ?", created_playlist.Pid).
			Find(&song_ids)

		var track_ids []netmodel.TrackId
		for _, songId := range song_ids {
			track_ids = append(track_ids, netmodel.TrackId{
				SongId: songId,
			})
		}
		for index, item := range song_ids {
			if index == 0 {
				song_ids_str += fmt.Sprintf("%v", item)
			} else {
				if index > 20 {
					break
				}
				song_ids_str += fmt.Sprintf(",%v", item)
			}
		}
		/* 请求曲目详细信息 */
		playlist_songs := remote.RequestSongDetail(client, song_ids_str)

		playlist_songs = remote.FilterSongs(utils.SqlDB, playlist_songs)

		/* 组合这些查询信息 */
		playlist_detail := netmodel.PlayListInfo{
			Id:          created_playlist.Pid,
			Name:        created_playlist.Name,
			Cover:       created_playlist.Cover,
			UserId:      created_playlist.Uid,
			Description: created_playlist.Description,
			Tags: func(tag_string string) []string {
				var original_text = strings.Replace(strings.Replace(strings.Replace(tag_string, "[", "", -1), "]", "", -1), " ", "", -1)
				return strings.Split(original_text, ",")
			}(created_playlist.Tags),
			CreatorUser: netmodel.Creator{
				Id:            user_model.Uid,
				NickName:      user_model.Name,
				Signatrue:     "",
				Description:   user_model.Description,
				AvatarUrl:     user_model.AvatarUrl,
				BackgroundUrl: user_model.AvatarUrl,
			},
			PlayList:    playlist_songs,
			PlayListIds: track_ids,
		}
		access_playlist_collection = append(access_playlist_collection, playlist_detail)
	}

	if err == nil {
		client_response.OK = 1
		client_response.Code = 200
		client_response.Message = "success"
		client_response.Data = access_playlist_collection
	} else {
		client_response.OK = 0
		client_response.Code = 500
		client_response.Message = "error"
		client_response.Data = nil
	}
	pc.Data["json"] = &client_response
	pc.ServeJSON()
}
