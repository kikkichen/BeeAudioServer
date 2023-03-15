package localmodel

import (
	"BeeAudioServer/models/responsemodel"
	"time"
)

type UserModel struct {
	Uid         uint64    `gorm:"column:uid;primaryKey" json:"uid"`
	Name        string    `gorm:"column:name" json:"screen_name"`
	Description string    `gorm:"column:description" json:"description"`
	AvatarUrl   string    `gorm:"column:avatar_url" json:"profile_image_url"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"create_at"`
}

func (UserModel) TableName() string {
	return "user_table"
}

/* UserModel 映射为 携带关注关系属性的 ResponseUser类型 */
func (u *UserModel) MapToResponseUser() responsemodel.ResponseUser {
	return responsemodel.ResponseUser{
		Uid:         u.Uid,
		Name:        u.Name,
		Description: u.Description,
		AvatarUrl:   u.AvatarUrl,
		CreatedAt:   u.CreatedAt,
	}
}
