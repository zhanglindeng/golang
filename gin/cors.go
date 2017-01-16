package main

import (
	"fmt"
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/gin-contrib/gzip"
	"time"
	"gopkg.in/gin-contrib/cors.v1"
)

func main() {
	fmt.Println("hello gin")

	r := gin.New()

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:63342"}
	r.Use(cors.New(config))
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	r.Run(":8080")
}
