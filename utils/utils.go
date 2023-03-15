package utils

import (
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const NET_TOKEN = "bd28b9d4c22dcd0ffd026840822f6406c07f1bc9033852b594e17f11dc5ba274993166e004087dd3462d4851999c23629ce2a6ad2e9c3f3f3c7baf0e78175d9fa138882140a3d995a89fe7c55eac81f3"
const NET_ACCESS_ADDRESS = "http://127.0.0.1:3000"

/* 数据库连接字段 */
const MYSQL_DSN = "shigure:123666@tcp(127.0.0.1:3306)/mblog?charset=utf8mb4&parseTime=True&loc=Local"

const LOCAL_NETCLOUD_ADDRESS = "http://127.0.0.1:3000"
const LOCAL_OAUTH2_SERVER = "http://127.0.0.1:9096"

const BASIC_AUTH_USERNAME = "test_client_1"
const BASIC_AUTH_PASSWORD = "test_secret_1"

var once sync.Once

func GetGormDB(db *gorm.DB) *gorm.DB {
	once.Do(func() {
		var err error
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN:                       MYSQL_DSN,
			DefaultStringSize:         256,
			DisableDatetimePrecision:  true,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			SkipInitializeWithVersion: false,
		}), &gorm.Config{})

		if err != nil {
			log.Fatal(err)
		}
	})
	return db
}
