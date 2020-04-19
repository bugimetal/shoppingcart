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

type AuthService interface {
	Authenticate(context.Context, User) (User, error)
}

type Auth struct {
}

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

func New() AuthService {
	return &Auth{}
}

func (a *Auth) Authenticate(ctx context.Context, user User) (User, error) {
	if err := user.Validate(); err != nil {
		return user, err
	}
	// faking user ID
	user.ID = 1

	return user, nil
}
