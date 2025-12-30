package v2handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UsersHandler struct {}




func NewUsersHandler() *UsersHandler {
	return &UsersHandler{}
}



func (u *UsersHandler) GetUsersV2(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func (u *UsersHandler) GetUserByIdV2(c *gin.Context) {
	
	c.JSON(http.StatusOK, gin.H{
		"message": "User ID:" ,
	})
}

func (u *UsersHandler) CreateUserV2(c *gin.Context) {
	name := c.PostForm("name")
	c.JSON(http.StatusCreated, gin.H{
		"message": "User Name: " + name,
	})
}

func (u *UsersHandler) UpdateUserV2(c *gin.Context) {
	id := c.Param("id")
	name := c.PostForm("name")
	c.JSON(http.StatusOK, gin.H{
		"message": "User ID: " + id + " User Name: " + name,
	})
}

func (u *UsersHandler) DeleteUserV2(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusNoContent, gin.H{
		"message": "User ID: " + id + " deleted",
	})
}