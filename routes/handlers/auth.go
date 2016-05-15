package handlers

import (
    "time"

    "github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go"

    "github.com/gobricks/facecontrol/config"
    "github.com/gobricks/facecontrol/classes/response"
    "github.com/gobricks/facecontrol/classes/credentials"
)

// Auth gets uid/password pair and returns JWT token
func Auth(c *gin.Context) {
    uid := c.Query("uid")
    password := c.Query("password")

    if uid == "" || password == "" {
        c.JSON(400, response.BaseResponse{
            Success: false,
            Error: "Empty credentials",
        })
        return
    }

    // getting credeantials storage
    credentials := c.MustGet("credentials").(*credentials.Credentials)

    // looking for user
    found := credentials.Authenticate(uid, password)

    if found == false {
        c.JSON(404, response.BaseResponse{
            Success: false,
            Error: "User not found",
        })
        return
    }

    // creating token
    token := jwt.New(config.JWTSigningMethod)

    // adding claims
    token.Claims["uid"] = uid
    token.Claims["exp"] = time.Now().Add(config.JWTExpireTime).Unix()

    // signing token
    tokenString, err := token.SignedString([]byte(config.JWTSecret))

    if err != nil {
        c.JSON(500, response.BaseResponse{
            Success: false,
            Error: "Cannot sign token",
        })
        return
    }

    c.JSON(200, response.TokenSignResponse{
        BaseResponse: response.BaseResponse{
            Success: true,
        },
        Token: tokenString,
    })
}