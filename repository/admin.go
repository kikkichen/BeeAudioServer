package repository

import (
	"BeeAudioServer/models/localmodel"
	"BeeAudioServer/utils"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ServerTotal struct {
	/* 用户总数 */
	User uint64 `json:"user_count"`
	/* 博文发文量 */
	Blog uint64 `json:"blog_count"`
	/* 评论总量 */
	Comment uint64 `json:"comment_count"`
	/* 点赞总量 */
	Attitude uint64 `json:"attitude_count"`
	/* 音频数量 */
	Song uint64 `json:"song_count"`
	/* 专辑总量 */
	Album uint64 `json:"album_count"`
	/* 艺人总量 */
	Artist uint64 `json:"artist_count"`
	/* 歌单总量 */
	Playlist uint64 `json:"playlist_count"`
}

/**	创建新的管理员
 *	@param	db	gorm数据库连接对象
 *	@param	name	管理员昵称
 *	@param	password	密码明文
 */
func CreateAdministrator(db *gorm.DB, name, password string) (uint64, error) {
	var err error
	var new_admin_id uint64
	/* 循环生成Uid,并验证其是空闲的 */
	for {
		new_admin_id_tail := utils.RandomString(4, utils.DefaultNumber)
		new_admin_id_string := "80" + new_admin_id_tail
		new_admin_id = utils.StringParseToUint64(new_admin_id_string)
		/* 判断新生成的Uid没有被使用 */
		exist_signal, _ := IsExistInAdminModel(db, new_admin_id)
		if !exist_signal {
			break
		}
	}

	new_admin := localmodel.Administrator{
		Id: new_admin_id,
		Name: func(string) string {
			if len(name) == 0 {
				return fmt.Sprintf("管理员%v", new_admin_id)
			} else {
				return name
			}
		}(name),
		Password:  utils.GenerateStringByMD5(password),
		AdminType: 0,
	}

	/* 创建新管理员逻辑 */
	err = db.Model(&localmodel.Administrator{}).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&new_admin).
		Error
	if err == nil {
		return new_admin_id, nil
	} else {
		return 0, err
	}
}

/**	验证该管理员是否存在 - 通过管理员ID
 *	@param	db	gorm数据库连接对象
 *	@param	adminId	查询目标管理员ID
 */
func IsExistInAdminModel(db *gorm.DB, adminId uint64) (bool, error) {
	var err error
	var exist int = 0
	err = db.Model(&localmodel.Administrator{}).
		Select("COUNT(*)").
		Where("id = ?", adminId).
		First(&exist).Error
	if exist > 0 {
		return true, err
	} else {
		return false, err
	}
}

/**	修改管理员信息
 *	@param	db	gorm数据库连接对象
 *	@param	admin	目标修改管理员数据对象
 */
func UpdateAdminDetail(db *gorm.DB, admin localmodel.Administrator) error {
	var err error
	original_admin, err := FindAdministartorByID(db, admin.Id)

	/* 修改昵称 */
	err = db.Model(&localmodel.Administrator{}).
		Where("id = ?", admin.Id).
		Update("name", admin.Name).Error

	/* 修改管理员类型 */
	err = db.Model(&localmodel.Administrator{}).
		Where("id = ?", admin.Id).Update("type", admin.AdminType).Error

	/* 修改管理员密码 */
	if admin.Password != original_admin.Password {
		/* 如果密文有被更改，则重新加密 */
		err = db.Model(&localmodel.Administrator{}).
			Where("id = ?", admin.Id).
			Update("password", utils.GenerateStringByMD5(admin.Password)).Error
	}
	return err
}

/**	删除管理员
 *	@param	db	gorm数据库连接对象
 *	@param	adminId	查询目标管理员ID
 */
func DeleteTargetAdmin(db *gorm.DB, adminId uint64) error {
	err := db.Model(&localmodel.Administrator{}).
		Where("id = ?", adminId).
		Delete(&localmodel.Administrator{}).Error
	return err
}

