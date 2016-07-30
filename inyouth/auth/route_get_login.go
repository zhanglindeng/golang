package auth

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "log"
    "dzlin.com/inyouth/ct"
    "dzlin.com/inyouth/err"
)

func getLogin(c *gin.Context) {

    log.Println(ct.CT_ONE)
    log.Println(err.ErrEmail)
    log.Println(authConf)

    log.Println("Referer", c.Request.Referer())

    c.HTML(http.StatusOK, "login.html", nil)

    //c.JSON(http.StatusOK, gin.H{
    //    "code":0,
    //    "data":"hello",
    //})
}
