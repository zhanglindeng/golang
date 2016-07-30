package auth

import (
    "dzlin.com/inyouth/conf"
)

func AuthRoute(conf *conf.Conf) {

    authConf = conf
    conf.Router.LoadHTMLGlob("data/html/auth/*")
    auth := conf.Router.Group("/auth")
    auth.Use()
    {
        auth.GET("login", getLogin)
        auth.POST("login", postLogin)
    }
}

