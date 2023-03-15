package repository

import (
	"BeeAudioServer/models/localmodel"

	"gorm.io/gorm"
)

/* 处理用户之间的订阅（关注）关系 */

/*	关注指定Uid用户
 *	@param db	gorm链接对象
 *	@param	uid 用户UID
 *	@param	be_follow_uid	关注对象Uid
 *
 *	返回布尔值为 true 表示处理“关注”关系， false 表示处理“解除关注”关系
 */
func ActionFollowUser(
	db *gorm.DB,
	uid uint64,
	be_follow_uid uint64,
) (bool, error) {
	/* 查询有无记录 */
	var follow_record_count []localmodel.FollowModel
	db.Model(&localmodel.FollowModel{}).
		Where("follow_uid = ? AND be_follow_uid = ?", uid, be_follow_uid).
		Find(&follow_record_count)

	if len(follow_record_count) == 0 {
		/* 如果没有记录，则未建立起订阅（关注）关系，执行关注流传 */
		follow := localmodel.FollowModel{
			FollowUid:   uid,
			BeFollowUid: be_follow_uid,
		}
		if err := db.Model(&localmodel.FollowModel{}).Create(&follow).Error; err != nil {
			return true, err
		}
		return true, nil
	} else {
		/* 已存在关注关系， 则解除该关注关系 */
		if err := db.Delete(&localmodel.FollowModel{}, "follow_uid = ? AND be_follow_uid = ?", uid, be_follow_uid).Error; err != nil {
			return false, err
		}
		/* 返回 nil 空对象， 执行事务 */
		return false, nil
	}
}

/*	查询关注列表
 *	@params db	gorm链接对象
 *  @params	uid	查询目标uid用户的关注列表
 */
func MySubscribe(
	db *gorm.DB,
	uid uint64,
) []localmodel.UserModel {
	var subscribe_uid_list []localmodel.FollowModel
	var be_subscribe_user_list []localmodel.UserModel
	/* 查询目标 uid 关注订阅的用户 */
	db.Find(&subscribe_uid_list, "follow_uid = ?", uid)

	/* 获取关注的uid */
	var uids []uint64
	for _, item := range subscribe_uid_list {
		uids = append(uids, item.BeFollowUid)
	}

	/* 关注用户信息查询 */
	db.Table("user_table").Where("uid IN ?", uids).Find(&be_subscribe_user_list)

	return be_subscribe_user_list
}

/*	获取粉丝列表
 *	@params db	gorm链接对象
 *  @params	uid	查询目标uid用户的关注列表
 */
func MyFans(
	db *gorm.DB,
	uid uint64,
) []localmodel.UserModel {
	var subscribe_uid_list []localmodel.FollowModel
	var be_subscribe_user_list []localmodel.UserModel
	/* 查询订阅 uid 目标的用户 */
	db.Find(&subscribe_uid_list, "be_follow_uid = ?", uid)

	/* 获取粉丝的uid */
	var uids []uint64
	for _, item := range subscribe_uid_list {
		uids = append(uids, item.FollowUid)
	}

	/* 关注用户信息查询 */
	db.Table("user_table").Where("uid IN ?", uids).Find(&be_subscribe_user_list)

	return be_subscribe_user_list
}

/*	获取互粉列表
 *	@params db	gorm链接对象
 *  @params	uid	查询目标uid用户的关注列表
 */
func MyFriends(
	db *gorm.DB,
	uid uint64,
) []localmodel.UserModel {
	/* 查询我的关注 */
	var my_subscribe_uids []localmodel.FollowModel
	db.Find(&my_subscribe_uids, "follow_uid = ?", uid)
	/* 查询关注我的 */
	var my_fans_uids []localmodel.FollowModel
	db.Find(&my_fans_uids, "be_follow_uid = ?", uid)

	/* 查询我的关注里有无我的粉丝 */
	var my_friends_uid []uint64
	for _, fan := range my_fans_uids {
		for _, focus := range my_subscribe_uids {
			if focus.BeFollowUid == fan.FollowUid {
				my_friends_uid = append(my_friends_uid, focus.BeFollowUid)
			}
		}
	}

	/* 为我的互粉列表填充信息 */
	var my_friends []localmodel.UserModel
	db.Table("user_table").Where("uid IN ?", my_friends_uid).Find(&my_friends)
	return my_friends
}

/*	获取互粉列表 V2 采用InnoDB表连接方法
 *	@params db	gorm链接对象
 *  @params	uid	查询目标uid用户的关注列表
 *
 *	(SELECT be_follow_uid FROM mblog.follow_model fm  left join mblog.user_model um on fm.be_follow_uid = um.uid and fm.be_follow_uid = um.uid where fm.follow_uid = 2448497652)
 *	SELECT follow_uid, name, description, avatar_url, be_follow_uid  FROM mblog.follow_model fm LEFT JOIN mblog.user_model um2 ON fm.follow_uid = um2.uid where be_follow_uid = 2448497652 and follow_uid IN  (?)
 */
func MyFriendsV2(
	db *gorm.DB,
	uid uint64,
) []localmodel.UserModel {
	var my_focus_uids []uint64
	/* 我的关注 */
	db.Model(&localmodel.FollowModel{}).
		Select("be_follow_uid").
		Joins("left join mblog.user_table um on follow_table.be_follow_uid = um.uid and follow_table.be_follow_uid = um.uid").
		Where("follow_table.follow_uid = ?", uid).
		Scan(&my_focus_uids)

	/* 在我的粉丝用户联表列表中筛选 我的关注 */
	var my_friends []localmodel.UserModel
	db.Model(&localmodel.FollowModel{}).
		Select("uid, name, description, avatar_url, be_follow_uid").
		Joins("LEFT JOIN mblog.user_table um2 ON follow_table.follow_uid = um2.uid").
		Where("be_follow_uid = ? and follow_uid IN  (?)", uid, my_focus_uids).
		Scan(&my_friends)
	return my_friends
}
