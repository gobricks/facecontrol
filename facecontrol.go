package facecontrol

import (
    "github.com/gin-gonic/gin"

    "github.com/gobricks/facecontrol/config"
    "github.com/gobricks/facecontrol/middlewares"
    "github.com/gobricks/facecontrol/classes/credentials"
    "github.com/gobricks/facecontrol/classes/syncer"
    "github.com/gobricks/facecontrol/classes/helpers"

    apiRoutes "github.com/gobricks/facecontrol/routes"
)

// Run starts facecontrol
func Run(generator credentials.Generator) {
    r := gin.New()

    // injecting users credentials into every request context
    credentials := generator()
    r.Use(middlewares.CredentialsInjector(&credentials))

    // initializing background credentials sync
    go syncer.SyncCredentials(&credentials, generator)

    // enabling security
    r.Use(middlewares.SecurityChecker())

    // adding route handlers
    apiRoutes.InitGroups(r)
    
    // last checks before startup
    helpers.CheckBootstrap()

    // ignition!
    if config.EnableSSL == true {
        r.RunTLS(config.RunAddr + config.RunPortSSL, config.SSLCertFile, config.SSLKeyFile)
    } else {
	    r.Run(config.RunAddr + config.RunPort)
    }
}