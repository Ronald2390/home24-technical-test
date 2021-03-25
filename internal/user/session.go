package user

import (
	"context"
	"home24-technical-test/internal/user/model"
	"time"
)

// session types
const (
	LoginSessionType = "login"
)

// SessionStorage represents the session storage interface
type SessionStorage interface {
	FindByTokenAndType(ctx context.Context, token string, sessType string) (*model.Session, error)
	Insert(ctx context.Context, session *model.Session) error
	Update(ctx context.Context, session *model.Session) error
	Delete(ctx context.Context, token string, sessType string) error
	UpdateByUserID(ctx context.Context, session *model.Session) error
	DeleteByUserID(ctx context.Context, userID int) error
}

// SessionServiceInterface represents the user session service interface
type SessionServiceInterface interface {
	GetSession(ctx context.Context, token string) (*model.Session, error)
	RemoveSession(ctx context.Context, token string) error
	ExtendingSessionTimeout(ctx context.Context, token string) (*model.Session, error)
	CreateSession(ctx context.Context, user *model.User, loginToken string) error
	UpdateSession(ctx context.Context, user *model.User) error
	DeleteSession(ctx context.Context, userID int) error
}

// SessionService is the domain logic implementation of user session service interface
type SessionService struct {
	sessionStorage SessionStorage
}

// GetSession get session by token given
func (s *SessionService) GetSession(ctx context.Context, token string) (*model.Session, error) {
	session, err := s.sessionStorage.FindByTokenAndType(ctx, token, LoginSessionType)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// RemoveSession remove session by token given
func (s *SessionService) RemoveSession(ctx context.Context, token string) error {
	session, err := s.GetSession(ctx, token)
	if err != nil {
		return err
	}

	if session != nil && session.User != nil {
		err = s.sessionStorage.Delete(ctx, token, LoginSessionType)
		if err != nil {
			return err
		}
	}

	return nil
}

// ExtendingSessionTimeout extends session expiration date for 2d
func (s *SessionService) ExtendingSessionTimeout(ctx context.Context, token string) (*model.Session, error) {
	session, err := s.sessionStorage.FindByTokenAndType(ctx, token, LoginSessionType)
	if err != nil {
		return nil, err
	}

	session.ExpiredAt = time.Now().Add(48 * time.Hour)
	if err := s.sessionStorage.Update(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// UpdateSession updates user session
func (s *SessionService) UpdateSession(ctx context.Context, user *model.User) error {
	session := &model.Session{}
	session.User = user
	err := s.sessionStorage.UpdateByUserID(ctx, session)
	if err != nil {
		return err
	}
	return nil
}

// DeleteSession delete user session
func (s *SessionService) DeleteSession(ctx context.Context, userID int) error {
	err := s.sessionStorage.DeleteByUserID(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

// CreateSession creates user session
func (s *SessionService) CreateSession(ctx context.Context, user *model.User, loginToken string) error {
	err := s.sessionStorage.DeleteByUserID(ctx, user.ID)
	if err != nil {
		return err
	}

	session := &model.Session{
		ID:        loginToken,
		Type:      LoginSessionType,
		ExpiredAt: time.Now().Add(48 * time.Hour),
		Info: map[string]interface{}{
			"UserID": user.ID,
		},
		User: user,
	}

	if err := s.sessionStorage.Insert(ctx, session); err != nil {
		return err
	}

	return nil
}

// NewSessionService creates a new user session service
func NewSessionService(
	sessionStorage SessionStorage,
) *SessionService {
	return &SessionService{
		sessionStorage: sessionStorage,
	}
}
