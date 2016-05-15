package helpers

import (
    "github.com/gobricks/facecontrol/config"
)

// CheckBootstrap checks if all startup preparations has been made
func CheckBootstrap()  {
    if config.JWTSecret == "" {
        panic("Required environment variable FC_JWT_SECRET has not been set")
    }
}