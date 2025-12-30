package v1handler

import (
	"net/http"

	"golang-backend-fundamental-1/utils"

	"github.com/gin-gonic/gin"
)



type CategoriesHandler struct {}

type CreateCategoryV1Params struct {
	Name string `form:"name" binding:"required,min=3,max=50"`
	Status string `form:"status" binding:"required,oneof=1 2"`
}


type GetCategoryByCategoryCodeV1Params struct {
	CategoryCode string `uri:"category_code" binding:"oneof=1 2 3 4 5 6"`
}

func NewCategoriesHandler() *CategoriesHandler {
	return &CategoriesHandler{}
}


func (ca *CategoriesHandler) GetCategoriesV1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Categories",
	})
}

func (ca *CategoriesHandler) GetCategoryByIdV1(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Category ID: " + id,
	})
}	

func (ca *CategoriesHandler) CreateCategoryV1(c *gin.Context) {
	var params CreateCategoryV1Params
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.HandleValidationErrors(err),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Category Name: " + params.Name,
		"status": params.Status,
	})
}

func (ca *CategoriesHandler) UpdateCategoryV1(c *gin.Context) {
	id := c.Param("id")
	name := c.PostForm("name")
	c.JSON(http.StatusOK, gin.H{
		"message": "Category ID: " + id + " Category Name: " + name,
	})
}

func (ca *CategoriesHandler) DeleteCategoryV1(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusNoContent, gin.H{
		"message": "Category ID: " + id + " deleted",
	})
}

func (ca *CategoriesHandler) GetCategoryByCategoryCodeV1(c *gin.Context) {
	var params GetCategoryByCategoryCodeV1Params
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.HandleValidationErrors(err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Category Code: " + params.CategoryCode,
	})
}