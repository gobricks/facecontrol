package handlers

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/gobricks/facecontrol/classes/credentials"
	"github.com/gobricks/facecontrol/classes/response"
	"github.com/gobricks/facecontrol/config"
)

// Check accepts JWT token string and returns user credentials
func Check(c *gin.Context) {
	tokenString := c.Param("token")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWTSecret), nil
	})

	if err != nil || token.Valid == false {
		c.JSON(400, response.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	credentials := c.MustGet("credentials").(*credentials.Credentials)
	uid := token.Claims["uid"].(string)

	user, err := credentials.Get(uid)
	if err != nil {
		c.JSON(400, response.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(200, response.UserResponse{
		BaseResponse: response.BaseResponse{
			Success: true,
		},
		User: user.Payload,
	})
}
