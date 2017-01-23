r.GET("/file", func(c *gin.Context) {
		file := "./file/test.pdf"
		f, err := os.Open(file)
		if err != nil {
			c.AbortWithStatus(404)
			return
		}
		defer f.Close()
		log.Println("404")

		fi, _ := f.Stat()
		if fi.IsDir() {
			c.AbortWithStatus(404)
			return
		}
		http.ServeContent(c.Writer, c.Request, fi.Name(), fi.ModTime(), f)
	})
