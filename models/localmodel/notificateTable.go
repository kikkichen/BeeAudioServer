package localmodel

import "time"

type notificateModel struct {
	ActionUid      uint64    `gorm:"column:action_id"`
	NotificateUid  uint64    `gorm:"column:notificate_uid,index:,class:FULLTEXT,option:WITH PARSER ngram INVISIBLE"`
	NotificateType int       `gorm:"column:notificate_type"`
	NotificateId   uint64    `gorm:"column:notificate_id"`
	NotificateTime time.Time `gorm:"column:notificate_time"`
	isReaded       bool      `gorm:"column:readed"`
}

func (notificateModel) TableName() string {
	return "notificate_table"
}
