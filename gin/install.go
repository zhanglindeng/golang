package main

import (
	"fmt"
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/Unknwon/com"
	"github.com/go-xweb/log"
)

type config struct {
	installed bool
}

func main() {
	fmt.Println("hello gin")

	router := gin.New()

	// cfg
	cfg := &config{installed:false}

	// 是否有 install.lock 文件
	if com.IsExist("./install.lock") {
		cfg.installed = true
	}

	router.Use(func(c *gin.Context) {
		c.Set("cfg", cfg)
		c.Next()
	})

	// 正常
	r := router.Group("")
	{
		// 是否安裝
		r.Use(func(c *gin.Context) {
			cfg1 := c.MustGet("cfg").(*config)
			if !cfg1.installed {
				c.Redirect(302, "install")
				return
			}
			c.Next()
		})
		r.GET("", func(c *gin.Context) {
			c.String(200, "index")
		})
		r.GET("test", func(c *gin.Context) {
			c.String(200, "test")
		})

		// 重新安裝
		r.GET("installed", func(c *gin.Context) {
			cfg1 := c.MustGet("cfg").(*config)
			cfg1.installed = false;
			c.Redirect(302, "install")
		})
	}


	// install
	install := router.Group("install")
	{
		install.Use(func(c *gin.Context) {
			cfg1 := c.MustGet("cfg").(*config)
			if cfg1.installed {
				c.AbortWithStatus(404)
				return
			}
			log.Println("no install")
			c.Next()
		})
		install.GET("", func(c *gin.Context) {
			c.String(200, "install")
		})
		install.GET("do", func(c *gin.Context) {
			// 安裝完成后
			cfg1 := c.MustGet("cfg").(*config)
			cfg1.installed = true;
			c.String(200, "installing...")
		})
	}

	router.Run(":8080")
}
