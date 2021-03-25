package service

import (
	"context"
	"crypto/rand"
	"fmt"

	"home24-technical-test/internal/user"
	"home24-technical-test/internal/user/model"
	"home24-technical-test/internal/user/public"

	"golang.org/x/crypto/bcrypt"
)

// ServiceInterface represents the user application service interface
type ServiceInterface interface {
	ListUsers(ctx context.Context, params *public.FindAllUsersParams) ([]*model.User, error)
	GetUser(ctx context.Context, userID int) (*model.User, error)
	CreateUser(ctx context.Context, params *public.CreateUserParams) error
	UpdateUser(ctx context.Context, params *public.UpdateUserParams) (*model.User, error)
	DeleteUser(ctx context.Context, userID int) error
	ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error
	Login(ctx context.Context, params *public.LoginParams) (*public.LoginResponse, error)
	Logout(ctx context.Context, token string) error
	GetLoginSession(ctx context.Context, token string) (*model.Session, error)
}

// Service is the domain logic implementation of user Service interface
type Service struct {
	userService        user.ServiceInterface
	userSessionService user.SessionServiceInterface
}

// GetLoginSession gets the user login session from the session storage
func (s *Service) GetLoginSession(ctx context.Context, token string) (*model.Session, error) {
	session, err := s.userSessionService.GetSession(ctx, token)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// ListUsers is listing all Users
func (s *Service) ListUsers(ctx context.Context, params *public.FindAllUsersParams) ([]*model.User, error) {
	users, err := s.userService.ListUsers(ctx, params)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUser get user by its id
func (s *Service) GetUser(ctx context.Context, userID int) (*model.User, error) {
	user, err := s.userService.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates users data
func (s *Service) UpdateUser(ctx context.Context, params *public.UpdateUserParams) (*model.User, error) {
	updatedUser, err := s.userService.UpdateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	err = s.userSessionService.UpdateSession(ctx, updatedUser)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// ChangePassword changes user's password
func (s *Service) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	return s.userService.ChangePassword(ctx, userID, oldPassword, newPassword)
}

//DeleteUser deleting user and its session
func (s *Service) DeleteUser(ctx context.Context, userID int) error {
	err := s.userService.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	err = s.userSessionService.DeleteSession(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

// CreateUser creates a new user
func (s *Service) CreateUser(ctx context.Context, params *public.CreateUserParams) error {
	return s.userService.CreateUser(ctx, params)
}

// Login gets the user logged in the system
func (s *Service) Login(ctx context.Context, params *public.LoginParams) (*public.LoginResponse, error) {
	loggedUser, err := s.userService.GetUserByEmail(ctx, params.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(loggedUser.Password), []byte(params.Password)); err != nil {
		return nil, user.ErrWrongPassword
	}

	loginToken, errGenerateToken := generateToken()
	if errGenerateToken != nil {
		return nil, errGenerateToken
	}

	err = s.userSessionService.CreateSession(ctx, loggedUser, loginToken)
	if err != nil {
		return nil, err
	}

	return &public.LoginResponse{
		SessionID: loginToken,
		User:      loggedUser,
	}, nil
}

// Logout gets the user logged out from system
func (s *Service) Logout(ctx context.Context, token string) error {
	return s.userSessionService.RemoveSession(ctx, token)
}

func generateToken() (string, error) {
	buff := make([]byte, 32)
	_, err := rand.Read(buff)
	if err != nil {
		return "", err
	}
	token := fmt.Sprintf("%x", buff)
	return token, nil
}

// NewService creates a new user AppService
func NewService(
	userService user.ServiceInterface,
	userSessionService user.SessionServiceInterface,
) *Service {
	return &Service{
		userService:        userService,
		userSessionService: userSessionService,
	}
}
