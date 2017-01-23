package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"dengzhanglin/gin/static"
	"os"
	"net/http"
	"fmt"
)

func main() {
	r := gin.Default()

	// if Allow DirectoryIndex
	//r.Use(static.Serve("/", static.LocalFile("/tmp", true)))
	// set prefix
	r.Use(static.Serve("/asset/", static.LocalFile("./assets", false)))

	//r.Use(static.Serve("/", static.LocalFile("/tmp", false)))
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "test")
	})

	r.GET("/file", func(c *gin.Context) {
		filepath := "./file/test.pdf"
		file(c, filepath)
	})

	r.GET("/download", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("%s; filename=%s", "attachment", "test.pdf"))
		file(c, "./file/test.pdf")
	})
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

func file(c *gin.Context, file string) {
	f, err := os.Open(file)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}
	defer f.Close()

	fi, _ := f.Stat()
	if fi.IsDir() {
		c.AbortWithStatus(404)
		return
	}
	http.ServeContent(c.Writer, c.Request, fi.Name(), fi.ModTime(), f)
}
