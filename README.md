[![Build Status](https://travis-ci.org/gobricks/facecontrol.svg?branch=master)](https://travis-ci.org/gobricks/facecontrol)
[![Go Report Card](https://goreportcard.com/badge/github.com/gobricks/facecontrol)](https://goreportcard.com/report/github.com/gobricks/facecontrol)


# Facecontrol

Simple authentication, single sign-on and (optional) authorization solution.

# Basic example

``` go
package main

import (
    "time"
    "errors"
    "net/http"

    "github.com/gobricks/facecontrol"
)

type MyUser struct {
    Login string    `json:"login"`
    FullName string `json:"fullname"`
    IsAdmin bool    `json:"is_admin"`
    CanEdit []string `json:"can_edit"`
}

func main() {
    fc, _ := facecontrol.New(facecontrol.Config{
        RunAt: ":8080",
        JwtSecret: "OpenSesame",
        JwtTTL: 24 * time.Hour,
        Validator: findUser,
    })
    
    fc.Run()
}

func findUser(r *http.Request) (facecontrol.Payload, error) {
    login := r.URL.Query().Get("login")
    password := r.URL.Query().Get("password")

    if login != "admin" && password != "12345" {
        return nil, errors.New("Invalid credentials")
    }

    return MyUser{
        Login: "admin",
        FullName: "Johnny Mnemonic",
        IsAdmin: true,
        CanEdit: []string{"posts", "comments"},
    }, nil
}
```

# Configuration

Use `facecontrol.Config` struct to customize Facecontrol behavior. Available fields are:

``` go
RunAt     string // defines address of running facecontrol instance. Example: "127.0.0.1:6000". Required
EnableSSL bool   // forces facecontrol to run in HTTPS mode
SSLCert   string // path to corresponding SSL file. Required if EnableSSL is true
SSLKey    string // path to corresponding SSL file. Required if EnableSSL is true
JwtSecret string // will be used to sign auth tokens. Required
JwtTTL    time.Duration // token expiration time
Validator CredentialsValidator // user define credentials validation function
```

# Validator function

A function with signature of `func(*http.Request) (facecontrol.Payload, error)` can be passed to `facecontrol.Config`.
If so every incoming HTTP request for token issuing will be passed to this function.
You can use this function to find user in your database or any other credential storage.
If given function return non-nil error user will be declined from acquiring token.

# Token issuing and validation

After calling ```facecontrol.Run()``` a web server will startup, allowing you to call two URLs:
* ```GET /issue``` - for token issuing
* ```GET /validate``` - for validating previously issued token

**Token validation example**:

```curl -X POST -F "login=admin" -F "password=d41d8cd98f00b204e9800998ecf8427e" "http://127.0.0.1:6000/issue"```

Returns:

```
eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjU5MTI3NjksInVpZCI6ImdvYnJpY2tzIn0.ITqJ1uMdNZXb9XfqbNVF-qy7hVTnPr5ZUk3SHf77y6MDb6_nBCxXN01Fo5M3jxP9o5DnCYV3Ic4OnIybb9qs1
```

**Token validation example**:

```curl -X GET -H "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE0ODYzMjAwODYsImlzcyI6ImZhY2Vjb250cm9sIn0.dZB-v4fx2x155YarTze17sQsq1HRpz0rYdIxF3hUG469-0l3N1RzE9ES1MFz8kPSWLaKUvXBAqXXDEEmNEb-DA" "http://127.0.0.1:6000/validate"```

Returns:

``` json
{
  "iat": 1486320086,
  "iss": "facecontrol",
  "data": {
      "login": "admin",
      "is_admin": true
  }
}
```

# How it fits into your infastructure
![How it works](http://i.imgur.com/Cn2ImqX.jpg)

# How to achieve single sign-on

Just make session cookie available to any service hosted on your domain (e.g. *.mysite.com).

# How to achieve authorization

You can pass user priveleges into token payload using `Validator` function.
All your services will get this priveleges back after user authentication.
See basic example.

Upon receiving user data from facecontrol your service can check if user can perform certain action based on available priveleges. 

# Important security notices

* It is highly recomended to use HTTPS for any facecontrol communications.
* You must never save into or pass to facecontrol user password in plain text. Use hashed version of password instead.
* Do not share JWT secret with any other services. It is ment to be unknown for everyone except facecontrol service.


