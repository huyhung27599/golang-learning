package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
)

func SimPleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		log.Println("Starting Simple Middleware")
		c.Writer.Write([]byte("Starting Simple Middleware"))
		c.Next()
		c.Writer.Write([]byte("Ending Simple Middleware"))
		log.Println("Ending Simple Middleware")
	}
}