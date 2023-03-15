package repository

import (
	"BeeAudioServer/models/localmodel"
	"BeeAudioServer/utils"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*	获取用户的套餐信息
 *	@params	db		gorm 目标连接对象
 *	@params	uid		目标查询用户的Uid
 *
 *	返回 用户的最新一次Premiumn订单对象
 */
func GetUserPremiumInfo(
	db *gorm.DB,
	uid uint64,
) localmodel.PremiumModel {
	var premium_info localmodel.PremiumModel
	db.Model(&localmodel.PremiumModel{}).
		Where("uid = ?", uid).
		Order("server_in DESC").
		Last(&premium_info)

	return premium_info
}

/**	通过cardID获取Premium套餐信息
 *	@params	db		gorm 目标连接对象
 *	@params	card_id	加入Premium套餐订单卡号ID
 */
func GetPremiumOrderInfo(
	db *gorm.DB,
	card_id string,
) localmodel.PremiumModel {
	var premium_info localmodel.PremiumModel
	db.Model(&localmodel.PremiumModel{}).
		Where("card_id = ?", card_id).
		Order("server_in DESC").
		Last(&premium_info)

	return premium_info
}

/*	获取用户的近20次套餐信息
 *	@params	db		gorm 目标连接对象
 *	@params	uid		目标查询用户的Uid
 *
 *	返回 用户的近 20次 Premiumn订单对象
 */
func GetUserPremiumInfoList(
	db *gorm.DB,
	uid uint64,
) []localmodel.PremiumModel {
	var premium_info_list []localmodel.PremiumModel
	db.Model(&localmodel.PremiumModel{}).
		Where("uid = ?", uid).
		Order("server_in DESC").
		Limit(20).
		Find(&premium_info_list)

	return premium_info_list
}

/*	过渡用户类型
 *	@params	db		gorm 目标连接对象
 *	@params	card_id		Premium 订单卡号
 *	@params	isPremium	判断Premium过度类型, true为由普通用户过度到Premium会员
 *	@params	card_type	Premium 套餐类型
 *	@params	uid		执行更新的目标用户Uid
 */
func TransitionUserType(
	db *gorm.DB,
	card_id string,
	isPremium bool,
	card_type int,
	uid uint64,
) bool {
	/* 由普通用户过渡到Premium会员 */
	if isPremium {
		err := db.Transaction(func(tx *gorm.DB) error {
			/* 注册一个新的 Premium 订单 */
			new_premium := localmodel.PremiumModel{
				Uid:      uid,
				CardId:   card_id,
				CardType: card_type,
				ServerIn: time.Now(),
				ServerExpired: func(t int) time.Time {
					if card_type == 0 {
						/* 若是个人Premium套餐，则服务期限为1年 */
						return time.Now().Add(365 * 24 * time.Hour)
					} else {
						var count int = 0
						db.Model(&localmodel.PremiumModel{}).
							Select("COUNT(*)").
							Where("card_id = ?", card_id).
							Order("server_in ASC").
							Last(&count)
						if count == 0 {
							/* 若是当前用户是家庭组的第一位用户，则服务期限为1年 */
							return time.Now().Add(365 * 24 * time.Hour)
						} else {
							/* 若是要加入家庭组套餐， 以第一位家庭成员的服务结束日期为准 */
							var expired_time time.Time
							db.Model(&localmodel.PremiumModel{}).
								Select("server_expired").
								Where("card_id = ?", card_id).
								Order("server_in ASC").
								Last(&expired_time)
							return expired_time
						}
					}
				}(card_type),
			}
			/* 实现插入Premium卡表 */
			if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&new_premium).Error; err != nil {
				return err
			}
			/* Premium 订单注册成功后 更新用户类型 */
			tx.Model(&localmodel.UserPermissionModel{}).
				Where("uid = ?", uid).
				Update("user_type", 1)
			return nil
		})
		if err != nil {
			return false
		}
	} else {
		/* 由Premium会员变回普通会员 */
		result := db.Model(&localmodel.UserPermissionModel{}).
			Where("uid = ?", uid).
			Update("user_type", 0)
		if result.Error != nil {
			return false
		}
	}

	return true
}

/*	生成 Premium 卡号
 *	@params	db	gorm 数据库连接对象
 *
 */
func GeneratePremiumCardId(db *gorm.DB) string {
	var card_id string

	for {
		/* 随即生成25位数大写字母兼数字卡号 */
		card_id = utils.RandomString(25, utils.DefaultUpperLetters)
		/* ID重复查询 */
		result := db.Model(&localmodel.PremiumModel{}).
			Where("card_id = ?", card_id).
			Order("server_in DESC ").
			First(&localmodel.PremiumModel{})
		if result.RowsAffected == 0 {
			break
		}
	}

	/* 返回随机生成的卡号 */
	return card_id
}

