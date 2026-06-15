package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Hari-Krishna-Moorthy/multi-tenant-IAM/domain/session"
	"github.com/redis/go-redis/v9"
)

type sessionRepository struct {
	client *redis.Client
}

func NewSessionRepository(client *redis.Client) session.Repository {
	return &sessionRepository{client: client}
}

func (r *sessionRepository) Save(ctx context.Context, s *session.Session) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	// Store session details
	sessionKey := fmt.Sprintf("session:%s", s.ID)
	err = r.client.Set(ctx, sessionKey, data, time.Until(s.ExpiresAt)).Err()
	if err != nil {
		return err
	}

	// Add to user sessions list for Feature C (Single user can have multiple sessions)
	userSessionsKey := fmt.Sprintf("user_sessions:%s", s.UserID)
	return r.client.SAdd(ctx, userSessionsKey, s.ID).Err()
}

func (r *sessionRepository) GetByID(ctx context.Context, sessionID string) (*session.Session, error) {
	sessionKey := fmt.Sprintf("session:%s", sessionID)
	data, err := r.client.Get(ctx, sessionKey).Bytes()
	if err != nil {
		return nil, err
	}

	var s session.Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *sessionRepository) Delete(ctx context.Context, sessionID string) error {
	s, err := r.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}

	sessionKey := fmt.Sprintf("session:%s", sessionID)
	userSessionsKey := fmt.Sprintf("user_sessions:%s", s.UserID)

	// Remove from both the session store and the user's session list
	pipe := r.client.Pipeline()
	pipe.Del(ctx, sessionKey)
	pipe.SRem(ctx, userSessionsKey, sessionID)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *sessionRepository) GetByUserID(ctx context.Context, userID string) ([]*session.Session, error) {
	userSessionsKey := fmt.Sprintf("user_sessions:%s", userID)
	sessionIDs, err := r.client.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		return nil, err
	}

	sessions := make([]*session.Session, 0, len(sessionIDs))
	for _, id := range sessionIDs {
		s, err := r.GetByID(ctx, id)
		if err == nil {
			sessions = append(sessions, s)
		}
	}

	return sessions, nil
}
