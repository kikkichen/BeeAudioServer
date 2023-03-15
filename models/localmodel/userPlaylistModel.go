package localmodel

import "time"

type UserPlayListModel struct {
	Uid         uint64    `gorm:"column:uid" json:"uid"`
	Pid         uint64    `gorm:"column:pid" json:"pid"`
	Name        string    `gorm:"column:name" json:"name"`
	Cover       string    `gorm:"column:coverImgUrl" json:"coverImgUrl"`
	CreateTime  time.Time `gorm:"column:createTime" json:"createTime"`
	Description string    `gorm:"column:description" json:"description"`
	Tags        string    `gorm:"column:tags" json:"tags"`
	Public      int       `gorm:"column:public" json:"public"`
}

func (UserPlayListModel) TableName() string {
	return "user_playlist_table"
}
