package conf

import (
    "github.com/gin-gonic/gin"
    "github.com/astaxie/beego/session"
)

type Conf struct {
    Router  *gin.Engine
    Session *session.Manager
}