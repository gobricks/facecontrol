package facecontrol

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

// CredentialsValidator is an interface that defines communication
// logic with custom credentials storage. It accepts map of credentials
// and returns payload data if credentials are valid.
// If given credentials are invalid it must return nil.
type CredentialsValidator func(map[string]string) Payload

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

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Cannot parse form inputs: %s", err)))
		return
	}

	var inputs map[string]string
	if len(r.PostForm) > 0 {
		inputs = make(map[string]string)
		for k, v := range r.PostForm {
			inputs[k] = v[0]
		}
	}

	var credentials Payload
	if f.conf.Validator != nil {
		credentials = f.conf.Validator(inputs)
		if credentials == nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Invalid credentials given"))
			return
		}
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

	authHeader := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(authHeader) != 2 || authHeader[0] != "Bearer" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Malformed "Authorization" header`))
		return
	}

	token, err := jwt.Parse(authHeader[1], func(token *jwt.Token) (interface{}, error) {
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

	if token.Valid == false {
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
