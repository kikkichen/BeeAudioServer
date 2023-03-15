package repository

import (
	"BeeAudioServer/models/localmodel"
	"BeeAudioServer/utils"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*	为目标博文点赞
 *	@param	db	gorm 连接对象
 *	@param	bid	目标动态博文 bid
 *	@param	uid	当前触发请求用户的 uid
 *	@param	source	点赞来源设备信息
 *
 *	返回值中的布尔类型，true表示完成“点赞”操作， false 表示完成 “取消点赞”操作
 */
func ActionAttitude(
	db *gorm.DB,
	bid uint64,
	uid uint64,
	source string,
) (bool, error) {
	if IsPullAttitude(db, bid, uid) {
		return false, db.Transaction(func(tx *gorm.DB) error {
			/* 若已经点赞，则再执行该方法表示撤销点赞 */
			attitude_record := GetPullAttitude(tx, bid, uid)
			/* 删除原有的点赞记录信息 */
			if err := tx.Delete(&attitude_record).Error; err != nil {
				return err
			}
			/* 更新原博文动态的点赞数量值 */
			/* 获取原博文的点赞数据 */
			orignal_blog_info := GetBlogDetail(db, attitude_record.Bid)
			/* 更新博文表中的点赞数 -1 */
			return tx.Model(&localmodel.BlogModel{}).
				Where("bid = ?", orignal_blog_info.Bid).
				Update("attitudes_count", orignal_blog_info.AttitudesCount-1).
				Error
		})

	} else {
		/* 若没有点赞，则新增一条点赞记录 */
		/* 随机生成 Attitude Id */
		var new_aid_tail string
		var new_aid_string string
		var new_aid uint64

		/* 循环生成 Aid,并验证其是空闲的 */
		for {
			new_aid_tail = utils.RandomString(10, utils.DefaultNumber)
			new_aid_string = "890000" + new_aid_tail
			new_aid = utils.StringParseToUint64(new_aid_string)
			/* 判断新生成的Bid没有被使用 */
			if !IsPullAttitudeByAid(db, new_aid) {
				break
			}
		}

		new_attitude := localmodel.AttitudeModel{
			Aid:     new_aid,
			Bid:     bid,
			Uid:     uid,
			Created: time.Now(),
			Source:  source,
		}

		/* 触发点赞 事务逻辑 */
		return true, db.Transaction(func(tx *gorm.DB) error {
			/* 插入记录 */
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&new_attitude).Error; err != nil {
				return err
			}
			/* 获取原博文的点赞数据 */
			orignal_blog_info := GetBlogDetail(db, bid)
			/* 更新博文表中的点赞数 +1 */
			return tx.Model(&localmodel.BlogModel{}).
				Where("bid = ?", orignal_blog_info.Bid).
				Update("attitudes_count", orignal_blog_info.AttitudesCount+1).
				Error
		})
	}
}

/*	检测是否对当前博文动态点赞
 *	@param	db	gorm 连接对象
 *	@param	bid	目标动态博文 bid
 *	@param	uid	当前触发请求用户的 uid
 *
 *	若当前请求用户已对当前博文动态点赞，返回值为true; 若未对当前博文动态点赞，返回值为false
 */
func IsPullAttitude(
	db *gorm.DB,
	bid uint64,
	uid uint64,
) bool {
	var attitude_record localmodel.AttitudeModel
	db.Model(&localmodel.AttitudeModel{}).
		Where("uid = ? AND bid = ?", uid, bid).
		Find(&attitude_record)

	if attitude_record.Aid == 0 {
		return false
	} else {
		return true
	}
}

/*	依据Aid,检测是否存在对当前博文动态点赞的记录
 *	@param	db	gorm 连接对象
 *	@param	aid	点赞记录 Aid
 *
 *	若当前请求用户已对当前博文动态点赞，返回值为true; 若未对当前博文动态点赞，返回值为false
 */
func IsPullAttitudeByAid(
	db *gorm.DB,
	aid uint64,
) bool {
	var attitude_record localmodel.AttitudeModel
	db.Model(&localmodel.AttitudeModel{}).
		Where("aid = ?", aid).
		Find(&attitude_record)

	if attitude_record.Aid == 0 {
		return false
	} else {
		return true
	}
}

/*	获取评论记录信息
 *	@param	db	gorm 连接对象
 *	@param	bid	目标动态博文 bid
 *	@param	uid	当前触发请求用户的 uid
 */
func GetPullAttitude(
	db *gorm.DB,
	bid uint64,
	uid uint64,
) localmodel.AttitudeModel {
	var attitude_record localmodel.AttitudeModel
	db.Model(&localmodel.AttitudeModel{}).
		Where("uid = ? AND bid = ?", uid, bid).
		Find(&attitude_record)
	return attitude_record
}

/*
 *
 *
 */
