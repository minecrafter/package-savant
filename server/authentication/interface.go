package authentication

import (
	"errors"
	"net/http"
)

// ErrIncorrectCredentials indicates that invalid credentials were provided.
var ErrIncorrectCredentials = errors.New("Invalid credentials.")

// Provider defines an authentication backend.
type Provider interface {
	Authenticate(r *http.Request) error
}
