package auth

import (
	"context"
	"errors"
)

// dummyOAuth2Provider is a placeholder implementation.
// In a real application, this would call the /userinfo endpoint
// of providers like Google, Okta, or Auth0.
type dummyOAuth2Provider struct {
}

func NewDummyOAuth2Provider() *dummyOAuth2Provider {
	return &dummyOAuth2Provider{}
}

func (p *dummyOAuth2Provider) VerifyToken(ctx context.Context, token string) (string, error) {
	if token == "valid-mock-token" {
		return "user@example.com", nil
	}
	return "", errors.New("invalid token")
}
