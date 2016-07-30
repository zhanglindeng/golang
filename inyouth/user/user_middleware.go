package user

import (
    "github.com/gin-gonic/gin"
    "log"
)

func authRequired() gin.HandlerFunc {

    return func(c *gin.Context) {
        log.Println("before authRequired")

        sess, _ := userConf.Session.SessionStart(c.Writer, c.Request)
        defer sess.SessionRelease(c.Writer)
        log.Println("SessionStart")

        user := sess.Get("user")
        if user == nil {
            c.Redirect(302, "/auth/login")
            c.AbortWithStatus(302)
            log.Println("Redirect /auth/login")
        } else {
            log.Println("session user", user.(string))
        }

        c.Next()
        log.Println("after authRequired")
    }
}
