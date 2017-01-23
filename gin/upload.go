package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"dengzhanglin/gin/static"
	"os"
	"net/http"
	"fmt"
	"io"
	"log"
	"time"
	"strings"
	"path/filepath"
	"mime"
)

// 获取大小的借口
// io.SectionReader.Size() int64
type sizer interface {
	Size() int64
}

// os.File.Stat()(os.FileInof,error)
type stater interface {
	Stat() (os.FileInfo, error)
}

func main() {
	r := gin.Default()

	// if Allow DirectoryIndex
	//r.Use(static.Serve("/", static.LocalFile("/tmp", true)))
	// set prefix
	r.Use(static.Serve("/asset/", static.LocalFile("./assets", false)))

	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.String(200, "pong " + fmt.Sprint(time.Now().Unix()))
	})

	r.GET("/html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

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

	r.POST("/upload", func(c *gin.Context) {

		err := c.Request.ParseMultipartForm(3 << 20)
		if err != nil {
			log.Println(err, "file size")
			c.AbortWithStatus(500)
			return
		}

		if c.Request.MultipartForm.File == nil {
			log.Println(err, "MultipartForm")
			c.AbortWithStatus(500)
			return
		}

		file, header, err := c.Request.FormFile("file")

		if err != nil {
			log.Println(err, "form file")
			c.AbortWithStatus(500)
			return
		}
		var size int64
		switch t := file.(type) {
		case stater:
			stat, err := t.Stat()
			log.Println("stater")
			if err != nil {
				log.Println(err, "file type")
				c.AbortWithStatus(500)
				return
			}
			size = stat.Size()
		case sizer:
			size = t.Size()
		default:
			size = 0

		}

		// 文件大小限制
		if size >= (3 << 20) {
			log.Println(err, "file size error")
			c.AbortWithStatus(500)
			return
		}

		filename := header.Filename
		fmt.Println(header.Filename, "file name")

		// 文件后缀
		ext := strings.ToLower(filepath.Ext(filename))
		log.Println(ext, "ext")

		out, err := os.Create("./tmp/" + filename)
		if err != nil {
			log.Println(err, "create")
			c.AbortWithStatus(500)
			return
		}
		defer out.Close()

		fileSize, err := io.Copy(out, file)
		if err != nil {
			log.Println(err, "copy")
			c.AbortWithStatus(500)
			return
		}

		//fileinfo, err := out.Stat()
		//if err != nil {
		//	log.Println(err, "fileinfo")
		//}

		mimeType := mime.TypeByExtension(ext)

		log.Println("file size 1", size)
		log.Println("file size 2", fileSize)
		log.Println("file mime", mimeType)
		//log.Println("file mime", fileinfo)
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
