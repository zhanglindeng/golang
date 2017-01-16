package main

import (
	"fmt"
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/gin-contrib/gzip"
	"time"
)

func main() {
	fmt.Println("hello gin")

	r := gin.New()

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	r.Run(":8080")
}
