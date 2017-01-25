package main

import (
	"github.com/didip/tollbooth"
	"time"
	"./middleware"
	"github.com/gin-gonic/gin"
)

//func HelloHandler(w http.ResponseWriter, req *http.Request) {
//	w.Write([]byte("Hello, World!"))
//}
//
//func main() {
//	// Create a request limiter per handler.
//	http.Handle("/", tollbooth.LimitFuncHandler(tollbooth.NewLimiter(1, 1* time.Second), HelloHandler))
//	http.ListenAndServe(":12345", nil)
//}

func main() {
	r := gin.New()

	// Create a limiter struct.
	limiter := tollbooth.NewLimiter(1, time.Second)
	r.GET("/", middleware.LimitHandler(limiter), func(c *gin.Context) {
		c.String(200, "Hello, world!")
	})

	r.Run(":12345")
}

// 其他参考 https://github.com/VojtechVitek/ratelimit
