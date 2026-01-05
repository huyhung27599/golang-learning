package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)
type Client struct {
	limiter *rate.Limiter
	lastSeen time.Time
}


var (
	mu sync.Mutex
 clients = make(map[string]*Client)
)

func getClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	return ip
}

func getRateLimiter (ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
   client, exists := clients[ip]
   if !exists {
	limiter := rate.NewLimiter(rate.Every(1*time.Second), 10)
	client = &Client{limiter: limiter, lastSeen: time.Now()}
	clients[ip] = client
   }
   client.lastSeen = time.Now()
   return client.limiter
}

func CleanupOldClients() {
	for {
		time.Sleep(1*time.Minute)
		mu.Lock()
	for ip, client := range clients {
		if time.Since(client.lastSeen) > 3*time.Minute {
			delete(clients, ip)
		}
	}
	mu.Unlock()
}
}

func RateLimitingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getClientIP(c)
		limiter := getRateLimiter(ip)
		if limiter.Allow() {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"message": "Too many requests",
		})
	}
}