package main

import (
    "github.com/gin-gonic/gin"
    "dzlin.com/inyouth/conf"
    "dzlin.com/inyouth/auth"
    "log"
    "dzlin.com/inyouth/user"
    "github.com/astaxie/beego/session"
)

func main() {

    log.Println("inyouth")

    // router
    router := gin.New()

    // session
    sessionConfig := `{"cookieName":"PHPSESSID","gclifetime":1800,"ProviderConfig":"./data/tmp"}`
    session, err := session.NewManager("file", sessionConfig)
    if err != nil {
        log.Fatalln("error:", err.Error())
    }
    go session.GC()

    // config
    conf := &conf.Conf{Router:router, Session:session}

    // modules
    auth.AuthRoute(conf)
    user.UserRoute(conf)

    // run
    conf.Router.Run(":9090")
}
