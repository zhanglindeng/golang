package main

import (
	"fmt"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

func main() {
	r := gin.New()
	// gin.SetMode(gin.ReleaseMode)

	r.Use(setStartTime)
	r.Use(replaceResponseWriter)
	r.Use(afterRequest)

	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.String(200, "pong " + fmt.Sprint(time.Now().Unix()))
	})

	r.GET("/html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	r.GET("/json", func(c *gin.Context) {
		c.JSON(200, gin.H{"code":0, "message":"success"})
	})

	r.GET("/err", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code":40000, "message":"not found"})
	})

	r.Run(":8080")
}
