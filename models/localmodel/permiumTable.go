package localmodel

import "time"

/* Permium 会员 订单表 */
type PremiumModel struct {
	Uid           uint64    `gorm:"column:uid;primaryKey" json:"uid"`
	CardId        string    `gorm:"column:card_id;primaryKey" json:"card_id"`
	CardType      int       `gorm:"column:card_type" json:"card_type"`
	ServerIn      time.Time `gorm:"column:server_in" json:"server_in"`
	ServerExpired time.Time `gorm:"column:server_expired" json:"server_expired"`
}

func (PremiumModel) TableName() string {
	return "permium_table"
}
