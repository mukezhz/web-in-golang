package controllers

import (
	"fmt"
	"go-file-upload/models"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileController struct {
	DB *gorm.DB
}

func (c *FileController) UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filePath := filepath.Join("uploads", file.Filename)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	uuidV4 := uuid.New().String()
	fileMetadata := models.File{
		Filename: file.Filename,
		UUID:     uuidV4,
	}
	if err := c.DB.Create(&fileMetadata).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "Details": fileMetadata})
}

func (c *FileController) UploadFiles(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	files := form.File["files"]
	var fileModels []models.File
	for _, file := range files {
		filePath := filepath.Join("uploads", file.Filename)
		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		fileModels = append(fileModels, models.File{
			UUID:     uuid.New().String(),
			Filename: file.Filename,
		})
	}
	err = c.DB.Create(&fileModels).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file information"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"files":   fileModels,
	})
}

func (c *FileController) GetFile(ctx *gin.Context) {
	uuidV4 := ctx.Param("uuid")
	var file models.File
	err := c.DB.Where("uuid = ?", uuidV4).First(&file).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	filePath := filepath.Join("uploads", file.Filename)
	fileData, err := os.Open(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer func(fileData *os.File) {
		err := fileData.Close()
		if err != nil {

		}
	}(fileData)
	// Read the first 512 bytes of the file to determine its content type
	fileHeader := make([]byte, 512)
	_, err = fileData.Read(fileHeader)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	fileContentType := http.DetectContentType(fileHeader)
	fileInfo, err := fileData.Stat()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
		return
	}

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Filename))
	ctx.Header("Content-Type", fileContentType)
	ctx.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	ctx.File(filePath)
}

func (c *FileController) GetFileMetadata(ctx *gin.Context) {
	uuidV4 := ctx.Param("uuid")
	var file models.File
	err := c.DB.Where("uuid = ?", uuidV4).First(&file).Error

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"file":    file,
	})
}

func (c *FileController) GetAllFileMetadata(ctx *gin.Context) {
	var files []models.File
	_, err := c.DB.Find(&files).Rows()
	for _, file := range files {
		log.Println(file)
	}
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Error while marshalling"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"files":   files,
	})
}

func (c *FileController) DeleteFile(ctx *gin.Context) {
	uuidV4 := ctx.Param("uuid")
	var file models.File
	err := c.DB.Where("uuid = ?", uuidV4).First(&file).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	filePath := filepath.Join("uploads", file.Filename)
	err = os.Remove(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from upload folder"})
		return
	}

	err = c.DB.Delete(&file).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from database"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "File " + file.Filename + " deleted successfully",
	})
}
