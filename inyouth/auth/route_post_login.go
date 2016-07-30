package auth

import (
    "github.com/gin-gonic/gin"
    "io/ioutil"
    "log"
    "encoding/json"
)

func postLogin(c *gin.Context) {
    log.Println("postLogin")
    
    sess, _ := authConf.Session.SessionStart(c.Writer, c.Request)
    defer sess.SessionRelease(c.Writer)

    username := c.Request.FormValue("username")
    sess.Set("user", username)

    c.String(200, "hello %s", username)
}

func postLogin1(c *gin.Context) {

    log.Println(authConf)

    body, err := ioutil.ReadAll(c.Request.Body)
    if err != nil {
        log.Println(err)
        c.Abort()
    }

    sess, _ := authConf.Session.SessionStart(c.Writer, c.Request)
    defer sess.SessionRelease(c.Writer)
    log.Println("SessionStart")

    log.Println(string(body))

    data := make(map[string]string)
    err = json.Unmarshal(body, &data)
    if err != nil {
        c.AbortWithError(500, err)
    }

    log.Println(data)
    sess.Set("user", data["user"])

}


