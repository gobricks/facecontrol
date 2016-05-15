package middlewares

import (
    "github.com/gin-gonic/gin"
    "github.com/unrolled/secure"

    "github.com/gobricks/facecontrol/config"
)

var secureMiddleware = secure.New(secure.Options{
    AllowedHosts: config.SecurityAllowedHosts,
    SSLRedirect: config.SecuritySSLRedirect,
    SSLProxyHeaders: config.SecuritySSLProxyHeaders,
    FrameDeny: config.SecurityFrameDeny,

    IsDevelopment: config.SecurityIsDevelopment,
})

// SecurityChecker performs various security checks before request processing
func SecurityChecker() gin.HandlerFunc {
    return func(c *gin.Context) {
        err := secureMiddleware.Process(c.Writer, c.Request)

        // If there was an error, do not continue.
        if err != nil {
            c.Abort()
        }

        c.Next()
    }
}