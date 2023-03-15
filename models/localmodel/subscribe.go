package localmodel

type SubscribeModel struct {
	Uid  uint64 `gorm:"column:uid" json:"uid"`
	Data string `gorm:"column:subscribe_data" json:"subscribe_data"`
}

func (SubscribeModel) TableName() string {
	return "subscribe_table"
}
