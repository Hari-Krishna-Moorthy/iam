package auth

import (
	"context"
	"fmt"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/permission"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/session"
	"github.com/golang-jwt/jwt/v5"
)

type jwtProvider struct {
	secretKey []byte
}

func NewJWTProvider(secret string) *jwtProvider {
	return &jwtProvider{secretKey: []byte(secret)}
}

type claims struct {
	UserID      string   `json:"user_id"`
	TenantID    string   `json:"tenant_id"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

func (p *jwtProvider) GenerateToken(ctx context.Context, s *session.Session) (string, error) {
	perms := make([]string, len(s.Permissions))
	for i, perm := range s.Permissions {
		perms[i] = perm.String()
	}

	c := claims{
		UserID:      s.UserID,
		TenantID:    s.TenantID,
		Role:        s.Role,
		Permissions: perms,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        s.ID,
			ExpiresAt: jwt.NewNumericDate(s.ExpiresAt),
			IssuedAt:  jwt.NewNumericDate(s.CreatedAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(p.secretKey)
}

func (p *jwtProvider) ValidateToken(ctx context.Context, tokenStr string) (*session.Session, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return p.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if c, ok := token.Claims.(*claims); ok && token.Valid {
		perms := make([]permission.Permission, 0, len(c.Permissions))
		for _, pStr := range c.Permissions {
			perm, _ := permission.Parse(pStr)
			perms = append(perms, perm)
		}

		return &session.Session{
			ID:          c.ID,
			UserID:      c.UserID,
			TenantID:    c.TenantID,
			Role:        c.Role,
			Permissions: perms,
			ExpiresAt:   c.ExpiresAt.Time,
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}
