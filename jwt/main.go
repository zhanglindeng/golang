package main

import (
    "github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go"
    "time"
    "fmt"
)

var secret = []byte("rmrS1qRXLjIYGB0QN6bzssuMSm16eH3oIPgx6SiuDAr4x9qTNPlzRaFKsvWRVFCH")

func main() {

    router := gin.New()
    router.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
        c.Next()
    })

    router.OPTIONS("/*cors", func(c *gin.Context) {})

    router.POST("/token", func(c *gin.Context) {
        claims := &jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
            Issuer:    "abc",
            NotBefore: time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
        }
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

        //ss, err := token.SignedString([]byte("VbP3BNPPFBSfmXKd32HfZUEgXN0lx69A"))
        //fmt.Printf("%v %v", ss, err)

        if tokenString, err := token.SignedString(secret); err != nil {
            c.JSON(500, gin.H{"code":1, "message":"server error"})
        } else {
            c.JSON(200, gin.H{"code":0, "token":tokenString, "message":"ok"})
        }
    })

    router.GET("/auth", func(c *gin.Context) {
        tokenString := c.Request.Header.Get("Authorization")
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Don't forget to validate the alg is what you expect:
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
            }
            return secret, nil
        })
        if err != nil {
            c.JSON(200, gin.H{"code":1, "message":err.Error()})
        } else {
            if token.Valid {
                // todo token.Claims["iss"]确定user
                fmt.Println(token.Claims)
                c.JSON(200, gin.H{"code":0, "message":"ok"})
            } else {
                c.JSON(200, gin.H{"code":1, "message":err.Error()})
            }
        }
    })

    router.Run(":50000")
}
