package localmodel

type ArtistModel struct {
	Id        uint64 `gorm:"column:id" json:"id"`
	Cover     string `gorm:"column:cover" json:"cover"`
	Name      string `gorm:"column:name" json:"name"`
	BriefDesc string `gorm:"column:brief_desc" json:"brief_desc"`
}

func (ArtistModel) TableName() string {
	return "artist_table"
}
