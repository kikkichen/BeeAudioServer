package repository

import (
	"BeeAudioServer/models/localmodel"
	"BeeAudioServer/models/netmodel"
	"BeeAudioServer/remote"
	"BeeAudioServer/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"gorm.io/gorm/clause"
)

/**	插入 / 更新 一条曲目数据到历史记录
 *	@param	client	resty网络请求客户端对象
 *	@param	uid		当前更新历史记录用户ID
 *	@param	songId	用户当前触发更新历史记录元数据的曲目ID
 */
func AddPlayHistory(client *resty.Client, uid, songId uint64) error {
	var err error

	/* 请求获取音频详细信息 */
	var remote_result []netmodel.SongInfo = remote.RequestSongDetail(client, fmt.Sprintf("%v", songId))
	/* 若音频信息结果有效 */
	if len(remote_result) != 0 {
		song_info := remote_result[0]
		/* 构建一个用于填充的新的历史记录信息 */
		var new_history_record localmodel.HistorySong = localmodel.HistorySong{
			// 曲目ID
			SongId: song_info.SongId,
			// 曲目名字
			SongTitle: song_info.Name,
			// 曲目艺人ID列表
			ArtistIds: func(artist_list []netmodel.InnerSongInfoAr) []uint64 {
				var artist_ids []uint64
				for _, artist := range artist_list {
					artist_ids = append(artist_ids, artist.ArId)
				}
				return artist_ids
			}(song_info.Ar),
			// 曲目所属专辑ID
			AlbumIds: song_info.Al.AlId,
			// 当前时间
			PlayAt: time.Now().Unix(),
		}

		/* 向数据库中查询当前用户的历史播放记录 */
		/* 若数据库还不存在当前使用用户的播放列表记录数据， 则为其添加数据 */
		var isNewHand bool = func(uid uint64) bool {
			var record_number int
			// 插叙记录
			err = utils.SqlDB.Model(&localmodel.HistoryDataModel{}).Select("Count(*)").Where("uid = ?", uid).First(&record_number).Error
			if record_number == 0 {
				return true
			} else {
				return false
			}
		}(uid)

		if isNewHand {
			// 该用户还没拥有记录
			var new_record_list []localmodel.HistorySong
			new_record_list = append(new_record_list, new_history_record)
			// 将列表构建为json数据
			json_new_record_list, _ := json.Marshal(new_record_list)
			// 构建 gorm 历史播放记录模型对象
			record_data := localmodel.HistoryDataModel{
				Uid:         uid,
				HistoryData: string(json_new_record_list),
			}
			// 插入记录
			err = utils.SqlDB.Model(&localmodel.HistoryDataModel{}).
				Clauses(clause.OnConflict{UpdateAll: true}).
				Create(&record_data).
				Error
		} else {
			/* 该用户已经拥有了记录 */
			/* 获取用户的记录 */
			var origin_record_data localmodel.HistoryDataModel
			err = utils.SqlDB.Model(&localmodel.HistoryDataModel{}).
				Where("uid = ?", uid).
				First(&origin_record_data).Error
			/* 解析json记录数据 */
			var origin_data []localmodel.HistorySong
			err = json.Unmarshal([]byte(origin_record_data.HistoryData), &origin_data)
			if err != nil {
				return err
			}
			/* 首先将当前历史播放元数据记录插入 */
			var new_list_data []localmodel.HistorySong
			new_list_data = append(new_list_data, new_history_record)

			var count int = 1
			/* 将非该条数据的旧数据依次拷贝到新列表中 */
			for _, meta_history_data := range origin_data {
				/* 若历史播放记录条目超过500条，则按照队列形式更新历史播放记录数据 */
				if count >= 500 {
					break
				}
				if meta_history_data.SongId != new_history_record.SongId {
					new_list_data = append(new_list_data, meta_history_data)
					count++
				}
			}
			/* 将新数据解析为json数据 */
			json_new_list_data, _ := json.Marshal(new_list_data)

			/* gorm 模型数据对象 */
			model_data := localmodel.HistoryDataModel{
				Uid:         uid,
				HistoryData: string(json_new_list_data),
			}
			/* 将新数据更新到数据库中 */
			err = utils.SqlDB.Model(&localmodel.HistoryDataModel{}).
				Clauses(clause.OnConflict{UpdateAll: true}).
				Create(&model_data).Error
			if err != nil {
				return err
			}
		}
	}
	return err
}

