package localmodel

/* 播放历史记录 - 数据元 */
type HistorySong struct {
	SongId    uint64   `json:"sid"`
	SongTitle string   `json:"title"`
	ArtistIds []uint64 `json:"ar_ids"`
	AlbumIds  uint64   `json:"al_ids"`
	PlayAt    int64    `json:"play_at"`
}

/* 数据库中的历史记录模型 */
type HistoryDataModel struct {
	Uid         uint64 `gorm:"column:uid"`
	HistoryData string `gorm:"column:history_data"`
}

func (*HistoryDataModel) TableName() string {
	return "history_play_table"
}
