package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// переменные для запуска
var (
	RunMode    = get(os.Getenv("FC_MODE"), "debug")
	RunAddr    = bget(RunMode == gin.ReleaseMode, "0.0.0.0", "")
	RunPort    = ":" + get(os.Getenv("FC_PORT"), "8080")
	RunPortSSL = ":" + get(os.Getenv("FC_SSL_PORT"), "4430")

	EnableSSL   = os.Getenv("FC_SSL") == "enable"
	SSLCertFile = os.Getenv("FC_SSL_CERT")
	SSLKeyFile  = os.Getenv("FC_SSL_KEY")

	JWTSecret = os.Getenv("FC_JWT_SECRET")
)

// безопасность
var (
	SecurityAllowedHosts    = strings.Split(os.Getenv("FC_SECURITY_ALLOWED_HOSTS"), ",")
	SecuritySSLRedirect     = os.Getenv("FC_SECURITY_SSL_REDIRECT") == "enable"
	SecuritySSLProxyHeaders = map[string]string{"X-Forwarded-Proto": "https"}
	SecurityFrameDeny       = true

	SecurityIsDevelopment = (RunMode != gin.ReleaseMode)

	expireSeconds, _ = strconv.Atoi(get(os.Getenv("FC_JWT_EXPIRE"), "2592000"))
	JWTExpireTime    = time.Second * time.Duration(expireSeconds)

	JWTSigningMethod = jwt.SigningMethodHS512
)

// misc
var (
	syncSeconds, _      = strconv.Atoi(get(os.Getenv("FC_SYNC_INTERVAL"), "2"))
	SyncIntervalSeconds = time.Duration(syncSeconds)
)
