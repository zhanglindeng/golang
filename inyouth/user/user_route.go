package user

import "dzlin.com/inyouth/conf"

func UserRoute(conf *conf.Conf) {
    userConf = conf

    user := conf.Router.Group("/user")
    user.Use(authRequired())
    {
        user.GET("/:name", getIndex)
    }

}
