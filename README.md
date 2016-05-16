[![Build Status](https://travis-ci.org/gobricks/facecontrol.svg?branch=master)](https://travis-ci.org/gobricks/facecontrol)
[![Go Report Card](https://goreportcard.com/badge/github.com/gobricks/facecontrol)](https://goreportcard.com/report/github.com/gobricks/facecontrol)


# Facecontrol

Simple yet powerful authentication, single sign-on and (optinal) authorization solution.

# Basic example

Create file `main.go` and paste the following code into it:

``` go
package main

import (
    "github.com/gobricks/facecontrol"
    "github.com/gobricks/facecontrol/classes/credentials"
)

func main() {
    facecontrol.Run(getUsersData)
}

func getUsersData() credentials.Credentials {
    dataStorage := credentials.Credentials{}
    // key of map must be unique ID which will be used alongside
    // with password to identify user (e.g. username)
    dataStorage["gobricks"] = credentials.User{
        Password: "a963a333881c43fbca256730ff7d59235a3bd4ea",
        Payload: map[string]interface{}{
            "username": "gobricks",
            "fullname": "John Doe",
            "age": 24,
            "email": "jdoe@gobricks.io",
        },
    }
    
    // in real life you need to replace code above with
    // actual users data generator (e.g. mysql query or something like this)
    
    return dataStorage
}
```

# Build and run

```
$ go build main.go
$ FC_JWT_SECRET="replace_me_with_something_stronk" ./main
```

# Start up variables

Facecontrol can use certain environment variables to customize its behavour.

**Required**

* **FC_JWT_SECRET** _string_ - secret passphrase for token signing.

**Additional**

* **FC_MODE** _string_ - facecontrol run mode. Available options: `debug`, `release`. Default: `debug`
* **FC_HOST** _string_ - facecontrol run host. Default: `0.0.0.0`
* **FC_PORT** _int_ - facecontrol HTTP run port. Default: `8080`
* **FC_SSL** _string_ - enable facecontrol in HTTPS mode. Available options: `enable`, `disable`. Default: `disable`
* **FC_SSL_PORT** _int_ - facecontrol HTTPS run port. Default: `4430`
* **FC_SSL_CERT** _string_ - path to SSL cert file. Default: `""`
* **FC_SSL_KEY** _string_ - path to SSL key file. Default: `""`
* **FC_SYNC_INTERVAL** _int_ - interval between users credentials syncronization in seconds. Default: `2`
* **FC_JWT_EXPIRE** _int_ - JWT token expiration time in seconds. Default: `2592000`
* **FC_SECURITY_ALLOWED_HOSTS** _string_ - comma separated list of allowed domains. Default: `""`
* **FC_SECURITY_SSL_REDIRECT** _string_ - force SSL redirect. Available options: `enable`, `disable`. Default: `disable`

# Token issuing and validation

After calling ```facecontrol.Run()``` a web server will startup, allowing you to call two URLs:
* ```GET /token/``` - for token issuing
* ```GET /token/<token>/``` - for checking previously issued token

**Token validation example**:

```GET /token/?uid=gobricks&password=a963a333881c43fbca256730ff7d59235a3bd4ea```

Returns:

``` json
{
  "success": true,
  "token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjU5MTI3NjksInVpZCI6ImdvYnJpY2tzIn0.ITqJ1uMdNZXb9XfqbNVF-qy7hVTnPr5ZUk3SHf77y6MDb6_nBCxXN01Fo5M3jxP9o5DnCYV3Ic4OnIybb9qs1A"
}
```

**Token validation example**:

```GET /token/eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjU5MTI3NjksInVpZCI6ImdvYnJpY2tzIn0.ITqJ1uMdNZXb9XfqbNVF-qy7hVTnPr5ZUk3SHf77y6MDb6_nBCxXN01Fo5M3jxP9o5DnCYV3Ic4OnIybb9qs1A```

Returns:

``` json
{
  "success": true,
  "user": {
    "age": 24,
    "email": "jdoe@gobricks.io",
    "fullname": "John Doe",
    "username": "gobricks"
  }
}
```

# How it fits into your infastructure
![How it works](http://i.imgur.com/Cn2ImqX.jpg)

# How to achieve single sign-on

Just make session cookie available to any service hosted on your domain (e.g. *.mysite.com).

# How to achieve authorization

You can pass map with user priveleges alongside with primary authentication data in your ```main.go``` implementation. All your services will get this priveleges back after user authentication.

Example:

``` go
dataStorage["gobricks"] = credentials.User{
    Password: "a963a333881c43fbca256730ff7d59235a3bd4ea",
    Payload: map[string]interface{}{
        "username": "gobricks",
        "fullname": "John Doe",
        "age": 24,
        "email": "jdoe@gobricks.io",
        "priveleges": map[string]bool{
            "mail.read": true,
            "mail.delete": false,
        },
    },
}
```

Upon receiving user data from facecontrol your service can check if user can perform certain action based on available priveleges. 

# Important security notices

* It is highly recomended to use HTTPS for any facecontrol communications.
* You must never save into or pass to facecontrol user password in plain text. Use hashed version of password instead.
* Do not share JWT secret with any other services. It is ment to be unknown for everyone except facecontrol service.


