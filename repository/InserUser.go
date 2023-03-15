package repository

import (
	"BeeAudioServer/models/localmodel"
	"BeeAudioServer/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

/* 新增一位新用户， 用户创建成功返回true, 创建失败返回false
 *	@params db gorm 链接对象
 *	@params	password 用户预设密码
 *
 */
func AddUser(
	db *gorm.DB,
	password string,
	email string,
	phone string,
) uint64 {
	/* 生成 Uid */
	var new_uid_tail string
	var new_uid_string string
	var new_uid uint64

	/* 循环生成Uid,并验证其是空闲的 */
	for {
		new_uid_tail = utils.RandomString(6, utils.DefaultNumber)
		new_uid_string = "9900" + new_uid_tail
		new_uid = utils.StringParseToUint64(new_uid_string)
		/* 判断新生成的Uid没有被使用 */
		if !IsExistInUserModel(db, new_uid) {
			break
		}
	}

	/* 构建新用户 gorm 模型 */
	new_user := localmodel.UserModel{Uid: new_uid, Name: fmt.Sprintf("用户%v", new_uid_string), CreatedAt: time.Now()}
	new_user_permission := localmodel.UserPermissionModel{Uid: new_uid, Password: password, UserType: 0, Email: email, Phone: phone}
	/* 先填充 user_model， 再填充 user_permission */
	result := db.Model(&localmodel.UserModel{}).Create(&new_user)
	result = db.Model(&localmodel.UserPermissionModel{}).Create(&new_user_permission)

	/* 出现新用户数据插入错误，则返回字符串0 */
	if result.Error != nil {
		new_uid = 0
	}
	return new_uid
}
