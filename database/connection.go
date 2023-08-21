package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"go-file-upload/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var GlobalDB *gorm.DB

func InitDatabase() (err error) {
	config, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error reading .env file")
	}

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
