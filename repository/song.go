package repository

import (
	"BeeAudioServer/models/localmodel"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

/* 用户管理端显示的曲目格式 */
type SongTrackDetail struct {
	SongId     uint64 `gorm:"column:id" json:"song_id"`
	SongName   string `gorm:"column:name" json:"song_name"`
	ArtistId   uint64 `gorm:"column:ar_id" json:"ar_id"`
	ArtistName string `gorm:"column:ar_name" json:"ar_name"`
	AlbumId    uint64 `gorm:"column:al_id" json:"al_id"`
	AlbumName  string `gorm:"column:al_name" json:"al_name"`
	LocalPath  string `gorm:"column:local_path" json:"local_path"`
	Privilege  int    `gorm:"column:privilege" json:"privilege"`
	Quality    string `gorm:"column:quality" json:"quality"`
	Useful     int    `gorm:"column:useful" json:"useful"`
	Source     string `gorm:"column:source" json:"source"`
}

/**	获取本地音频Url
 *	@param	db	gorm数据库连接对象
 *	@param	song_id	目标获取url的曲目ID
 */
func GetLocalAudioUrl(
	db *gorm.DB,
	song_id uint64,
) string {
	var local_path string
	db.Model(&localmodel.SongTable{}).
		Select("local_path").
		Where("id = ?", song_id).
		First(&local_path)

	path_args := strings.Split(local_path, "/")
	var url string = ""
	if len(path_args) != 0 {
		url = fmt.Sprintf("http://192.168.43.240:9093/%v/%v/%v/%v/%v.m3u8", path_args[0], path_args[1], path_args[2], path_args[3], path_args[3])
	}
	return url
}

/**	浏览曲目列表 - 分页
 *	@param	db	gorm数据库连接对象
 *	@param	page	页码
 *	@param	size	单页容量
 *	@param	sort	按照曲目ID大小排序	true为由小到大 false为又大到小
 */
func BrowserSongs(
	db *gorm.DB,
	page, size int,
	sort bool,
) ([]SongTrackDetail, error) {
	var err error

	var song_list []SongTrackDetail
	err = db.Model(&localmodel.SongTable{}).
		Select("song_table.id, song_table.name, song_table.ar_id, artist_table.name as 'ar_name' ,song_table.al_id, album_table.name as 'al_name', local_path, privilege, quality, useful, source").
		Joins("LEFT JOIN mblog.album_table ON song_table.al_id = album_table.id LEFT JOIN mblog.artist_table ON song_table.ar_id = artist_table.id").
		Order(func(bool) string {
			if sort {
				return "song_table.id ASC"
			} else {
				return "song_table.id DESC"
			}
		}(sort)).
		Limit(size).
		Offset((page - 1) * size).
		Find(&song_list).
		Error
	return song_list, err
}

/**	通过曲目名字浏览曲目列表 - 分页
 *	@param	db	gorm数据库连接对象
 *	@param	name	曲目名字关键字字符串
 *	@param	page	页码
 *	@param	size	单页容量
 *	@param	sort	按照曲目ID大小排序	true为由小到大 false为又大到小
 */
func BrowserSongsBySongName(
	db *gorm.DB,
	name string,
	page, size int,
	sort bool,
) ([]SongTrackDetail, error) {
	var err error

	var song_list []SongTrackDetail
	err = db.Model(&localmodel.SongTable{}).
		Select("song_table.id, song_table.name, song_table.ar_id, artist_table.name as 'ar_name' ,song_table.al_id, album_table.name as 'al_name', local_path, privilege, quality, useful, source").
		Joins("LEFT JOIN mblog.album_table ON song_table.al_id = album_table.id LEFT JOIN mblog.artist_table ON song_table.ar_id = artist_table.id").
		Where("song_table.name LIKE ?", "%"+name+"%").
		Order(func(bool) string {
			if sort {
				return "song_table.id ASC"
			} else {
				return "song_table.id DESC"
			}
		}(sort)).
		Limit(size).
		Offset((page - 1) * size).
		Find(&song_list).
		Error
	return song_list, err
}

/**	通过专辑名字浏览曲目列表 - 分页
 *	@param	db	gorm数据库连接对象
 *	@param	name	专辑名字关键字字符串
 *	@param	page	页码
 *	@param	size	单页容量
 *	@param	sort	按照曲目ID大小排序	true为由小到大 false为又大到小
 */
func BrowserSongsByAlbumName(
	db *gorm.DB,
	name string,
	page, size int,
	sort bool,
) ([]SongTrackDetail, error) {
	var err error

	var song_list []SongTrackDetail
	err = db.Model(&localmodel.SongTable{}).
		Select("song_table.id, song_table.name, song_table.ar_id, artist_table.name as 'ar_name' ,song_table.al_id, album_table.name as 'al_name', local_path, privilege, quality, useful, source").
		Joins("LEFT JOIN mblog.album_table ON song_table.al_id = album_table.id LEFT JOIN mblog.artist_table ON song_table.ar_id = artist_table.id").
		Where("album_table.name LIKE ?", "%"+name+"%").
		Order(func(bool) string {
			if sort {
				return "song_table.id ASC"
			} else {
				return "song_table.id DESC"
			}
		}(sort)).
		Limit(size).
		Offset((page - 1) * size).
		Find(&song_list).
		Error
	return song_list, err
}

/**	通过艺人名字浏览曲目列表 - 分页
 *	@param	db	gorm数据库连接对象
 *	@param	name	艺人名字关键字字符串
 *	@param	page	页码
 *	@param	size	单页容量
 *	@param	sort	按照曲目ID大小排序	true为由小到大 false为又大到小
 */
func BrowserSongsByArtistName(
	db *gorm.DB,
	name string,
	page, size int,
	sort bool,
) ([]SongTrackDetail, error) {
	var err error

	var song_list []SongTrackDetail
	err = db.Model(&localmodel.SongTable{}).
		Select("song_table.id, song_table.name, song_table.ar_id, artist_table.name as 'ar_name' ,song_table.al_id, album_table.name as 'al_name', local_path, privilege, quality, useful, source").
		Joins("LEFT JOIN mblog.album_table ON song_table.al_id = album_table.id LEFT JOIN mblog.artist_table ON song_table.ar_id = artist_table.id").
		Where("artist_table.name LIKE ?", "%"+name+"%").
		Order(func(bool) string {
			if sort {
				return "song_table.id ASC"
			} else {
				return "song_table.id DESC"
			}
		}(sort)).
		Limit(size).
		Offset((page - 1) * size).
		Find(&song_list).
		Error
	return song_list, err
}

/**	通过ID检索到曲目
 *	@param	db	gorm数据库连接对象
 *	@param	song_id	检索目标曲目ID
 */
func SelectSongById(
	db *gorm.DB,
	song_id uint64,
) (SongTrackDetail, error) {
	var err error
	var exist_signal int
	var song SongTrackDetail

	err = db.Model(&localmodel.SongTable{}).
		Select("COUNT(*)").
		Where("id = ?", song_id).
		First(&exist_signal).Error

	if exist_signal == 0 {
		return SongTrackDetail{}, err
	} else {
		err = db.Model(&localmodel.SongTable{}).
			Select("song_table.id, song_table.name, song_table.ar_id, artist_table.name as 'ar_name' ,song_table.al_id, album_table.name as 'al_name', local_path, privilege, quality, useful, source").
			Joins("LEFT JOIN mblog.album_table ON song_table.al_id = album_table.id LEFT JOIN mblog.artist_table ON song_table.ar_id = artist_table.id").
			Where("song_table.id = ?", song_id).
			First(&song).
			Error
		return song, err
	}
}

/**	修改目标曲目的可用性、收听等级、来源
 *	@param	db	gorm数据库连接对象
 *	@param	song_id	修改目标曲目ID
 *	@param	privilege	权限
 *	@param	useful	可用性
 *	@param	source	音频来源
 */
func UpdateSongDetail(
	db *gorm.DB,
	song_id uint64,
	privilege int,
	useful int,
	source string,
) (SongTrackDetail, error) {
	var err error
	var exist_signal int
	var song SongTrackDetail

	err = db.Model(&localmodel.SongTable{}).
		Select("COUNT(*)").
		Where("id = ?", song_id).
		First(&exist_signal).Error

	if exist_signal == 0 {
		return SongTrackDetail{}, err
	} else {
		/* 修改曲目信息 */
		/* 修改收听等级权限 */
		err = db.Model(&localmodel.SongTable{}).
			Where("id = ?", song_id).
			Update("privilege", privilege).
			Error
		/* 修改音频可用性 */
		err = db.Model(&localmodel.SongTable{}).
			Where("id = ?", song_id).
			Update("useful", useful).
			Error
		/* 修改音频来源 */
		err = db.Model(&localmodel.SongTable{}).
			Where("id = ?", song_id).
			Update("source", source).
			Error

		/* 查询曲目信息 */
		err = db.Model(&localmodel.SongTable{}).
			Select("song_table.id, song_table.name, song_table.ar_id, artist_table.name as 'ar_name' ,song_table.al_id, album_table.name as 'al_name', local_path, privilege, quality, useful, source").
			Joins("LEFT JOIN mblog.album_table ON song_table.al_id = album_table.id LEFT JOIN mblog.artist_table ON song_table.ar_id = artist_table.id").
			Where("song_table.id = ?", song_id).
			First(&song).
			Error
		return song, err
	}
}
