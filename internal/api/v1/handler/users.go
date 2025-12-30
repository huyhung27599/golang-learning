package v1handler

import (
	"golang-backend-fundamental-1/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UsersHandler struct {}


type GetUserByIdV1Params struct {
	ID int `uri:"id" binding:"gt=0"`
}

type GetUserByUuidV1Params struct {
	UUID string `uri:"uuid" binding:"uuid"`
}

func NewUsersHandler() *UsersHandler {
	return &UsersHandler{}
}



func (u *UsersHandler) GetUsersV1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func (u *UsersHandler) GetUserByIdV1(c *gin.Context) {
	var params GetUserByIdV1Params


	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.HandleValidationErrors(err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		
	})
}

func (u *UsersHandler) CreateUserV1(c *gin.Context) {
	name := c.PostForm("name")
	c.JSON(http.StatusCreated, gin.H{
		"message": "User Name: " + name,
	})
}

func (u *UsersHandler) UpdateUserV1(c *gin.Context) {
	id := c.Param("id")
	name := c.PostForm("name")
	c.JSON(http.StatusOK, gin.H{
		"message": "User ID: " + id + " User Name: " + name,
	})
}

func (u *UsersHandler) DeleteUserV1(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusNoContent, gin.H{
		"message": "User ID: " + id + " deleted",
	})
}

func (u *UsersHandler) GetUserByUuidV1(c *gin.Context) {
	var params GetUserByUuidV1Params
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": utils.HandleValidationErrors(err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User UUID: " + params.UUID,
	})
}