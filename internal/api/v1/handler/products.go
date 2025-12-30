package v1handler

import (
	"fmt"
	"net/http"

	"golang-backend-fundamental-1/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductsHandler struct {}

type GetProductsBySlugV1Params struct {
	Slug string `uri:"slug" binding:"slug.min=3,slug.max=20"`
}

type GetProductsV1Params struct {
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
	Search string `form:"search" binding:"required,search"`
	Email string `form:"email" binding:"omitempty,email"`
	Date string `form:"date" binding:"omitempty,datetime=2006-01-02"`
}


type ProductImage struct {
	ImageName string `json:"image_name" binding:"required,min=3,max=50"`
	ImageLink string `json:"image_link" binding:"required,file_ext=jpg,png,jpeg,gif,webp"`
}

type ProductAttribute struct {
	AttributeName string `json:"attribute_name" binding:"required,min=3,max=50"`
	AttributeValue string `json:"attribute_value" binding:"required,min=3,max=50"`
}

type ProductInfo struct {
	InfoKey string `json:"info_key" binding:"required,min=3,max=50"`
	InfoValue string `json:"info_value" binding:"required,min=3,max=50"`
}

type CreateProductV1Params struct {
	Name string `json:"name" binding:"required,min=3,max=50"`
	Price int `json:"price" binding:"required,min_int=100000,max_int=1000000"`
	Display *bool `json:"display" binding:"omitempty"`
	ProductImage ProductImage `json:"product_image" binding:"required"`
	Tags []string `json:"tags" binding:"required,gt=3,lt=5"`
	ProductAttributes []ProductAttribute `json:"product_attributes" binding:"required,gt=0,dive"`
	ProductInfo map[string]ProductInfo `json:"product_info" binding:"required,dive"`
	ProductMetadata map[string]any `json:"product_metadata" binding:"omitempty"`
}

 


func NewProductsHandler() *ProductsHandler {
	return &ProductsHandler{}
}



func (p *ProductsHandler) GetProductsV1(c *gin.Context) {


var params GetProductsV1Params
if err := c.ShouldBindQuery(&params); err != nil {
	c.JSON(http.StatusBadRequest, gin.H{
		"message": utils.HandleValidationErrors(err),
	})
	return
}

if params.Limit == 0 {
	params.Limit = 10
}

if params.Email != "" {
	params.Email = "No Email"
}

if params.Date != "" {
	params.Date = "No Date"
}

c.JSON(http.StatusOK, gin.H{
	"message": "Products",
	"limit": params.Limit,
	"search": params.Search,
	"email": params.Email,
	"date": params.Date,
	})
}

func (p *ProductsHandler) GetProductByIdV1(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Product ID: " + id,
	})
}

func (p *ProductsHandler) CreateProductV1(c *gin.Context) {
	
	var params CreateProductV1Params
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.HandleValidationErrors(err),
		})
		return
	}


	for key, _ := range params.ProductInfo {
		if _, err := uuid.Parse(key); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": gin.H{
					"product_info": fmt.Sprintf("Invalid UUID: %s", key),
				}	,
			})
			return
		}	
	}

	if params.Display == nil {
	defaultDisplay := true
	params.Display = &defaultDisplay
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product Name: " ,
		"data": params,	
	})
}

func (p *ProductsHandler) UpdateProductV1(c *gin.Context) {
	id := c.Param("id")
	name := c.PostForm("name")
	c.JSON(http.StatusOK, gin.H{
		"message": "Product ID: " + id + " Product Name: " + name,
	})
}

func (p *ProductsHandler) DeleteProductV1(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusNoContent, gin.H{
		"message": "Product ID: " + id + " deleted",
	})
}

func (p *ProductsHandler) GetProductBySlugV1(c *gin.Context) {
	var params GetProductsBySlugV1Params
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.HandleValidationErrors(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product Slug: " + params.Slug,
	})
}