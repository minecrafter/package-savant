package authentication

import (
    "net/http"
)

// TestProvider compares the provided details to a set of static credentials: "sage" as the username, and "egas" as the password. Using it is not recommended.
type TestProvider struct {
}

func (TestProvider) Authenticate(r *http.Request) error {
    username, password, ok := r.BasicAuth()
    if !ok {
        return ErrIncorrectCredentials
    }
	if username == "sage" && password == "egas" {
		return nil
	}
	return ErrIncorrectCredentials
}
