package conn

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	lib "majoo/lib"
)

var (
	DBConn *gorm.DB
)

func DBEstablish() *gorm.DB {
	dsn := lib.GetConfig("DB_USER") + ":" + lib.GetConfig("DB_PASS") + "@tcp(" + lib.GetConfig("DB_HOST") + ":" + lib.GetConfig("DB_PORT") + ")/" + lib.GetConfig("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func DBConnection() *gorm.DB {
	return DBConn
}
