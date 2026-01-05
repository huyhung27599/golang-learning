package middlewares

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ApiKeyMiddleware() gin.HandlerFunc {

	expectedApiKey := os.Getenv("API_KEY")
	if expectedApiKey == "" {
		expectedApiKey = "secret"
	}

	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "API key is required",
			})
			return
		}

		if apiKey != expectedApiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid API key",
			})
			return
		}

		c.Next()
	}
}