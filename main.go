package main

import (
	"github.com/joho/godotenv"
	"go-file-upload/controllers"
	"go-file-upload/database"
	"go-file-upload/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error while reading .env")
	}
	err = database.InitDatabase(config)
	if err != nil {
		log.Fatal("Error initializing database")
	}
	db := database.GlobalDB
	log.Println("Connected to database!")

	err = db.AutoMigrate(&models.File{})
	if err != nil {
		log.Fatal("Error migrating database")
	} else {
		log.Println("Migration Sucessful!")
	}

	fileController := &controllers.FileController{DB: db}

	r := gin.Default()
	r.POST("/file", fileController.UploadFile)
	r.POST("/files", fileController.UploadFiles)
	r.GET("/file/:uuid", fileController.GetFile)
	r.GET("/metadata/:uuid", fileController.GetFileMetadata)
	r.GET("/metadata", fileController.GetAllFileMetadata)
	r.DELETE("/file/:uuid", fileController.DeleteFile)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Error starting server")
	}
}
