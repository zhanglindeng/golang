package user

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "log"
)

func getIndex(c *gin.Context) {
    log.Println("user index")
    c.String(http.StatusOK, "Hello %s", c.Param("name"))
}
