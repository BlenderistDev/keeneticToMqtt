package auth

import (
	"errors"
	"fmt"

	"keeneticToMqtt/internal/clients/keenetic/auth"
)

type Auth struct {
	authClient *auth.Auth
}

func NewAuth(authClient *auth.Auth) *Auth {
	return &Auth{authClient: authClient}
}

func (a *Auth) RefreshAuth() error {
	realm, challenge, err := a.authClient.CheckAuth()
	switch {
	case errors.Is(err, auth.ErrUnauthorized):
		if err := a.authClient.Auth(realm, challenge); err != nil {
			return fmt.Errorf("error while keenetic auth: %w", err)
		}
	case err != nil:
		return fmt.Errorf("error while keenetic check auth: %w", err)
	}

	return nil
}
