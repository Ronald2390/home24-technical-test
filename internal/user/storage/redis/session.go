package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"home24-technical-test/internal/user"
	"home24-technical-test/internal/user/model"

	"github.com/go-redis/redis"
)

// SessionStorage represents the implementation of
type SessionStorage struct {
	redisClient *redis.Client
}

// FindByTokenAndType finds a session by its token & type
func (ss *SessionStorage) FindByTokenAndType(ctx context.Context, token string, sessType string) (*model.Session, error) {
	val, err := ss.redisClient.Get(fmt.Sprintf("%s:%s", sessType, token)).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var session model.Session
	if err := json.Unmarshal([]byte(val), &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// Insert inserts a new session
func (ss *SessionStorage) Insert(ctx context.Context, session *model.Session) error {
	sessionBytes, err := json.Marshal(session)
	if err != nil {
		return err
	}

	if err := ss.redisClient.Set(
		fmt.Sprintf("%s:%s", session.Type, session.ID),
		sessionBytes,
		session.ExpiredAt.Sub(time.Now())).Err(); err != nil {
		return err
	}

	return nil
}

// Update updates the user session
func (ss *SessionStorage) Update(ctx context.Context, session *model.Session) error {
	updatedSession, err := ss.FindByTokenAndType(ctx, session.ID, session.Type)
	if err != nil {
		return err
	}

	updatedSession.ExpiredAt = session.ExpiredAt
	updatedSession.User = session.User

	err = ss.Insert(ctx, updatedSession)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the session from storage
func (ss *SessionStorage) Delete(ctx context.Context, token string, sessType string) error {
	_, err := ss.redisClient.Del(fmt.Sprintf("%s:%s", sessType, token)).Result()
	if err != nil {
		return err
	}

	return nil
}

// UpdateByUserID  updates the user session based on user id
func (ss *SessionStorage) UpdateByUserID(ctx context.Context, session *model.Session) error {
	loginKeys := ss.redisClient.Keys(`*login*`)

	for _, loginKey := range loginKeys.Val() {
		splitLoginKey := strings.Split(loginKey, ":")

		currentSession, err := ss.FindByTokenAndType(ctx, splitLoginKey[1], user.LoginSessionType)
		if err != nil {
			return err
		}

		if currentSession != nil && currentSession.User.ID == session.User.ID {
			currentSession.User = session.User
			err = ss.Update(ctx, currentSession)
			if err != nil {
				return err
			}

			session = currentSession
		}
	}

	return nil
}

// DeleteByUserID delete session user
func (ss *SessionStorage) DeleteByUserID(ctx context.Context, userID int) error {
	loginKeys := ss.redisClient.Keys(`*login*`)

	for _, loginKey := range loginKeys.Val() {
		splitLoginKey := strings.Split(loginKey, ":")

		session, err := ss.FindByTokenAndType(ctx, splitLoginKey[1], user.LoginSessionType)
		if err != nil {
			return err
		}

		if session != nil && session.User.ID == userID {
			err = ss.Delete(ctx, splitLoginKey[1], user.LoginSessionType)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// NewSessionStorage creates a new session storage
func NewSessionStorage(redisClient *redis.Client) *SessionStorage {
	return &SessionStorage{
		redisClient,
	}
}
