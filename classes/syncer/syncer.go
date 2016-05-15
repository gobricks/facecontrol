package syncer

import (
    "time"

    "github.com/gobricks/facecontrol/classes/credentials"
    "github.com/gobricks/facecontrol/config"
)

// SyncCredentials syncs users credentials
func SyncCredentials(с *credentials.Credentials, generator credentials.Generator) {
    for range time.Tick(config.SyncIntervalSeconds * time.Second) {
        credentials := generator()
        *с = credentials
    }
}