/*	验证是否存在卡号
 *	@params	CardId	25位卡号 Id 字符串
 *	@params	db	gorm 数据库连接对象
 */
func IsExistPermiumCard(
	db *gorm.DB,
	card_id string,
) bool {
	result := db.Model(&localmodel.PremiumModel{}).
		Where("card_id = ?", card_id).
		Order("server_in DESC ").
		First(&localmodel.PremiumModel{})
	if result.RowsAffected == 0 {
		return false
	} else {
		return true
	}
}

/*	验证套餐是否满载
 *	@params	CardId	25位卡号 Id 字符串
 *	@params	db	gorm 数据库连接对象
 *
 *	返回值为 true 表示未满载， false表示满载
 */
func IsOverLoaderCardID(
	db *gorm.DB,
	card_id string,
) bool {
	if GetCardType(db, card_id) == 1 {
		var counter int64
		db.Model(&localmodel.PremiumModel{}).
			Where("card_id = ? And card_type = ?", card_id, 0).
			Count(&counter)
		if counter < 6 {
			return true
		} else {
			/* 用户满载 */
			return false
		}
	} else {
		/* 若类型为非家庭组类型，存在的情况下则为满载 */
		return false
	}
}

/*	验证Premium套餐类型
 *	@params	db	gorm 数据库连接对象
 *	@params	card_id	卡号ID
 *
 *	类型 0 为Premium个人套餐, 类型Premium 1 为家庭套餐
 */
func GetCardType(
	db *gorm.DB,
	card_id string,
) int {
	var card_type int
	db.Model(&localmodel.PremiumModel{}).
		Select("card_type").
		Where("card_id = ?", card_id).
		Order("server_in ASC ").
		Last(&card_type)

	return card_type
}

/*	验证Premium套餐 是否过期
 *	@params	db	gorm 数据库连接对象
 *	@params	card_id	卡号ID
 *
 *	返回值为true代表没过期， 返回值为 false 代表已经过期
 */
func isServerTimePremiumn(
	db *gorm.DB,
	card_id string,
) bool {
	var premium_server_expired time.Time
	db.Model(&localmodel.PremiumModel{}).
		Select("server_expired").
		Where("card_id = ?", card_id).
		Order("server_in ASC ").
		Last(&premium_server_expired)

	if premium_server_expired.Unix() > time.Now().Unix() {
		return true
	} else {
		return false
	}
}

/*	Premium 过期执行 (每次上线执行判断)
 *	@params	db	gorm连接对象
 *	@params	uid	目标用户Uid
 *
 */
func DealWithExpiredPremium(
	db *gorm.DB,
	uid uint64,
) {
	/* 验证是否为 Premium 会员 */
	var permium_signal int
	db.Model(&localmodel.UserPermissionModel{}).
		Select("user_type").
		Where("uid = ?", uid).
		Scan(&permium_signal)

	/* 若非 Premium 会员，则跳过下列执行 */
	if permium_signal == 0 {
		return
	} else {
		premium_info := GetUserPremiumInfo(db, uid)
		/* 过期机制判断 */
		expired_signal := isServerTimePremiumn(db, premium_info.CardId)
		if !expired_signal {
			return
		} else {
			/* 若过期执行用户类型过渡 */
			TransitionUserType(db, premium_info.CardId, false, premium_info.CardType, uid)
		}
	}
}

/*	成为 个人 Premium会员
 *	@params	db	gorm连接对象
 *	@params	uid	目标用户Uid
 *
 */
func BecomePersonPremium(
	db *gorm.DB,
	uid uint64,
) string {
	/* 生成新的 Premium 卡号ID */
	new_card_id := GeneratePremiumCardId(db)
	TransitionUserType(db, new_card_id, true, 0, uid)
	return new_card_id
}

/*	成为 家庭组 Premium会员 - 成为该家庭组的首位用户，即管理员
 *	@params	db gorm连接对象
 *	@params	uid 目标用户Uid
 */
func BecomeFamilyPremium(
	db *gorm.DB,
	uid uint64,
) string {
	new_card_id := GeneratePremiumCardId(db)
	TransitionUserType(db, new_card_id, true, 1, uid)
	return new_card_id
}

/*	加入 家庭组， 成功Premium会员 - 管理员同意后的操作
 *	@params	db gorm连接对象
 *	@params	uid 目标用户Uid
 *	@params	card_id	加入Premium家庭组的卡号ID
 */
