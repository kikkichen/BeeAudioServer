package controllers

import (
	model "BeeAudioServer/models"
	"BeeAudioServer/models/localmodel"
	"BeeAudioServer/remote"
	"BeeAudioServer/repository"
	"BeeAudioServer/utils"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
)

type SongController struct {
	MainController
}

/* 获取曲目详情 */
func (sc *SongController) GetSongDetail() {
	var song_ids string = sc.GetString("song_ids")
	client := resty.New()
	song_detail_list := remote.RequestSongDetail(client, song_ids)
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
	sc.Data["json"] = &client_response
	sc.ServeJSON()
}

/* 获取音频URL */
func (sc *SongController) GetSongUrl() {
	var song_ids string = sc.GetString("song_ids")
	client := resty.New()
	song_urls := remote.RequestAudioUrl(client, song_ids, "standard")

	/* ID字符转列表转换为为长整形列表 */
	var song_id_list_string string = strings.Replace(strings.Replace(strings.Replace(song_ids, "[", "", -1), "]", "", -1), " ", "", -1)
	var song_id_list []string = strings.Split(song_id_list_string, ",")
	var song_id_list_uint []uint64
	for _, id := range song_id_list {
		song_id_list_uint = append(song_id_list_uint, utils.StringParseToUint64(id))
	}

	var err error
	for index, song_id := range song_id_list_uint {
		/* 数据库匹配覆盖 */
		var exist_signal int
		err = utils.SqlDB.Model(&localmodel.SongTable{}).
			Select("COUNT(*)").
			Where("id = ?", song_id).
			First(&exist_signal).Error

		if exist_signal > 0 {
			local_url := repository.GetLocalAudioUrl(utils.SqlDB, song_id)
			song_urls[index].Url = local_url
		}
	}

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if err == nil {
		if len(song_urls) == 0 {
			/* 若没有获取到有效的艺人信息， 返回错误处理 */
			client_response.OK = 0
			client_response.Message = "No Effective Link"
			client_response.Data = nil
		} else {
			/* 若请求数据有效，返回有效的艺人详细信息数据 */
			client_response.OK = 1
			client_response.Message = "success"
			client_response.Data = song_urls
		}
	} else {
		client_response.OK = 0
		client_response.Message = "server error"
		client_response.Data = nil
	}
	sc.Data["json"] = &client_response
	sc.ServeJSON()
}

/* 获取单个音频本地Url */
func (sc *SongController) GetLocalAudioUrl() {
	song_id, err := sc.GetUint64("song_id")
	song_local_url := repository.GetLocalAudioUrl(utils.SqlDB, song_id)

	if err != nil {
		log.Fatal(err)
	}

	client_response := model.ResponseBody{}
	if len(song_local_url) == 0 {
		/* 若没有获取到有效的艺人信息， 返回错误处理 */
		client_response.OK = 0
		client_response.Message = "No Effective Link"
		client_response.Data = nil
	} else {
		/* 若请求数据有效，返回有效的艺人详细信息数据 */
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = song_local_url
	}
	sc.Data["json"] = &client_response
	sc.ServeJSON()
}

/**	管理员浏览曲目列表 - 分页
 *
 */
func (sc *SongController) GetLocalSongList() {
	sort, err := sc.GetBool("sort", false)
	page, err := sc.GetInt("page", 1)
	size, err := sc.GetInt("size", 20)

	/* 查询 */
	song_list, err := repository.BrowserSongs(utils.SqlDB, page, size, sort)

	client_response := model.ResponseBody{}
	if err == nil {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = song_list
	} else {
		sc.Ctx.ResponseWriter.WriteHeader(500)
		client_response.OK = 0
		client_response.Message = "error"
		client_response.Data = []repository.SongTrackDetail{}
	}
	sc.Data["json"] = &client_response
	sc.ServeJSON()
}

/**	管理员通过曲目名字检索曲目列表 - 分页
 *
 */
