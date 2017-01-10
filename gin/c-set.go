package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"time"
	"log"
	"fmt"
)

type Config struct {
	Name  string
	Mode  int
	Start time.Time
}

func main() {
	router := gin.New()

	config := &Config{
		Name:"gin web app",
		Mode:0,
		Start:time.Now(),
	}

	fmt.Printf("%p\n", config)
	router.Use(func(c *gin.Context) {
		c.Set("config", config)
	})

	router.GET("/", func(c *gin.Context) {
		config := c.MustGet("config").(*Config)
		fmt.Printf("path: %s, %p\n", c.Request.URL.Path, config)
		log.Println(config.Start.UnixNano())

		c.String(200, "running...")
	})

	router.GET("/login", func(c *gin.Context) {
		config := c.MustGet("config").(*Config)
		fmt.Printf("%p\n", config)
		log.Println("before", config.Start.UnixNano())
		config.Start = time.Now()
		c.Next()
	}, func(c *gin.Context) {
		config := c.MustGet("config").(*Config)
		fmt.Printf("path: %s, %p\n", c.Request.URL.Path, config)
		log.Println("after", config.Start.UnixNano())
		c.String(200, "login")
	})

	router.Run(":8089")
}
