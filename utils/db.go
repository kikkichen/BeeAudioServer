package utils

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var SqlDB *gorm.DB
var err error

func init() {
	SqlDB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       MYSQL_DSN,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}
}
