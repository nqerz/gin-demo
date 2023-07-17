package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	secret := os.Getenv("SECRET_KEY_VALUE")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong " + secret,
			"time": time.Now().Format(time.RFC3339),
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
