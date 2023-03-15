package localmodel

type SongTable struct {
	Id        string `gorm:"id" json:"id"`
	Name      string `gorm:"name" json:"name"`
	ArId      string `gorm:"ar_id" json:"ar_id"`
	AlId      string `gorm:"al_id" json:"al_id"`
	LocalPath string `gorm:"local_path" json:"local_path"`
	Privilege int    `gorm:"privilege" json:"privilege"`
	Quality   string `gorm:"quality" json:"quality"`
	Useful    int    `gorm:"useful" json:"useful"`
	Source    string `gorm:"source" json:"source"`
}

func (SongTable) TableName() string {
	return "song_table"
}
