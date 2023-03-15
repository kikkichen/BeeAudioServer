package localmodel

type AlbumModel struct {
	Id     uint64 `gorm:"column:id" json:"id"`
	ArId   uint64 `gorm:"column:ar_id" json:"ar_id"`
	Name   uint64 `gorm:"column:name" json:"name"`
	PicStr string `gorm:"column:pic_str" json:"pic_str"`
	Pic    uint64 `gorm:"column:pic" json:"pic"`
}

func (AlbumModel) TableName() string {
	return "album_table"
}
