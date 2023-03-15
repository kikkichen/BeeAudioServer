package localmodel

import "time"

type PlayListTable struct {
	Pid        uint64    `gorm:"column:id"`
	Uid        uint64    `gorm:"column:uid"`
	SongId     uint64    `gorm:"column:song_id"`
	CreateTime time.Time `gorm:"column:create_at"`
}

func (PlayListTable) TableName() string {
	return "playlist_table"
}