func JoinFamilyPremium(
	db *gorm.DB,
	uid uint64,
	target_user_id uint64,
	card_id string,
) bool {
	/* 判断其卡号存在， 且家庭组不为满载状态 */
	if IsExistPermiumCard(db, card_id) && IsOverLoaderCardID(db, card_id) {
		/* 确保当前操作是由该家庭组的管理员执行 */
		var premium_info localmodel.PremiumModel
		db.Model(&localmodel.PremiumModel{}).
			Where("card_id = ?", card_id).
			Order("server_in ASC").
			First(&premium_info)

		if premium_info.Uid == uid {
			return TransitionUserType(db, card_id, true, 1, target_user_id)
		} else {
			return false
		}
	} else {
		/* 家庭组不存在， 或者家庭组满载 */
		return false
	}
}

/**	提交加入家庭去的申请
 *	@params	db gorm连接对象
 *	@params	uid 目标用户Uid
 *	@params	card_id	加入Premium家庭组的卡号ID
 */
func PostJoinFamilyPremiumApply(
	db *gorm.DB,
	uid uint64,
	card_id string,
) bool {
	/* 判断其卡号存在， 且家庭组不为满载状态 */
	if IsExistPermiumCard(db, card_id) && IsOverLoaderCardID(db, card_id) {
		err := db.Model(&localmodel.PremiumModel{}).
			Clauses(clause.OnConflict{UpdateAll: true}).
			Create(&localmodel.PremiumModel{
				Uid:           uid,
				CardId:        card_id,
				CardType:      10,
				ServerIn:      time.Now(),
				ServerExpired: time.Now(),
			}).Error
		if err == nil {
			return true
		} else {
			return false
		}
	} else {
		/* 家庭组不存在， 或者家庭组满载 */
		return false
	}
}

/**	获取Premium家庭组中的成员列表
 *	@params	db gorm连接对象
 *	@params	card_id	加入Premium家庭组的卡号ID
 */
func GetPremiumFamilyUserGroup(
	db *gorm.DB,
	card_id string,
) ([]localmodel.UserModel, bool) {
	/* 判断其卡号存在 */
	if IsExistPermiumCard(db, card_id) {
		/* 获取Premium家庭组内成员的订单信息 */
		var premium_info_list []localmodel.PremiumModel
		db.Model(&localmodel.PremiumModel{}).
			Where("card_id = ? And card_type = ?", card_id, 1).
			Order("server_in ASC").
			Find(&premium_info_list)
		/* 检索用户信息 */
		var family_group []localmodel.UserModel
		for _, premium_order := range premium_info_list {
			var temp_user localmodel.UserModel
			db.Model(&localmodel.UserModel{}).
				Where("uid = ?", premium_order.Uid).
				First(&temp_user)
			family_group = append(family_group, temp_user)
		}
		return family_group, true
	} else {
		/* 家庭组不存在 */
		return []localmodel.UserModel{}, false
	}
}

/**	获取Premium家庭组 申请加入成员列表
 *	@params	db gorm连接对象
 *	@params	card_id	加入Premium家庭组的卡号ID
 */
func GetPremiumFamilyPostApplyGroup(
	db *gorm.DB,
	card_id string,
) ([]localmodel.UserModel, bool) {
	/* 判断其卡号存在 */
	if IsExistPermiumCard(db, card_id) {
		/* 获取Premium家庭组内成员的订单信息 */
		var premium_info_list []localmodel.PremiumModel
		db.Model(&localmodel.PremiumModel{}).
			Where("card_id = ? And card_type = ?", card_id, 10).
			Order("server_in ASC").
			Find(&premium_info_list)
		/* 检索用户信息 */
		var family_group []localmodel.UserModel
		for _, premium_order := range premium_info_list {
			var temp_user localmodel.UserModel
			db.Model(&localmodel.UserModel{}).
				Where("uid = ?", premium_order.Uid).
				First(&temp_user)
			family_group = append(family_group, temp_user)
		}
		return family_group, true
	} else {
		/* 家庭组不存在 */
		return []localmodel.UserModel{}, false
	}
}

/**	管理员移除家庭组成员 - （包括移除审核申请）
 *	@params	db gorm连接对象
 *	@params	uid 目标用户Uid
 *	@params	card_id	加入Premium家庭组的卡号ID
 */
func RemoveUserFromPremiumFamilyGroup(
	db *gorm.DB,
	uid uint64,
	card_id string,
) error {
	var err error
	err = db.Where("card_id = ? AND uid = ?", card_id, uid).
		Delete(&localmodel.PremiumModel{}).Error
	err = db.Model(&localmodel.UserPermissionModel{}).
		Where("uid = ?", uid).
		Update("user_type", 0).Error
	return err
}
