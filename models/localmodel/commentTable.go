package localmodel

import (
	"time"
)

type CommentModel struct {
	Cid    uint64    `gorm:"column:cid"`
	RootId uint64    `gorm:"column:root_id"`
	Text   string    `gorm:"column:text"`
	Source string    `gorm:"column:source"`
	Uid    uint64    `gorm:"column:uid"`
	Bid    uint64    `gorm:"column:bid"`
	PostAt time.Time `gorm:"column:post_at"`
	BeLike int       `gorm:"column:be_like"`
}

func (CommentModel) TableName() string {
	return "comment_table"
}