func (sc *SongController) GetLocalSongListBySongName() {
	sort, err := sc.GetBool("sort", false)
	page, err := sc.GetInt("page", 1)
	size, err := sc.GetInt("size", 20)
	name := sc.GetString("song_name")

	/* 查询 */
	song_list, err := repository.BrowserSongsBySongName(utils.SqlDB, name, page, size, sort)

	client_response := model.ResponseBody{}
	if err == nil {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = song_list
	} else {
		sc.Ctx.ResponseWriter.WriteHeader(500)
		client_response.OK = 0
		client_response.Message = "error"
		client_response.Data = []repository.SongTrackDetail{}
	}
	sc.Data["json"] = &client_response
	sc.ServeJSON()
}

/**	管理员通过专辑名字检索曲目列表 - 分页
 *
 */
func (sc *SongController) GetLocalSongListByAlbumName() {
	sort, err := sc.GetBool("sort", false)
	page, err := sc.GetInt("page", 1)
	size, err := sc.GetInt("size", 20)
	name := sc.GetString("album_name")

	/* 查询 */
	song_list, err := repository.BrowserSongsByAlbumName(utils.SqlDB, name, page, size, sort)

	client_response := model.ResponseBody{}
	if err == nil {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = song_list
	} else {
		sc.Ctx.ResponseWriter.WriteHeader(500)
		client_response.OK = 0
		client_response.Message = "error"
		client_response.Data = []repository.SongTrackDetail{}
	}
	sc.Data["json"] = &client_response
	sc.ServeJSON()
}

/**	管理员通过艺人名字检索曲目列表 - 分页
 *
 */
func (sc *SongController) GetLocalSongListByArtistName() {
	sort, err := sc.GetBool("sort", false)
	page, err := sc.GetInt("page", 1)
	size, err := sc.GetInt("size", 20)
	name := sc.GetString("artist_name")

	/* 查询 */
	song_list, err := repository.BrowserSongsByArtistName(utils.SqlDB, name, page, size, sort)

	client_response := model.ResponseBody{}
	if err == nil {
		/* 修改正常情况 */
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = song_list
	} else {
		sc.Ctx.ResponseWriter.WriteHeader(500)
		client_response.OK = 0
		client_response.Message = "server error"
		client_response.Data = []repository.SongTrackDetail{}
	}
	sc.Data["json"] = &client_response
	sc.ServeJSON()
}

/**	管理员通过ID检索到曲目
 *
 */
func (sc *SongController) SelectSongById() {
	song_id, err := sc.GetUint64("song_id")
	/* 检索 */
	target_song, err := repository.SelectSongById(utils.SqlDB, song_id)

	client_response := model.ResponseBody{}
	if err == nil {
		if target_song.SongId != 0 {
			client_response.OK = 1
			client_response.Message = "success"
			client_response.Data = target_song
		} else {
			client_response.OK = 0
			client_response.Message = "no result"
			client_response.Data = nil
		}
	} else {
		/* 修改错误情况 */
		sc.Ctx.ResponseWriter.WriteHeader(500)
		client_response.OK = 0
		client_response.Message = "server error"
		client_response.Data = []repository.SongTrackDetail{}
	}
	sc.Data["json"] = &client_response
	sc.ServeJSON()
}

/**	管理员修改音频信息
 *
 */
func (sc *SongController) ModifierSongDetail() {
	song_id, err := sc.GetUint64("song_id")
	privilege, err := sc.GetInt("privilege")
	useful, err := sc.GetInt("useful")
	source := sc.GetString("source")

	target_song, err := repository.UpdateSongDetail(utils.SqlDB, song_id, privilege, useful, source)

	client_response := model.ResponseBody{}
	if err == nil {
		if target_song.SongId != 0 {
			client_response.OK = 1
			client_response.Message = "success"
			client_response.Data = target_song
		} else {
			client_response.OK = 0
			client_response.Message = "no result"
			client_response.Data = nil
		}
	} else {
		/* 修改错误情况 */
		sc.Ctx.ResponseWriter.WriteHeader(500)
		client_response.OK = 0
		client_response.Message = "server error"
		client_response.Data = []repository.SongTrackDetail{}
	}
	sc.Data["json"] = &client_response
	sc.ServeJSON()
}
