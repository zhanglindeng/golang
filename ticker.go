import (
    "log"
    "net/http"
     "time"
)

func main() {
    r.GET("", func(ctx *gin.Context) {
		start := time.Now()
		t := time.NewTicker(3 * time.Second)
		counter := 0
		maxCounter := 5
		defer t.Stop()
		for {
			<-t.C
			log.Println("time:", time.Now())
			counter++
			if counter >= maxCounter {
				break
			}
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "OK",
			"start":   start.Format("2006-01-02 15:04:05"),
			"end":     time.Now().Format("2006-01-02 15:04:05"),
		})
	})
}
