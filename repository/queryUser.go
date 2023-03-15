package repository

import (
	"BeeAudioServer/models/localmodel"
	"BeeAudioServer/utils"
	"time"

	"gorm.io/gorm"
)

/* 携带密码与用户类型信息的 用户类型 */
type UserInfoWithTypeAndPassword struct {
	Uid         uint64    `gorm:"column:uid;primaryKey" json:"uid"`
	Name        string    `gorm:"column:name" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	AvatarUrl   string    `gorm:"column:avatar_url" json:"avatar_url"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	Password    string    `gorm:"column:password" json:"password"`
	UserType    int       `gorm:"column:user_type" json:"user_type"`
}

type UserInfoWithoutPassword struct {
	Uid         uint64    `gorm:"column:uid;primaryKey" json:"uid"`
	Name        string    `gorm:"column:name" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	AvatarUrl   string    `gorm:"column:avatar_url" json:"avatar_url"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UserType    int       `gorm:"column:user_type" json:"user_type"`
	Email       string    `gorm:"column:email" json:"email"`
	Phone       string    `gorm:"column:phone" json:"phone"`
	Birthday    time.Time `gorm:"column:birthday" json:"birthday"`
}

/*	获取用户信息
 *	@params db	gorm链接对象
 *	@params	uid	目标查询用户uid
 */
func GetUserInfo(
	db *gorm.DB,
	uid uint64,
) localmodel.UserModel {
	var userinfo localmodel.UserModel
	db.Model(&localmodel.UserModel{}).
		Where("uid = ?", uid).
		Find(&userinfo)
	return userinfo
}

/*	获取用户信息 包含用户类型 （是否会员）
 *	@params db	gorm链接对象
 *	@params	uid	目标查询用户uid
 */
func GetUserInfoV2(
	db *gorm.DB,
	uid uint64,
) UserInfoWithTypeAndPassword {
	var user_info UserInfoWithTypeAndPassword
	db.Model(&localmodel.UserModel{}).
		Select("user_table.uid, name, description, avatar_url, created_at, user_type").
		Joins("LEFT JOIN user_permission_table ON user_table.uid = user_permission_table.uid").
		Where("user_table.uid = ?", uid).
		First(&user_info)
	return user_info
}

/*	关键字查找相关用户
 *	@params db			gorm链接对象
 *	@params	keyWord		关键字
 *	@params	page		请求页数
 *	@params	size		请求单页大小
 */
func SearchUserByKeyWord(
	db *gorm.DB,
	keyWord string,
	page int,
	size int,
) []localmodel.UserModel {
	var result_users []localmodel.UserModel
	db.Model(&localmodel.UserModel{}).
		Where("uid = ? OR name LIKE ?", utils.StringParseToUint64(keyWord), "%"+keyWord+"%").
		Limit(size).
		Offset((page - 1) * size).
		Find(&result_users)
	return result_users
}

/*	依据 uid判断用户是否存在
 *	@params db	gorm链接对象
 *	@params	uid 用户UID
 *
 *	若 uid 存在，返回true, 反之返回false
 */
func IsExistInUserModel(
	db *gorm.DB,
	uid uint64,
) bool {
	var result []localmodel.UserModel
	db.Model(&localmodel.UserModel{}).Where("uid = ?", uid).Find(&result)
	if len(result) > 0 {
		return true
	} else {
		return false
	}
}

/*	获取用户密钥
 *	@params	db	gorm链接对象
 *	@params	uid	目标查询用户UID
 *
 *	若查找不到该用户，则返回空字符串
 */
func GetUserPasswordString(
	db *gorm.DB,
	uid uint64,
) string {
	/* 声明一个 携带用户类型与用户密码的变量 */
	var user_info UserInfoWithTypeAndPassword
	result := db.Model(&localmodel.UserModel{}).
		Select("user_table.uid, name, description, avatar_url, created_at, password, user_type").
		Joins("LEFT JOIN user_permission_table ON user_table.uid = user_permission_table.uid").
		Where("user_table.uid = ?", uid).
		First(&user_info)

	/* 如果查找不到用户 */
	if result.Error != nil {
		return ""
	}

	/* 查找到有效用户则返回密码密文 */
	return user_info.Password
}

