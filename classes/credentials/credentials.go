package credentials

import (
	"errors"
)

// Generator is a function type which must be passed to facecontrol Run function
type Generator func() Credentials

// Credentials main users data storage
type Credentials map[string]User

// User single user data record
type User struct {
	Payload  interface{}
	Password string
}

// Authenticate checks uid/password pair against data storage
func (c Credentials) Authenticate(uid string, password string) bool {
	if user, ok := c[uid]; ok {
		return user.Password == password
	}

	return false
}

// Get fetches single user record from storage
func (c Credentials) Get(uid string) (User, error) {
	if user, ok := c[uid]; ok {
		return user, nil
	}

	return User{}, errors.New("User not found")
}
