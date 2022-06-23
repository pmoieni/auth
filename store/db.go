package store

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	gormLogger "gorm.io/gorm/logger"
)

var db *gorm.DB

type DBConfig struct {
	DBDriver   string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
}

// initialize a db connection
func InitDB(c *DBConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})
	if err != nil {
		return
	}

	db.AutoMigrate(&User{})
	return
}
