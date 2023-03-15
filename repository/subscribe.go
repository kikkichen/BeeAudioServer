package repository

import (
	"BeeAudioServer/models/localmodel"
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/**	同步项目项目订阅信息
 *	@param	db	gorm 数据库连接对象
 *	@param	uid	目标操作用户ID
 *	@param	data	订阅列表数据
 */
func SyncSubscribeData(
	db *gorm.DB,
	uid uint64,
	data string,
) error {
	return db.Model(&localmodel.SubscribeModel{}).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Where("uid = ?", uid).
		Create(&localmodel.SubscribeModel{
			Uid:  uid,
			Data: data,
		}).Error
}

/**	获取用户项目订阅信息
 *	@param	db	gorm 数据库连接对象
 *	@param	uid	目标操作用户ID
 */
func GetSubscribeData(
	db *gorm.DB,
	uid uint64,
) (localmodel.SubscribeModel, error) {
	var subscribe_model_result localmodel.SubscribeModel
	var exist_signal int
	err := db.Model(&localmodel.SubscribeModel{}).
		Select("Count(*)").
		Where("uid = ?", uid).
		First(&exist_signal).Error

	/* 当前用户是否存在订阅列表信息记录 */
	if exist_signal > 0 {
		err = db.Model(&localmodel.SubscribeModel{}).
			Where("uid = ?", uid).
			First(&subscribe_model_result).
			Error
		return subscribe_model_result, err
	} else {
		/* 初次构建音频订阅列表数据，需要将用户自己的默认喜爱歌单添加 */
		favoritePlayListId := QueryUserFavoritePlayList(db, uid).Pid
		favoritePlayListInfo := GetUserPlayListInfo(db, favoritePlayListId)
		myInfo := GetUserInfo(db, uid)

		temp_map := map[string]interface{}{
			"cover_url":     favoritePlayListInfo.Cover,
			"creator":       myInfo.Name,
			"is_my_created": true,
			"is_top":        true,
			"item_id":       favoritePlayListId,
			"title":         favoritePlayListInfo.Name,
			"item_type":     1000,
			"weight":        65535,
		}

		jsonPlaylistData, err := json.Marshal(temp_map)

		subscribe_model_result := localmodel.SubscribeModel{
			Uid:  uid,
			Data: "[" + string(jsonPlaylistData) + "]",
		}
		err = db.Model(&localmodel.SubscribeModel{}).
			Clauses(clause.OnConflict{UpdateAll: true}).
			Create(&subscribe_model_result).
			Error
		return subscribe_model_result, err
	}
}
