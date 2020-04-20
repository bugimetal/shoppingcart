package auth

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrNoUserName = errors.New("no user name provided")
	ErrNoPassword = errors.New("no password provided")
)

// AuthService describes an interface to Authenticate users
type AuthService interface {
	Authenticate(context.Context, User) (User, error)
}

type Auth struct {
}

// User describes basic user structure
type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Password []byte `json:"password"`
}

// Validate validates user for required params
func (u User) Validate() error {
	switch {
	case strings.TrimSpace(u.Name) == "":
		return ErrNoUserName
	case len(u.Password) == 0:
		return ErrNoPassword
	}
	return nil
}

// New returns a new authorization service
func New() AuthService {
	return &Auth{}
}

// Authenticate authenticates user by given name and password.
func (a *Auth) Authenticate(ctx context.Context, user User) (User, error) {
	if err := user.Validate(); err != nil {
		return user, err
	}

	switch user.Name {
	case "test":
		user.ID = 1
	case "hacker":
		user.ID = 2
	default:
		user.ID = 3
	}

	return user, nil
}
