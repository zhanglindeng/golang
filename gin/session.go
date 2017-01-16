package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	fmt.Println("hello gin")

	router := gin.New()

	store := sessions.NewCookieStore([]byte("rMDSbOoerUUXxHyZHGBFlLUoEDfCXBBH"))
	store.Options(sessions.Options{
		MaxAge:0,
		HttpOnly:true,
	})
	router.Use(sessions.Sessions("SESSIONID", store))

	router.GET("sess", func(c *gin.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count += 1
		}

		session.Set("count", count)
		session.Save()
		c.JSON(200, gin.H{"count": count})
	})

	router.Run(":8080")
}