/**	依据ID查询管理员
 *	@param	db	gorm数据库连接对象
 *	@param	adminId	查询目标管理员ID
 */
func FindAdministartorByID(db *gorm.DB, adminId uint64) (localmodel.Administrator, error) {
	var target_admin localmodel.Administrator
	err := db.Model(&localmodel.Administrator{}).
		Where("id = ?", adminId).
		First(&target_admin).Error
	return target_admin, err
}

/**	管理员登陆核验
 *	@param	db	gorm数据库连接对象
 *	@param	adminId	查询目标管理员ID
 *	@param	password	登陆密码
 */
func AdministratorLoginVerify(db *gorm.DB, adminId uint64, password string) (int, error) {
	var err error
	var select_admin localmodel.Administrator
	/* 判断管理员用户是否存在 */
	exist_signal, _ := IsExistInAdminModel(db, adminId)
	if !exist_signal {
		return -1, err
	}

	/* 核对密码 */
	err = db.Model(&localmodel.Administrator{}).
		Where("id = ?", adminId).
		First(&select_admin).Error
	if select_admin.Password == utils.GenerateStringByMD5(password) {
		return 0, err // 管理员用户存在， 且密码核对正确
	} else {
		return 1, err // 管理员用户存在， 但密码核对不正确
	}
}

/**	查询管理员 - 分页
 *	@params	db	gorm数据库连接对象
 *	@param	page	页码
 *	@param	size	每页大小容量
 */
func BrowserAdministorsByPage(db *gorm.DB, page, size int) ([]localmodel.Administrator, error) {
	var err error
	var admin_list []localmodel.Administrator
	err = db.Model(&localmodel.Administrator{}).
		Order("id DESC").
		Limit(size).
		Offset((page - 1) * size).
		Find(&admin_list).Error
	return admin_list, err
}

/**	查询当前管理员容量
 *	@params	db	gorm数据库连接对象
 */
func GetAdminTotalNumber(db *gorm.DB) (int, error) {
	var total int
	err := db.Model(&localmodel.Administrator{}).
		Select("COUNT(*)").
		First(&total).
		Error
	return total, err
}

/**	获取服务端总汇信息
 *	@params	db	gorm数据库连接对象
 */
func GetServerTotalInfo(db *gorm.DB) (ServerTotal, error) {
	var err error
	/* 用户总数 */
	var user_count uint64
	/* 博文发文量 */
	var blog_count uint64
	/* 评论总量 */
	var comment_count uint64
	/* 点赞总量 */
	var attitude_count uint64
	/* 音频数量 */
	var song_count uint64
	/* 专辑总量 */
	var album_count uint64
	/* 艺人总量 */
	var artist_count uint64
	/* 歌单总量 */
	var playlist_count uint64

	err = db.Model(&localmodel.UserModel{}).
		Select("COUNT(*)").
		First(&user_count).Error

	err = db.Model(&localmodel.BlogModel{}).
		Select("COUNT(*)").
		First(&blog_count).Error

	err = db.Model(&localmodel.CommentModel{}).
		Select("COUNT(*)").
		First(&comment_count).Error

	err = db.Model(&localmodel.AttitudeModel{}).
		Select("COUNT(*)").
		First(&attitude_count).Error

	err = db.Model(&localmodel.SongTable{}).
		Select("COUNT(*)").
		First(&song_count).Error

	err = db.Model(&localmodel.AlbumModel{}).
		Select("COUNT(*)").
		First(&album_count).Error

	err = db.Model(&localmodel.ArtistModel{}).
		Select("COUNT(*)").
		First(&artist_count).Error

	err = db.Model(&localmodel.UserPlayListModel{}).
		Select("COUNT(*)").
		First(&playlist_count).Error

	if err == nil {
		return ServerTotal{
			User:     user_count,
			Blog:     blog_count,
			Comment:  comment_count,
			Attitude: attitude_count,
			Song:     song_count,
			Album:    album_count,
			Artist:   artist_count,
			Playlist: playlist_count,
		}, nil
	} else {
		return ServerTotal{}, err
	}
}
