package facecontrol

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Claims is a custom JWT claims
type Claims struct {
	Payload `json:"data,omitempty"`
	jwt.StandardClaims
}

// Payload is a custom data to be append to issued token.
type Payload interface{}

// CredentialsValidator is an function that defines incoming credentials
// validation and payload construction logic. It accepts raw HTTP Request
// and returns payload data if credentials are valid.
// If given credentials are invalid it must return non-nil error.
type CredentialsValidator func(*http.Request) (Payload, error)

// Facecontrol is an SSO service.
type Facecontrol struct {
	conf Config
}

// Config defines essential Faceontrol variables and parts.
type Config struct {
	// Webserver run address
	RunAt string
	// EnableSSL forces facecontrol to run in HTTPS mode
	EnableSSL bool
	// SSLCert and SSLKey are a paths to corresponding SSL files
	SSLCert string
	SSLKey  string
	// JwtSecret will be used to sign auth tokens
	JwtSecret string
	// JwtTTL will be used to set token expiration
	JwtTTL time.Duration
	// Validator is a credentials validating function.
	Validator CredentialsValidator
}

// New returns new instance of Facecontrol.
// Secret is a sign string for JWT token.
func New(conf Config) (*Facecontrol, error) {
	if conf.JwtSecret == "" {
		return nil, fmt.Errorf("JwtSecret cannot be empty")
	}

	if conf.Validator == nil {
		return nil, fmt.Errorf("Cannot start without validator function")
	}

	return &Facecontrol{
		conf: conf,
	}, nil
}

func (f Facecontrol) getExpiration() int64 {
	if f.conf.JwtTTL == 0 {
		return 0
	}

	return time.Now().Add(f.conf.JwtTTL).Unix()
}

func (f Facecontrol) issueToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Cannot parse form inputs: %s", err)))
		return
	}

	credentials, err := f.conf.Validator(r)
	if credentials == nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, Claims{
		Payload: credentials,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "facecontrol",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: f.getExpiration(),
		},
	})

	tokenString, err := token.SignedString([]byte(f.conf.JwtSecret))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Cannot sign token: %s", err)))
		return
	}

	w.Header().Set("Content-Type", "application/jwt")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenString))
}

func (f Facecontrol) validateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rawAuthHeader := r.Header.Get("Authorization")

	if len(rawAuthHeader) < 8 && rawAuthHeader[:7] != "Bearer " {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Malformed "Authorization" header`))
		return
	}

	tokenString := rawAuthHeader[7:len(rawAuthHeader)]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate the alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(f.conf.JwtSecret), nil
	})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if !token.Valid {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Cannot validate given token`))
		return
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Cannot find token payload`))
		return
	}

	payload, err := json.Marshal(token.Claims.(jwt.MapClaims))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Cannot marshal token payload`))
		return
	}

	w.Header().Set("Content-Type", "application/json+jwt")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

// Run starts Facecontrol instance.
func (f Facecontrol) Run() error {
	http.HandleFunc("/issue", f.issueToken)
	http.HandleFunc("/validate", f.validateToken)

	var err error
	if f.conf.EnableSSL {
		err = http.ListenAndServeTLS(f.conf.RunAt, f.conf.SSLCert, f.conf.SSLKey, nil)
	} else {
		err = http.ListenAndServe(f.conf.RunAt, nil)
	}

	return err
}
