package middlewares

import (
    "github.com/gobricks/facecontrol/classes/credentials"

    "github.com/gin-gonic/gin"
)

// CredentialsInjector injects pointer to users credentials storage into any request context
func CredentialsInjector(credentials *credentials.Credentials) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Set("credentials", credentials)
        c.Next()
    }
}