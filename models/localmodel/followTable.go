package localmodel

type FollowModel struct {
	FollowUid   uint64 `gorm:"column:follow_uid;primaryKey"`
	BeFollowUid uint64 `gorm:"column:be_follow_uid;primaryKey"`
}

func (FollowModel) TableName() string {
	return "follow_table"
}
