package repository

import (
	"BeeAudioServer/models/localmodel"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/**	查询用户的默认歌单
 *	@param	db	gorm 数据库连接对象
 *	@param	uid	查询目的用户ID
 */
func QueryUserFavoritePlayList(
	db *gorm.DB,
	uid uint64,
) localmodel.UserFavoritePlaylist {
	var favoritePlayList localmodel.UserFavoritePlaylist
	db.Model(&localmodel.UserFavoritePlaylist{}).
		Where("uid = ?", uid).
		First(&favoritePlayList)
	return favoritePlayList
}

/**	根据Pid获取对应的自建歌单信息
 *	@param	db	gorm 数据库连接对象
 *	@param	pid	查询自建歌单ID （一般ID以 6900 开头）
 */
func GetUserPlayListInfo(
	db *gorm.DB,
	pid uint64,
) localmodel.UserPlayListModel {
	var playlist localmodel.UserPlayListModel
	db.Model(&localmodel.UserPlayListModel{}).
		Where("pid = ?", pid).
		First(&playlist)
	return playlist
}

/**	依据用户ID访问用户歌单
 *	@param	db	gorm 数据库连接对象
 *	@param	uid	目标查询用户ID
 */
func AccessUserPublicPlayList(
	db *gorm.DB,
	uid uint64,
) ([]localmodel.UserPlayListModel, error) {
	var user_created_playlist_collection []localmodel.UserPlayListModel
	err := db.Model(&localmodel.UserPlayListModel{}).
		Where("uid = ? AND public > 0", uid).
		Find(&user_created_playlist_collection).
		Error
	return user_created_playlist_collection, err
}

/**	将歌曲添加到歌单
 *	@param	db	gorm 数据库连接对象
 *	@param	uid	目的用户ID
 *	@param	pid	目的歌单ID
 *	@param	songId	曲目ID
 */
func AddSongIntoPlayList(
	db *gorm.DB,
	uid uint64,
	pid uint64,
	songId uint64,
) (bool, error) {
	insert_record := localmodel.PlayListTable{
		Pid:        pid,
		Uid:        uid,
		SongId:     songId,
		CreateTime: time.Now(),
	}
	err := db.Model(&localmodel.PlayListTable{}).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&insert_record).Error
	return true, err
}

/**	将歌曲从目标歌单中移除 （即便目标歌单中原本没有该曲目收录）
 *	@param	db	gorm 数据库连接对象
 *	@param	pid	目的歌单ID
 *	@param	songId	曲目ID
 */
func RemoveSongInfoPlayList(
	db *gorm.DB,
	pid uint64,
	songId uint64,
) error {
	err := db.Where("id = ? AND song_id = ?", pid, songId).
		Delete(&localmodel.PlayListTable{}).
		Error
	return err
}

/**	查询当前曲目在目标用户创建的哪些歌单中存在
 *	@param	db	gorm 数据库连接对象
 *	@param	uid	查询目标自建歌单创建者ID
 *	@param	sid	查询曲目ID
 */
func FindCurrentSongExistInMyPlayList(
	db *gorm.DB,
	uid uint64,
	sid uint64,
) []uint64 {
	var exist_playlist_ids []uint64
	db.Model(&localmodel.PlayListTable{}).
		Select("id").
		Where("uid = ? AND song_id = ?", uid, sid).
		Find(&exist_playlist_ids)
	return exist_playlist_ids
}

/**	获取我所有的自建歌单
 *	@param	db	gorm 数据库连接对象
 *	@param	uid	查询目标自建歌单创建者ID
 */
func FindMyPlayList(
	db *gorm.DB,
	uid uint64,
) []uint64 {
	var my_playlist_ids []uint64
	db.Model(&localmodel.UserPlayListModel{}).
		Select("pid").
		Where("uid = ?", uid).
		Find(&my_playlist_ids)
	return my_playlist_ids
}

/**	获取歌单中的歌曲ID
 *	@param	db	gorm 数据库连接对象
 *	@param	pid	目的歌单ID
 */
func GetPlayListAllSongID(
	db *gorm.DB,
	pid uint64,
	page int,
	size int,
) []uint64 {
	var song_ids []uint64
	db.Model(&localmodel.PlayListTable{}).
		Select("song_id").
		Where("id = ?", pid).
		Offset((page - 1) * size).
		Limit(size).
		Order("create_at DESC").
		Scan(&song_ids)
	return song_ids
}

/**	创建我的一个歌单
 *	@param	db	gorm 数据库连接对象
 *	@param	new_playlist	新建歌单信息
 */
func CreateMyPlayList(
	db *gorm.DB,
	new_playlist localmodel.UserPlayListModel,
) error {
	return db.Model(&localmodel.UserPlayListModel{}).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(new_playlist).Error
}

/* 删除我的一个歌单 */
func DeleteMyPlayList(
	db *gorm.DB,
	pid uint64,
) error {
	var err error
	err = db.
		Where("id = ?", pid).
		Delete(&localmodel.PlayListTable{}).Error
	if err != nil {
		return err
	}
	err = db.
		Where("pid = ?", pid).
		Delete(&localmodel.UserPlayListModel{}).Error
	if err != nil {
		return err
	}
	return nil
}
