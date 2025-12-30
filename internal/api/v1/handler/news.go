package v1handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"golang-backend-fundamental-1/utils"

	"github.com/gin-gonic/gin"
)

type NewsHandler struct{}


type CreateNewsV1Params struct {
	Title string `form:"title" binding:"required,min=3,max=50"`
	Status string `form:"status" binding:"required,oneof=1 2"`
}

func NewNewsHandler() *NewsHandler {
	return &NewsHandler{}
}

func (n *NewsHandler) GetNewsV1(c *gin.Context) {

	slug := c.Param("slug")

	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "News not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "News",
	})
}

func (n *NewsHandler) CreateNewsV1(c *gin.Context) {
	var params CreateNewsV1Params
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.HandleValidationErrors(err),
		})
		return
	}

	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Image is required",
		})
		return
	}

	if image.Size > 5<<20 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Image size is too large (max 5MB)",
		})
		return
	} 

	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create uploads directory",
		})
		return
	}

	des := fmt.Sprintf("./uploads/%s", filepath.Base(image.Filename))

	if err := c.SaveUploadedFile(image, des); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to save uploaded file",
		})
		return
	}

	
	c.JSON(http.StatusCreated, gin.H{
		"message": "News Name: " + params.Title,
		"status": params.Status,
		"image": des,
		"image_name": filepath.Base(image.Filename),
	})
}

func (n *NewsHandler) UploadFileNewV1(c *gin.Context) {
	var params CreateNewsV1Params
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.HandleValidationErrors(err),
		})
		return
	}


	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Image is required",
		})
		return
	}

	fileName, err := utils.ValidateAndSaveFile( image, "./uploads")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "News Name: " + params.Title,
		"status": params.Status,
		"image": fileName,
	})
}

func (n *NewsHandler) UploadMultipleFilesNewV1(c *gin.Context) {
	const publicUrl = "http://localhost:8080/uploads/"
	var params CreateNewsV1Params
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.HandleValidationErrors(err),
		})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get multipart form",
		})
		return
	}

	images := form.File["images"]
	if len(images) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Images are required",
		})
		return
	}


	var succesFiles []string
	var failedFiles []map[string]string
	for _, image := range images {
		fileName, err := utils.ValidateAndSaveFile(image, "./uploads")
		if err != nil {
			failedFiles = append(failedFiles, map[string]string{
				"filename": image.Filename,
				"error": err.Error(),
			})
			continue
		}
		publicImageUrl := publicUrl + fileName
		succesFiles = append(succesFiles, publicImageUrl)
	}
	
	res := gin.H{
		"message": "Create News",
		"status": params.Status,
		"images": succesFiles,
		
	}

	if len(failedFiles) > 0 {
		res["message"] = "Some files failed to upload"
		res["failed_files"] = failedFiles
	}
	c.JSON(http.StatusCreated, res)
}