/**	查看当前用户的播放历史记录， 返回 曲目信息、 历史信息 与 结果有效性标识
 *	@param	client	resty网络请求客户端对象
 *	@param	uid		当前请求历史播放记录的用户ID
 */
func BrowserPlayHistory(client *resty.Client, uid uint64, page, size int) ([]netmodel.SongInfo, []localmodel.HistorySong, error) {
	var err error
	/* 获取用户的记录 */
	var origin_record_data localmodel.HistoryDataModel
	/* 若该用户为新用户， ta的历史记录信息还未在数据库中初始化, 则直接返回空数据 */
	var isNewHand bool = func(uid uint64) bool {
		var record_number int
		// 插叙记录
		err = utils.SqlDB.Model(&localmodel.HistoryDataModel{}).Select("Count(*)").Where("uid = ?", uid).First(&record_number).Error
		if record_number == 0 {
			return true
		} else {
			return false
		}
	}(uid)

	if isNewHand {
		return []netmodel.SongInfo{}, []localmodel.HistorySong{}, nil
	} else {
		err = utils.SqlDB.Model(&localmodel.HistoryDataModel{}).
			Where("uid = ?", uid).
			First(&origin_record_data).Error
		/* 解析json记录数据 */
		var origin_data []localmodel.HistorySong
		err = json.Unmarshal([]byte(origin_record_data.HistoryData), &origin_data)
		if err != nil {
			// 错误则返回空数据
			return []netmodel.SongInfo{}, []localmodel.HistorySong{}, err
		}

		/* 判断历史记录是否为空 */
		if len(origin_data) != 0 {
			/* 对应页码、对应容量的历史记录数据 */
			var slice_origin_data []localmodel.HistorySong
			/* 根据分页参数，针对返回结果分类讨论 */
			if ((page-1)*size) >= len(origin_data) && (page*size) > len(origin_data) {
				/* 查询最大范围越界情况 */
				return []netmodel.SongInfo{}, []localmodel.HistorySong{}, nil

			} else if ((page-1)*size) < len(origin_data) && (page*size) > len(origin_data) {
				/* 查询最后一页数据 */
				slice_origin_data = origin_data[((page - 1) * size):]
			} else {
				/* 查询结果恰好都在范围内 */
				slice_origin_data = origin_data[((page - 1) * size):(page * size)]
			}

			/* 提取曲目ID */
			var song_ids []uint64
			for _, data := range slice_origin_data {
				song_ids = append(song_ids, data.SongId)
			}
			/* 获取字符串类型的曲目ID集合 */
			var song_ids_str string
			for index, item := range song_ids {
				if index == 0 {
					song_ids_str += fmt.Sprintf("%v", item)
				} else {
					song_ids_str += fmt.Sprintf(",%v", item)
				}
			}
			/* 获取曲目集合详细信息 */
			history_songs := remote.RequestSongDetail(client, song_ids_str)
			history_songs = remote.FilterSongs(utils.SqlDB, history_songs)
			/* 返回正确结果 */
			return history_songs, slice_origin_data, nil
		} else {
			// 正确请求无记录下， 返回空数据
			return []netmodel.SongInfo{}, []localmodel.HistorySong{}, nil
		}
	}
}

/**	清空历史记录
 *	@param	uid	当前清空历史播放记录的用户ID
 */
func ClearPlayHistory(uid uint64) error {
	var meta_empty_history []localmodel.HistorySong
	json_empty_data, err := json.Marshal(meta_empty_history)

	var new_empty_history localmodel.HistoryDataModel = localmodel.HistoryDataModel{
		Uid:         uid,
		HistoryData: string(json_empty_data),
	}
	/* 使用更新空数据的方式，清空用户数据 */
	err = utils.SqlDB.Model(&localmodel.HistoryDataModel{}).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&new_empty_history).Error
	return err
}
