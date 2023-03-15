package repository

import (
	"BeeAudioServer/models/localmodel"
	"BeeAudioServer/utils"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*	依据用户Uid，修改个人 邮箱 信息
 *	@params	db		gorm 目标连接对象
 *	@params	email	即将更新的邮箱字符串
 *	@params	uid		执行更新的目标用户Uid
 *
 */
func UpdateUserEmail(
	db *gorm.DB,
	email string,
	uid uint64,
) bool {
	action_result := db.Model(localmodel.UserPermissionModel{}).
		Where("uid = ?", uid).
		Update("email", email)

	if action_result.Error != nil {
		return false
	}

	return true
}

/*	依据用户Uid，修改 电话号码 信息
 *	@params	db		gorm 目标连接对象
 *	@params	phone	即将更新的电话号码字符串
 *	@params	uid		执行更新的目标用户Uid
 *
 */
func UpdateUserPhone(
	db *gorm.DB,
	phone string,
	uid uint64,
) bool {
	action_result := db.Model(localmodel.UserPermissionModel{}).
		Where("uid = ?", uid).
		Update("phone", phone)

	if action_result.Error != nil {
		return false
	}

	return true
}

/*	依据用户Uid，修改 电话号码 信息
 *	@params	db		gorm 目标连接对象
 *	@params	phone	即将更新的电话号码字符串
 *	@params	uid		执行更新的目标用户Uid
 *
 */
func UpdateUserBirthday(
	db *gorm.DB,
	birthday time.Time,
	uid uint64,
) bool {
	action_result := db.Model(localmodel.UserPermissionModel{}).
		Where("uid = ?", uid).
		Update("birthday", birthday)

	if action_result.Error != nil {
		return false
	}

	return true
}

/*	密码修改
 *	@params	db		gorm 目标连接对象
 *	@params	password	即将修改的新密码
 *	@params	uid		执行更新的目标用户Uid
 */
func UpdateUserPassword(
	db *gorm.DB,
	password string,
	uid uint64,
) bool {
	action_result := db.Model(localmodel.UserPermissionModel{}).
		Where("uid = ?", uid).
		Update("password", utils.GenerateStringByMD5(password))

	if action_result.Error != nil {
		return false
	}

	return true
}

/*	修改用户名
 *	@params	db		gorm 目标连接对象
 *	@params	name	即将修改的用户名
 *	@params	uid		执行更新的目标用户Uid
 */
func UpdateUserName(
	db *gorm.DB,
	name string,
	uid uint64,
) bool {
	action_result := db.Model(localmodel.UserModel{}).
		Where("uid = ?", uid).
		Update("name", name)

	if action_result.Error != nil {
		return false
	}

	return true
}

/*	修改用户头像路径
 *	@params	db		gorm 目标连接对象
 *	@params	avatar_url	修改后的头像信息地址
 *	@params	uid		执行更新的目标用户Uid
 */
func UpdateUserAvatarPath(
	db *gorm.DB,
	avatar_url string,
	uid uint64,
) bool {
	action_result := db.Model(localmodel.UserModel{}).
		Where("uid = ?", uid).
		Update("avatar_url", avatar_url)

	if action_result.Error != nil {
		return false
	}

	return true
}

/*	修改用户描述文字
 *	@params	db		gorm 目标连接对象
 *	@params	description	用户简介字符串字段
 *	@params	uid		执行更新的目标用户Uid
 */
func UpdateUserDescription(
	db *gorm.DB,
	description string,
	uid uint64,
) bool {
	action_result := db.Model(localmodel.UserModel{}).
		Where("uid = ?", uid).
		Update("description", description)

	if action_result.Error != nil {
		return false
	}

	return true
}

/**	管理员更改用户信息
 *	@param	db		gorm 目标连接对象
 *	@param	uid		修改目标用户ID
 *	@param	name	用户昵称
 *	@param	description	用户简介
 *	@param	email	邮件
 *	@param	phone	电话号码
 */
func ModifierUserDetailByAdmin(
	db *gorm.DB,
	uid uint64,
	name string,
	description string,
	email string,
	phone string,
) (UserInfoWithoutPassword, error) {
	var err error
	err = db.Model(&localmodel.UserModel{}).
		Where("uid = ?", uid).Error
	err = db.Model(&localmodel.UserModel{}).
		Where("uid = ?", uid).
		Update("name", name).Error
	err = db.Model(&localmodel.UserModel{}).
		Where("uid = ?", uid).
		Update("description", description).
		Error

	err = db.Model(&localmodel.UserPermissionModel{}).
		Where("uid = ?", uid).
		Update("email", email).Error
	err = db.Model(&localmodel.UserPermissionModel{}).
		Where("uid = ?", uid).
		Update("phone", phone).
		Error

	/* 查询修改后结果 */
	result_user := GetUserInfoWithoutPasword(db, uid)
	return result_user, err
}

/**	管理员注销用户
 *	@param	db		gorm 目标连接对象
 *	@param	uid		修改目标用户ID
 */
func LogoutUserByAdmin(
	db *gorm.DB,
	uid uint64,
) error {
	var err error
	/* 修改该用户的信息，修改为已注销状态 */
	if err = db.Model(&localmodel.UserPermissionModel{}).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&localmodel.UserPermissionModel{
			Uid:      uid,
			Password: "",
			UserType: -1,
			Email:    "",
			Phone:    "",
			Birthday: time.Now(),
		}).Error; err != nil {
		return err
	}

	/* 删除用户表记录 */
	if err = db.Model(&localmodel.UserModel{}).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&localmodel.UserModel{
			Uid:         uid,
			Name:        "用户已注销",
			Description: "",
			AvatarUrl:   "",
			CreatedAt:   time.Now(),
		}).Error; err != nil {
		return err
	}

	/* 删除关注关系 */
	db.Model(&localmodel.FollowModel{}).
		Where("follow_uid = ?", uid).
		Delete(&localmodel.FollowModel{})
	db.Model(&localmodel.FollowModel{}).
		Where("be_follow_uid = ?", uid).
		Delete(&localmodel.FollowModel{})

	/* 删除该用户全部的歌单 */
	db.Model(&localmodel.PlayListTable{}).
		Where("uid = ?", uid).
		Delete(&localmodel.PlayListTable{})
	db.Model(&localmodel.UserPlayListModel{}).
		Where("uid = ?", uid).
		Delete(&localmodel.UserPlayListModel{})
	db.Model(&localmodel.UserFavoritePlaylist{}).
		Where("uid = ?", uid).
		Delete(&localmodel.UserFavoritePlaylist{})

	/* 删除该用户的历史播放记录 */
	db.Model(&localmodel.HistoryDataModel{}).
		Where("uid = ?", uid).
		Delete(&localmodel.HistoryDataModel{})
	/* 删除该用户的订阅数据 */
	db.Model(&localmodel.SubscribeModel{}).
		Where("uid = ?", uid).
		Delete(&localmodel.SubscribeModel{})

	return err
}
