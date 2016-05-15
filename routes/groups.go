package routes

import (
    "github.com/gin-gonic/gin"

    "github.com/gobricks/facecontrol/routes/handlers"
)

// InitGroups enables all routes
func InitGroups(r *gin.Engine) {
    r.GET("/token/", handlers.Auth)
    r.GET("/token/:token/", handlers.Check)
}