/*	依据用户Uid, 获取不携带用户密钥的 用户信息
 *	包含 用户Uid、 用户名、 用户描述、 用户头像地址、 用户创建时间、 用户类型、用户邮箱、用户电话号码、生日日期(日期格式可能出现问题)
 *	@param	db	gorm连接对象
 *	@param	uid	请求目标用户Uid
 *
 */
func GetUserInfoWithoutPasword(
	db *gorm.DB,
	uid uint64,
) UserInfoWithoutPassword {
	var user_info UserInfoWithoutPassword
	db.Model(&localmodel.UserModel{}).
		Select("user_table.uid, user_table.name, user_table.description, user_table.avatar_url, user_table.created_at, user_permission_table.user_type, user_permission_table.email, user_permission_table.phone, birthday").
		Joins("LEFT JOIN mblog.user_permission_table ON user_table.uid = user_permission_table.uid").
		Where("user_permission_table.uid = ?", uid).
		Scan(&user_info)
	return user_info
}

/*	依据油壶Uid,获取其关注数、粉丝数 & 互粉数
 *	@param	db	gorm连接对象
*	@param	uid	请求目标用户Uid
*/
func GetCountOfFollowAndFans(
	db *gorm.DB,
	uid uint64,
) (int, int, int) {
	var follows, fans, friends int
	db.Model(&localmodel.FollowModel{}).
		Select("COUNT(*)").
		Where("follow_uid = ?", uid).
		First(&follows)

	db.Model(&localmodel.FollowModel{}).
		Select("COUNT(*)").
		Where("be_follow_uid = ?", uid).
		First(&fans)

	responseFriends := MyFriendsV2(db, uid)
	friends = len(responseFriends)

	return follows, fans, friends
}

/**	浏览用户目录 - 分页
 *	@params	db	gorm数据库连接对象
 *	@param	page	页码
 *	@param	size	每页大小容量
*	@param	sort	排序方式 true为由小到大， false为由大到小
*/
func BrowserAllUser(db *gorm.DB, page, size int, sort bool) ([]UserInfoWithoutPassword, error) {
	var err error
	var user_list []UserInfoWithoutPassword
	err = db.Model(&localmodel.UserModel{}).
		Select("user_table.uid, user_table.name, user_table.description, user_table.avatar_url, user_table.created_at, user_permission_table.user_type, user_permission_table.email, user_permission_table.phone, birthday").
		Joins("LEFT JOIN mblog.user_permission_table ON user_table.uid = user_permission_table.uid").
		Order(func(bool) string {
			if sort {
				return "user_table.uid ASC"
			} else {
				return "user_table.uid DESC"
			}
		}(sort)).
		Limit(size).
		Offset((page - 1) * size).
		Find(&user_list).
		Error
	return user_list, err
}

/**	通过关键字查询查询用户
 *	@param
 */
func SearchUserByKeyword(db *gorm.DB, keyword string, page, size int, sort bool) ([]UserInfoWithoutPassword, error) {
	var err error
	var result_user_list []UserInfoWithoutPassword
	var exist_signal int
	err = db.Model(&localmodel.UserModel{}).
		Select("COUNT(*)").
		Joins("LEFT JOIN mblog.user_permission_table ON user_table.uid = user_permission_table.uid").
		Where("user_table.name LIKE ?", `%`+keyword+`%`).
		First(&exist_signal).
		Error

	if exist_signal == 0 {
		return []UserInfoWithoutPassword{}, err
	} else {
		err = db.Model(&localmodel.UserModel{}).
			Select("user_table.uid, user_table.name, user_table.description, user_table.avatar_url, user_table.created_at, user_permission_table.user_type, user_permission_table.email, user_permission_table.phone, birthday").
			Joins("LEFT JOIN mblog.user_permission_table ON user_table.uid = user_permission_table.uid").
			Where("user_table.name LIKE ?", `%`+keyword+`%`).
			Order(func(bool) string {
				if sort {
					return "user_table.uid ASC"
				} else {
					return "user_table.uid DESC"
				}
			}(sort)).
			Limit(size).
			Offset((page - 1) * size).
			Find(&result_user_list).
			Error
		return result_user_list, err
	}
}

/**	查询用户量
 *	@params	db	gorm数据库连接对象
 */
func GetUserTotalNumber(db *gorm.DB) (int, error) {
	var total int
	err := db.Model(&localmodel.UserModel{}).
		Select("COUNT(*)").
		First(&total).Error
	return total, err
}
