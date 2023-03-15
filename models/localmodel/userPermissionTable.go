package localmodel

import "time"

/* 用户密码及权限与其他个人信息 */
type UserPermissionModel struct {
	Uid      uint64    `gorm:"column:uid;primaryKey"`
	Password string    `gorm:"column:password"`
	UserType int       `gorm:"column:user_type"`
	Email    string    `gorm:"column:email"`
	Phone    string    `gorm:"column:phone"`
	Birthday time.Time `gorm:"column:birthday"`
}

func (UserPermissionModel) TableName() string {
	return "user_permission_table"
}
