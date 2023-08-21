package database

import (
	"fmt"
	"go-file-upload/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GlobalDB *gorm.DB

func InitDatabase(config map[string]string) (err error) {
	dsn := fmt.Sprintf(
		"%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config["DB_USERNAME"],
		config["DB_PASSWORD"],
		config["DB_HOST"],
		config["DB_DATABASE"],
	)

	GlobalDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	err = GlobalDB.AutoMigrate(&models.File{})
	if err != nil {
		return
	}

	return
}
