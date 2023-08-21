package main

import (
	"go-file-upload/controllers"
	"go-file-upload/database"
	"go-file-upload/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = database.InitDatabase()
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
