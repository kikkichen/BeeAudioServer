package localmodel

type UserFavoritePlaylist struct {
	Uid uint64 `gorm:"column:uid;primaryKey"`
	Pid uint64 `gorm:"column:pid"`
}

func (UserFavoritePlaylist) TableName() string {
	return "user_favorite_playlist"
}
