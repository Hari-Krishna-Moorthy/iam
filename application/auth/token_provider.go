package auth

import (
	"context"
	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/session"
)

type TokenProvider interface {
	GenerateToken(ctx context.Context, session *session.Session) (string, error)
	ValidateToken(ctx context.Context, token string) (*session.Session, error)
}
