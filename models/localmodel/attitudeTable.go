package localmodel

import "time"

type AttitudeModel struct {
	Aid     uint64    `gorm:"column:aid;primaryKey"`
	Bid     uint64    `gorm:"column:bid;index:,class:FULLTEXT,option:WITH PARSER ngram INVISIBLE"`
	Uid     uint64    `gorm:"column:uid"`
	Created time.Time `gorm:"column:created_at"`
	Source  string    `gorm:"column:source"`
}

func (AttitudeModel) TableName() string {
	return "atititude_table"
}
